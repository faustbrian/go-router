package router

import "net/http"

// Middleware is the standard HTTP decorator shape. Construction happens during
// Compile, in reverse request order. Constructors must return a non-nil handler;
// callbacks own blocking, cancellation, re-entry, shared-state synchronization,
// partial responses and panics. The router adds no recovery or writer wrapper.
type Middleware = func(http.Handler) http.Handler

// MiddlewareExclusionPolicy controls whether a route may remove an inherited
// middleware layer. The zero value denies exclusion so security-sensitive
// middleware remains enforced unless its owner explicitly opts out.
type MiddlewareExclusionPolicy uint8

const (
	// MiddlewareExclusionDenied keeps inherited middleware non-excludable.
	MiddlewareExclusionDenied MiddlewareExclusionPolicy = iota
	// MiddlewareExclusionAllowed permits a route to exclude a named inherited layer.
	MiddlewareExclusionAllowed
)

// NamedMiddleware makes a middleware layer visible through introspection.
// Name may be empty when exclusion and duplicate detection are not needed.
// ExclusionPolicy defaults to MiddlewareExclusionDenied.
type NamedMiddleware struct {
	// Name identifies this layer for introspection and inherited exclusions.
	// Empty permits unnamed layers and skips name-based duplicate detection.
	Name string
	// Middleware constructs a non-nil handler at compile time. The function
	// and returned handler remain caller-owned and must support intended reuse.
	Middleware Middleware
	// ExclusionPolicy controls whether routes may remove this inherited layer.
	// The zero value denies removal; route-local middleware is never excluded.
	ExclusionPolicy MiddlewareExclusionPolicy
}

// Route is an explicit route descriptor. Builder.Register copies every slice
// and map before retaining it.
type Route struct {
	// Name is an optional stable generation/introspection name. Non-empty names
	// must be unique across the compiled table, including group prefixes.
	Name string
	// Methods is a non-empty set of distinct uppercase HTTP tokens. GET also
	// serves HEAD unless an explicit HEAD pattern wins; CONNECT is unsupported.
	Methods []string
	// Host is an optional ASCII hostname pattern without a port. Empty matches
	// any host; {name} captures exactly one label. IDNA conversion is caller-owned.
	Host string
	// Path is an absolute ServeMux literal/wildcard pattern, including {name},
	// {name...} and {$}. Invalid escapes and decoded dot segments are rejected.
	Path string
	// Handler serves the selected route and must not be nil or typed nil.
	// The router neither closes it nor recovers its panics or owns its state.
	Handler http.Handler
	// Middleware appends route-local layers after router and group layers.
	// Constructor and serving effects remain caller-owned.
	Middleware []NamedMiddleware
	// ExcludeMiddleware names inherited layers to remove only when each owner
	// allows exclusion. It never removes route-local layers; unknown names are
	// ignored, but empty, malformed or duplicate names fail compilation.
	ExcludeMiddleware []string
	// Metadata is copied descriptive data, not authorization state. Keys must
	// be non-empty, bounded and disjoint from inherited group metadata.
	Metadata map[string]string
	// Documentation is optional bounded introspection text, not executable policy.
	Documentation string
	// Operation is an optional bounded introspection identifier, not dispatch logic.
	Operation string
	// Source is an optional bounded caller label included in startup diagnostics.
	// It is not redacted; avoid secrets and personal data.
	Source string
}

func cloneRoute(route Route) Route {
	cloned := route
	cloned.Methods = append([]string(nil), route.Methods...)
	cloned.Middleware = append([]NamedMiddleware(nil), route.Middleware...)
	cloned.ExcludeMiddleware = append([]string(nil), route.ExcludeMiddleware...)
	if route.Metadata != nil {
		cloned.Metadata = make(map[string]string, len(route.Metadata))
		for key, value := range route.Metadata {
			cloned.Metadata[key] = value
		}
	}
	return cloned
}
