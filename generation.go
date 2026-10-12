package router

import (
	"net"
	"net/url"
	"strings"
	"unicode"
)

type parameterKind uint8

const (
	segmentParameter parameterKind = iota
	remainderParameter
	maxTrustedAuthorityBytes = 261
	maxRemainderSegments     = 65_536
)

// URLParameter is an explicitly typed named-route generation input. Construct
// values with Param or Remainder. The zero value is invalid during generation;
// values are reusable because constructors retain strings or copy segment slices
// and generation never mutates them.
type URLParameter struct {
	name      string
	values    []string
	kind      parameterKind
	oversized bool
}

// Param supplies one path segment or host label. It defers name/value validation
// to Path or URL. Path separators inside an otherwise valid segment are escaped
// as data, not expanded into remainder segments.
func Param(name, value string) URLParameter {
	return URLParameter{name: name, values: []string{value}, kind: segmentParameter}
}

// Remainder supplies explicit path segments for a remainder wildcard. Inputs
// above the package hard ceiling are rejected during generation without being
// copied. Smaller slices are defensively copied; strings remain immutable.
// Generation requires at least one segment, rejects empty/dot/slash-only segments
// and NUL/CR/LF, and still applies the router's lower parameter and byte budgets.
func Remainder(name string, segments ...string) URLParameter {
	if len(segments) > maxRemainderSegments {
		return URLParameter{name: name, kind: remainderParameter, oversized: true}
	}
	return URLParameter{name: name, values: append([]string(nil), segments...), kind: remainderParameter}
}

// BaseURL is an immutable absolute-URL origin constructed by NewBaseURL. Its
// zero value is invalid for URL. Validation is syntactic, not a decision that
// the origin is authorized or reachable; the caller owns that trust policy.
type BaseURL struct {
	scheme string
	host   string
}

// NewBaseURL validates one caller-trusted HTTP or HTTPS scheme and authority.
// It lowercases the scheme and accepts bounded ASCII authorities without
// userinfo, URL paths, controls or malformed ports. Authority spelling is retained;
// DNS resolution, IDNA conversion and proxy trust are not performed. Invalid
// syntax returns ErrGeneration; fixed byte ceilings return ErrLimitExceeded.
func NewBaseURL(scheme, authority string) (BaseURL, error) {
	if len(scheme) > len("https") {
		return BaseURL{}, generationError(ErrLimitExceeded, "scheme", "scheme is too long")
	}
	scheme = strings.ToLower(scheme)
	if scheme != "http" && scheme != "https" {
		return BaseURL{}, generationError(ErrGeneration, "scheme", "scheme must be http or https")
	}
	if len(authority) > maxTrustedAuthorityBytes {
		return BaseURL{}, generationError(ErrLimitExceeded, "host", "trusted authority is too long")
	}
	if authority == "" {
		return BaseURL{}, generationError(ErrGeneration, "host", "invalid trusted authority")
	}
	if invalidAuthority(authority) {
		return BaseURL{}, generationError(ErrGeneration, "host", "invalid trusted authority")
	}
	return BaseURL{scheme: scheme, host: authority}, nil
}

type generationRoute struct {
	host string
	path string
}

// Path generates a relative escaped path for a named route. Host wildcard
// values are intentionally not accepted by relative generation. Supply each path
// wildcard exactly once with the appropriate Param or Remainder kind; unknown,
// duplicate, missing, unused or unsafe values return ErrInvalidParameter. Names
// absent from the table return ErrGeneration, nil routers return ErrCompileState
// and budget excess returns ErrLimitExceeded. Failures return an empty string.
// Segments are escaped once with url.PathEscape (RFC 3986 sections 2.1/3.3;
// ROUTER-DEC-011), including literal percent and separator bytes as data.
func (r *Router) Path(name string, parameters ...URLParameter) (string, error) {
	if r == nil {
		return "", generationError(ErrCompileState, "router", "nil compiled router")
	}
	if len(name) > r.limits.MaxNameBytes {
		return "", generationError(ErrLimitExceeded, "name", "route name is too long")
	}
	route, exists := r.named[name]
	if !exists {
		return "", generationError(ErrGeneration, "name", "unknown route name")
	}
	provided, err := collectParameters(parameters, r.limits)
	if err != nil {
		return "", err
	}
	path, used, err := renderPath(route.path, provided)
	if err != nil {
		return "", err
	}
	if err := rejectUnused(provided, used); err != nil {
		return "", err
	}
	if len(path) > r.limits.MaxGeneratedURLBytes {
		return "", generationError(ErrLimitExceeded, "url", "generated path is too long")
	}
	return path, nil
}

