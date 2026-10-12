package router

import (
	"context"
	"maps"
	"net"
	"net/http"
	"net/url"
	pathpkg "path"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type routeContextKey struct{}
type pathValuesContextKey struct{}

// RouteInfo is a defensive snapshot of a compiled route, omitting handlers and
// function identity. Its exported fields are mutable caller-owned data: changing
// them does not change dispatch. Methods, Parameters, Middleware and Metadata
// are independently copied whenever information crosses the API.
type RouteInfo struct {
	// Name is the optional stable route name after group composition.
	Name string
	// Methods contains sorted explicitly registered methods, not implied HEAD
	// or automatic OPTIONS additions.
	Methods []string
	// Host is the registered composed host pattern, not the request authority.
	Host string
	// Pattern is the composed path pattern, not the concrete requested path.
	Pattern string
	// Parameters lists host wildcards left-to-right, then path wildcards in
	// pattern order, with no duplicate names.
	Parameters []string
	// Middleware lists resolved layers in request order; unnamed layers appear
	// as empty strings, and excluded inherited layers are omitted.
	Middleware []string
	// Metadata is a copied descriptive map with the composed route entries.
	Metadata map[string]string
	// Documentation is the retained bounded descriptive text.
	Documentation string
	// Operation is the retained optional identifier, not handler identity.
	Operation string
	// Source is the retained caller diagnostic label; it is not redacted.
	Source string
}

type compiledHost struct {
	pattern   string
	mux       *http.ServeMux
	methods   map[string]struct{}
	redirects []*http.ServeMux
}

// Router has immutable compiled indexes supporting concurrent ServeHTTP,
// Routes, Path and URL calls. Construct it through Builder.Compile; its zero
// value has no initialized routing budgets. Retained handlers/middleware must
// provide their own concurrency safety and resource lifetime. Router starts no
// goroutines, does no background I/O and has no Close operation.
type Router struct {
	limits           Limits
	hosts            []compiledHost
	table            []RouteInfo
	supported        map[string]struct{}
	named            map[string]generationRoute
	notFound         http.Handler
	methodNotAllowed http.Handler
	automaticOptions bool
	redirectPolicy   RedirectPolicy
}

// Routes returns independent snapshots sorted by name, host, path and method
// set, with methods sorted within each entry. It copies the full table and nested
// collections, so cost is linear in the published data. A nil Router returns nil;
// a compiled empty router returns a non-nil empty slice.
func (r *Router) Routes() []RouteInfo {
	if r == nil {
		return nil
	}
	table := make([]RouteInfo, len(r.table))
	for index, info := range r.table {
		table[index] = cloneRouteInfo(info)
	}
	return table
}

// MatchedRoute returns the route selected for request, when called from its
// handler or middleware chain. It returns a defensive snapshot, or zero/false
// for nil requests or absent package route context. Nested routers install their
// own matched metadata; this function reports the context on the supplied view.
func MatchedRoute(request *http.Request) (RouteInfo, bool) {
	if request == nil {
		return RouteInfo{}, false
	}
	info, ok := request.Context().Value(routeContextKey{}).(RouteInfo)
	if !ok {
		return RouteInfo{}, false
	}
	return cloneRouteInfo(info), true
}

// ServeHTTP dispatches synchronously without mutating compiled state. It
// rejects nil router/request/URL, invalid authorities/methods and CONNECT with
// 400, and excess targets with 414, before route middleware or custom miss
// handlers. It follows RedirectPolicy, explicit routes, automatic OPTIONS,
// known-path 405, unsupported-method 501 and ordinary 404 in that order.
//
// Handlers receive standard PathValue entries and matched metadata on a derived
// context. URL and RequestURI remain unchanged except for a mounted handler's
// optional cloned URL view; ServeMux owns its ordinary request path-value state.
// Supply a non-nil writer. The router does not recover panics, wrap the writer,
// stop canceled requests automatically or constrain caller callback execution.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if r == nil || request == nil || request.URL == nil {
		writeError(writer, http.StatusBadRequest, "bad request")
		return
	}
	if invalidAuthority(request.Host) {
		writeError(writer, http.StatusBadRequest, "bad request")
		return
	}
	if len(request.Method) > r.limits.MaxMethodBytes {
		writeError(writer, http.StatusBadRequest, "bad request")
		return
	}
	if !isToken(request.Method) || request.Method == http.MethodConnect {
		writeError(writer, http.StatusBadRequest, "bad request")
		return
	}
	if requestTargetTooLong(request, r.limits.MaxRequestTargetBytes) {
		writeError(writer, http.StatusRequestURITooLong, "request target too long")
		return
	}
	if request.RequestURI == "*" {
		if request.Method != http.MethodOptions {
			writeError(writer, http.StatusBadRequest, "bad request")
			return
		}
		if !r.automaticOptions {
			r.serveNotFound(writer, request)
			return
		}
		allowed := r.allMethods()
		if len(allowed) > 0 {
			writer.Header().Set("Allow", strings.Join(allowed, ", "))
		}
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	if r.redirectPolicy == FollowRedirects && nonCanonicalPath(request.URL.EscapedPath()) {
		servePinnedCanonicalRedirect(writer, request)
		return
	}

	hosts := r.matchingHosts(authorityHost(request.Host))
	if host := r.matchingHostForMethod(hosts, request); host != nil {
		if host.requiresSubtreeRedirect(request) {
			servePinnedSubtreeRedirect(writer, request)
			return
		}
		host.mux.ServeHTTP(writer, request)
		return
	}
	allowed := r.allowedMethods(hosts, request)
	if request.Method == http.MethodOptions && r.automaticOptions && len(allowed) > 0 {
		writer.Header().Set("Allow", strings.Join(allowed, ", "))
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	if len(allowed) > 0 {
		writer.Header().Set("Allow", strings.Join(allowed, ", "))
		if r.methodNotAllowed != nil {
			r.methodNotAllowed.ServeHTTP(writer, request)
		} else {
			http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		return
	}
	if _, supported := r.supported[request.Method]; !supported && request.Method != http.MethodOptions {
		writeError(writer, http.StatusNotImplemented, "method not implemented")
		return
	}
	r.serveNotFound(writer, request)
}

func requestTargetTooLong(request *http.Request, limit int) bool {
	if len(request.RequestURI) > limit || len(request.URL.Path) > limit {
		return true
	}
	if len(request.URL.RawPath) > limit || len(request.URL.RawQuery) > limit {
		return true
	}
	escapedPath := request.URL.EscapedPath()
	if len(escapedPath) > limit {
		return true
	}
	return len(request.URL.RawQuery) > limit-len(escapedPath)
}

func (r *Router) serveNotFound(writer http.ResponseWriter, request *http.Request) {
	if r.notFound != nil {
		r.notFound.ServeHTTP(writer, request)
		return
	}
	http.NotFound(writer, request)
}

// Route context must surround the entire resolved middleware chain. Preserve
// outer path values on mounted routers except names owned by this inner route;
// copying a separate parameter API would diverge from Request.PathValue.
// ROUTER-DEC-009/010 and mount_policy_test.go protect the collision rule.
func (r *Router) matchingHandler(infoIndex int, hostParameters []string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		info := cloneRouteInfo(r.table[infoIndex])
		current := make(map[string]struct{}, len(info.Parameters))
		for _, name := range info.Parameters {
			current[name] = struct{}{}
		}
		inherited, _ := request.Context().Value(pathValuesContextKey{}).(map[string]string)
		for name, value := range inherited {
			if _, shadowed := current[name]; !shadowed {
				request.SetPathValue(name, value)
			}
		}
		values, _ := matchHost(r.table[infoIndex].Host, authorityHost(request.Host))
		for index, name := range hostParameters {
			request.SetPathValue(name, values[index])
		}
		resolved := make(map[string]string)
		for name, value := range inherited {
			resolved[name] = value
		}
		for _, name := range info.Parameters {
			resolved[name] = request.PathValue(name)
		}
		ctx := context.WithValue(request.Context(), pathValuesContextKey{}, resolved)
		request = request.WithContext(context.WithValue(ctx, routeContextKey{}, info))
		handler.ServeHTTP(writer, request)
	})
}

func (r *Router) matchingHosts(host string) []*compiledHost {
	matches := make([]*compiledHost, 0, len(r.hosts))
	for index := range r.hosts {
		if _, matched := matchHost(r.hosts[index].pattern, host); matched {
			matches = append(matches, &r.hosts[index])
		}
	}
	return matches
}

func (r *Router) matchingHostForMethod(hosts []*compiledHost, request *http.Request) *compiledHost {
	for _, host := range hosts {
		if r.hostMatchesRequest(host, request) {
			return host
		}
	}
	return nil
}

func (r *Router) hostMatchesRequest(host *compiledHost, request *http.Request) bool {
	if r.redirectPolicy == RejectRedirects && host.requiresRedirect(request) {
		return false
	}
	_, pattern := host.mux.Handler(request)
	if pattern == "" {
		return false
	}
	return r.redirectPolicy == FollowRedirects || !nonCanonicalPath(request.URL.EscapedPath())
}

func (r *Router) allowedMethods(hosts []*compiledHost, request *http.Request) []string {
	allowed := make(map[string]struct{})
	for _, host := range hosts {
		for method := range host.methods {
			candidate := request.Clone(request.Context())
			candidate.Method = method
			if r.hostMatchesRequest(host, candidate) {
				allowed[method] = struct{}{}
				if method == http.MethodGet {
					allowed[http.MethodHead] = struct{}{}
				}
			}
		}
	}
	if len(allowed) == 0 {
		return nil
	}
	if r.automaticOptions {
		allowed[http.MethodOptions] = struct{}{}
	}
	methods := make([]string, 0, len(allowed))
	for method := range allowed {
		methods = append(methods, method)
	}
	sort.Strings(methods)
	return methods
}