// URL generates an absolute URL using a validated explicit base. A route host
// replaces the base hostname while retaining its trusted explicit port. No
// request or forwarding header supplies authority. Parameters obey Path's error
// rules and additionally require every host wildcard; a zero base or unknown name
// returns ErrGeneration. Nil query is accepted. Query keys are sorted by
// url.Values.Encode, values retain supplied order, and the query is not mutated;
// callers must not mutate its maps/slices concurrently. Empty value lists count
// toward admission budgets but emit nothing. Failures return an empty string.
func (r *Router) URL(name string, base BaseURL, query url.Values, parameters ...URLParameter) (string, error) {
	if r == nil {
		return "", generationError(ErrCompileState, "router", "nil compiled router")
	}
	if base.scheme == "" {
		return "", generationError(ErrGeneration, "base", "uninitialized base URL")
	}
	if base.host == "" {
		return "", generationError(ErrGeneration, "base", "uninitialized base URL")
	}
	if len(name) > r.limits.MaxNameBytes {
		return "", generationError(ErrLimitExceeded, "name", "route name is too long")
	}
	route, exists := r.named[name]
	if !exists {
		return "", generationError(ErrGeneration, "name", "unknown route name")
	}
	provided, err := collectParameters(parameters, r.limits)
	if err != nil {
		return "", err
	}
	path, used, err := renderPath(route.path, provided)
	if err != nil {
		return "", err
	}
	host, hostUsed, err := renderHost(route.host, base.host, provided)
	if err != nil {
		return "", err
	}
	for name := range hostUsed {
		used[name] = struct{}{}
	}
	if err := rejectUnused(provided, used); err != nil {
		return "", err
	}
	if err := validateQuery(query, r.limits); err != nil {
		return "", err
	}
	generated := base.scheme + "://" + host + path
	if encoded := query.Encode(); encoded != "" {
		generated += "?" + encoded
	}
	if len(generated) > r.limits.MaxGeneratedURLBytes {
		return "", generationError(ErrLimitExceeded, "url", "generated URL is too long")
	}
	return generated, nil
}

func collectParameters(parameters []URLParameter, limits Limits) (map[string]URLParameter, error) {
	if len(parameters) > limits.MaxURLParameters {
		return nil, generationError(ErrLimitExceeded, "parameters", "parameter count exceeded")
	}
	provided := make(map[string]URLParameter, len(parameters))
	bytes := 0
	values := 0
	for _, parameter := range parameters {
		if parameter.oversized {
			return nil, generationError(ErrLimitExceeded, "parameters", "remainder segment count exceeded")
		}
		if len(parameter.name) > limits.MaxWildcardNameBytes {
			return nil, generationError(ErrLimitExceeded, "parameters", "parameter bytes exceeded")
		}
		bytes += len(parameter.name)
		if bytes > limits.MaxURLParameterBytes {
			return nil, generationError(ErrLimitExceeded, "parameters", "parameter bytes exceeded")
		}
		values += len(parameter.values)
		if values > limits.MaxURLParameters {
			return nil, generationError(ErrLimitExceeded, "parameters", "parameter value count exceeded")
		}
		// Count raw names and every remainder value before escaping expands data.
		// MaxURLParameters caps both descriptors and total supplied segment values.
		for _, value := range parameter.values {
			bytes += len(value)
			if bytes > limits.MaxURLParameterBytes {
				return nil, generationError(ErrLimitExceeded, "parameters", "parameter bytes exceeded")
			}
		}
		if !validWildcardName(parameter.name) {
			return nil, generationError(ErrInvalidParameter, "parameters", "invalid parameter name")
		}
		if _, duplicate := provided[parameter.name]; duplicate {
			return nil, generationError(ErrInvalidParameter, "parameters", "duplicate parameter")
		}
		provided[parameter.name] = parameter
	}
	return provided, nil
}

func validWildcardName(value string) bool {
	if value == "" {
		return false
	}
	for index, character := range value {
		if unicode.IsLetter(character) || character == '_' ||
			index > 0 && unicode.IsDigit(character) {
			continue
		}
		return false
	}
	return true
}