func (h *compiledHost) requiresRedirect(request *http.Request) bool {
	if nonCanonicalPath(request.URL.EscapedPath()) {
		return true
	}
	return h.requiresSubtreeRedirect(request)
}

func (h *compiledHost) requiresSubtreeRedirect(request *http.Request) bool {
	for _, redirects := range h.redirects {
		if _, pattern := redirects.Handler(request); pattern != "" {
			_, matchedPattern := h.mux.Handler(request)
			_, matchedPath, ok := strings.Cut(matchedPattern, " ")
			return ok && redirectRoot(matchedPath) != ""
		}
	}
	return false
}

func redirectRoot(pattern string) string {
	if strings.HasSuffix(pattern, "/") {
		return strings.TrimSuffix(pattern, "/")
	}
	segments := strings.Split(pattern, "/")
	last := segments[len(segments)-1]
	if last == "{$}" || strings.HasSuffix(last, "...}") {
		return strings.TrimSuffix(strings.TrimSuffix(pattern, last), "/")
	}
	return ""
}

// ROUTER-DEC-001 pins Go 1.26.6 redirect serialization rather than the producer
// toolchain's form. Putting the cleaned escaped path in Path deliberately escapes
// existing percent signs again; switching to RawPath changes Location on the
// wire. TestPinnedServeMuxRedirectsPreserveGo126EscapedPathBehavior freezes it.
func servePinnedCanonicalRedirect(writer http.ResponseWriter, request *http.Request) {
	redirect := &url.URL{
		Path:     canonicalPath(request.URL.EscapedPath()),
		RawQuery: request.URL.RawQuery,
	}
	http.RedirectHandler(redirect.String(), http.StatusTemporaryRedirect).ServeHTTP(writer, request)
}

// The pinned subtree form uses decoded Path, unlike canonical redirects above.
// Preserve that observable distinction and the original query (ROUTER-DEC-001).
func servePinnedSubtreeRedirect(writer http.ResponseWriter, request *http.Request) {
	redirect := &url.URL{
		Path:     canonicalPath(request.URL.Path) + "/",
		RawQuery: request.URL.RawQuery,
	}
	http.RedirectHandler(redirect.String(), http.StatusTemporaryRedirect).ServeHTTP(writer, request)
}

func nonCanonicalPath(value string) bool {
	return canonicalPath(value) != value
}

func canonicalPath(value string) string {
	if value == "" {
		return "/"
	}
	if value[0] != '/' {
		value = "/" + value
	}
	cleaned := pathpkg.Clean(value)
	if strings.HasSuffix(value, "/") && cleaned != "/" {
		cleaned += "/"
	}
	return cleaned
}

func (r *Router) allMethods() []string {
	methods := slices.Collect(maps.Keys(r.supported))
	if len(methods) > 0 {
		methods = append(methods, http.MethodOptions)
	}
	slices.Sort(methods)
	methods = slices.Compact(methods)
	return methods
}

func routeInformation(route Route, middleware []string, limits Limits) (RouteInfo, []string, error) {
	parameters := make([]string, 0)
	hostParameters := make([]string, 0)
	seen := make(map[string]struct{})
	for _, label := range strings.Split(route.Host, ".") {
		if isWildcardLabel(label) {
			name := label[1 : len(label)-1]
			if _, duplicate := seen[name]; duplicate {
				return RouteInfo{}, nil, &Error{Kind: ErrInvalidRoute, Field: "host", Source: route.Source, Detail: "duplicate wildcard name"}
			}
			seen[name] = struct{}{}
			parameters = append(parameters, name)
			hostParameters = append(hostParameters, name)
		}
	}
	for _, segment := range strings.Split(route.Path, "/") {
		if len(segment) < 3 || segment[0] != '{' || segment[len(segment)-1] != '}' || segment == "{$}" {
			continue
		}
		name := strings.TrimSuffix(segment[1:len(segment)-1], "...")
		if _, duplicate := seen[name]; duplicate {
			return RouteInfo{}, nil, &Error{Kind: ErrInvalidRoute, Field: "path", Source: route.Source, Detail: "duplicate wildcard name"}
		}
		seen[name] = struct{}{}
		parameters = append(parameters, name)
	}
	if len(parameters) > limits.MaxURLParameters {
		return RouteInfo{}, nil, &Error{Kind: ErrLimitExceeded, Field: "parameters", Source: route.Source, Detail: "parameter count exceeded"}
	}
	return RouteInfo{
		Name:          route.Name,
		Methods:       append([]string(nil), route.Methods...),
		Host:          route.Host,
		Pattern:       route.Path,
		Parameters:    parameters,
		Middleware:    append([]string(nil), middleware...),
		Metadata:      cloneMetadata(route.Metadata),
		Documentation: route.Documentation,
		Operation:     route.Operation,
		Source:        route.Source,
	}, hostParameters, nil
}

func cloneRouteInfo(info RouteInfo) RouteInfo {
	cloned := info
	cloned.Methods = append([]string(nil), info.Methods...)
	cloned.Parameters = append([]string(nil), info.Parameters...)
	cloned.Middleware = append([]string(nil), info.Middleware...)
	cloned.Metadata = cloneMetadata(info.Metadata)
	return cloned
}

func cloneMetadata(metadata map[string]string) map[string]string {
	cloned := make(map[string]string, len(metadata))
	for key, value := range metadata {
		cloned[key] = value
	}
	return cloned
}

func matchHost(pattern, host string) ([]string, bool) {
	if pattern == "" {
		return nil, true
	}
	patternLabels := strings.Split(pattern, ".")
	hostLabels := strings.Split(host, ".")
	if len(patternLabels) != len(hostLabels) {
		return nil, false
	}
	values := make([]string, 0)
	for index, label := range patternLabels {
		if isWildcardLabel(label) {
			if hostLabels[index] == "" {
				return nil, false
			}
			values = append(values, hostLabels[index])
			continue
		}
		if !strings.EqualFold(label, hostLabels[index]) {
			return nil, false
		}
	}
	return values, true
}

func isWildcardLabel(label string) bool {
	return len(label) > 2 && label[0] == '{' && label[len(label)-1] == '}'
}

func authorityHost(authority string) string {
	host, _, err := net.SplitHostPort(authority)
	if err == nil {
		return strings.Trim(host, "[]")
	}
	return strings.Trim(authority, "[]")
}

func invalidAuthority(authority string) bool {
	if authority == "" {
		return false
	}
	if len(authority) > maxTrustedAuthorityBytes {
		return true
	}
	if !ascii(authority) || strings.ContainsAny(authority, "\x00\r\n\\/@") {
		return true
	}
	parsed, err := url.Parse("//" + authority)
	if err != nil || parsed.Host != authority || parsed.User != nil ||
		parsed.Hostname() == "" || parsed.Path != "" {
		return true
	}
	if port := parsed.Port(); port != "" {
		number, parseErr := strconv.Atoi(port)
		return parseErr != nil || number < 1 || number > 65_535
	}
	return strings.HasSuffix(authority, ":")
}

func writeError(writer http.ResponseWriter, status int, message string) {
	http.Error(writer, message, status)
}