func renderPath(pattern string, provided map[string]URLParameter) (string, map[string]struct{}, error) {
	segments := strings.Split(pattern, "/")
	output := make([]string, 0, len(segments))
	used := make(map[string]struct{})
	for _, segment := range segments {
		if segment == "{$}" {
			output = append(output, "")
			continue
		}
		name, kind, wildcard := pathWildcard(segment)
		if !wildcard {
			output = append(output, segment)
			continue
		}
		parameter, exists := provided[name]
		if !exists {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "missing path parameter")
		}
		if parameter.kind != kind {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "wrong parameter kind")
		}
		if len(parameter.values) == 0 {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "wrong parameter kind")
		}
		if kind == segmentParameter && len(parameter.values) != 1 {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "wrong parameter kind")
		}
		for _, value := range parameter.values {
			if !safePathSegment(value) {
				return "", nil, generationError(ErrInvalidParameter, "parameters", "unsafe path segment")
			}
			output = append(output, url.PathEscape(value))
		}
		used[name] = struct{}{}
	}
	return strings.Join(output, "/"), used, nil
}

func renderHost(pattern, baseAuthority string, provided map[string]URLParameter) (string, map[string]struct{}, error) {
	if pattern == "" {
		return baseAuthority, map[string]struct{}{}, nil
	}
	labels := strings.Split(pattern, ".")
	used := make(map[string]struct{})
	for index, label := range labels {
		if !isWildcardLabel(label) {
			continue
		}
		name := label[1 : len(label)-1]
		parameter, exists := provided[name]
		if !exists {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "missing host parameter")
		}
		if parameter.kind != segmentParameter {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "unsafe host label")
		}
		if len(parameter.values) != 1 {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "unsafe host label")
		}
		if !safeHostLabel(parameter.values[0]) {
			return "", nil, generationError(ErrInvalidParameter, "parameters", "unsafe host label")
		}
		labels[index] = parameter.values[0]
		used[name] = struct{}{}
	}
	host := strings.Join(labels, ".")
	parsed, _ := url.Parse("//" + baseAuthority)
	if port := parsed.Port(); port != "" {
		host = net.JoinHostPort(host, port)
	}
	return host, used, nil
}

func pathWildcard(segment string) (string, parameterKind, bool) {
	if len(segment) < 3 {
		return "", segmentParameter, false
	}
	if segment[0] != '{' {
		return "", segmentParameter, false
	}
	if segment[len(segment)-1] != '}' {
		return "", segmentParameter, false
	}
	name := segment[1 : len(segment)-1]
	if strings.HasSuffix(name, "...") {
		return strings.TrimSuffix(name, "..."), remainderParameter, true
	}
	return name, segmentParameter, true
}

func rejectUnused(provided map[string]URLParameter, used map[string]struct{}) error {
	for name := range provided {
		if _, exists := used[name]; !exists {
			return generationError(ErrInvalidParameter, "parameters", "unknown or unused parameter")
		}
	}
	return nil
}

func safePathSegment(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	if strings.Trim(value, "/") == "" {
		return false
	}
	return !strings.ContainsAny(value, "\x00\r\n")
}

func safeHostLabel(value string) bool {
	if value == "" {
		return false
	}
	if len(value) > 63 {
		return false
	}
	if value[0] == '-' || value[len(value)-1] == '-' {
		return false
	}
	for _, character := range value {
		if asciiLetterOrDigit(character) || character == '-' {
			continue
		}
		return false
	}
	return true
}

func validateQuery(values url.Values, limits Limits) error {
	if len(values) > limits.MaxQueryValues {
		return generationError(ErrLimitExceeded, "query", "query value count exceeded")
	}
	count := 0
	bytes := 0
	for key, entries := range values {
		bytes += len(key)
		if bytes > limits.MaxQueryBytes {
			return generationError(ErrLimitExceeded, "query", "query bytes exceeded")
		}
		// An empty list emits no query pair but still consumes one admission
		// slot. Otherwise many empty keys could bypass the count budget.
		if len(entries) == 0 {
			count++
		} else {
			count += len(entries)
			for _, entry := range entries {
				bytes += len(entry)
				if bytes > limits.MaxQueryBytes {
					return generationError(ErrLimitExceeded, "query", "query bytes exceeded")
				}
			}
		}
		if count > limits.MaxQueryValues {
			return generationError(ErrLimitExceeded, "query", "query value count exceeded")
		}
	}
	return nil
}

func generationError(kind error, field, detail string) error {
	return &Error{Kind: kind, Field: field, Detail: detail}
}

func ascii(value string) bool {
	for _, character := range value {
		if character > 127 {
			return false
		}
	}
	return true
}
