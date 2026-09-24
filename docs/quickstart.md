# Five-minute Quickstarts

## Route

Create a `Builder`, call `Register` with explicit methods, path, and an
`http.Handler`, then call `Compile`. Hand the returned `*Router` to
`http.Server` or `service`. Read wildcards with `request.PathValue`.

## Group

Call `Group` with `GroupOptions` and register routes in its callback. Host,
path, name, metadata, and middleware composition is transactional and flattened
into the compiled table. Nested callbacks execute outer middleware first.

```go
err := builder.Group(router.GroupOptions{
    PathPrefix: "/api/v1",
    NamePrefix: "api.v1.",
}, func(api *router.Builder) error {
    return api.Register(usersRoute)
})
```

## Middleware

Use `WithMiddleware` for router-wide layers and `Route.Middleware` for route
layers. A `NamedMiddleware` exposes a stable identifier in `Router.Routes`.
Request order is router, outer group, inner group, route; responses unwind in
reverse. `Route.ExcludeMiddleware` can explicitly remove named inherited
layers only when the layer owner sets
`ExclusionPolicy: router.MiddlewareExclusionAllowed`. Keep authentication and
authorization middleware at the default non-excludable policy.

## Mount

Call `Mount` with any `http.Handler`. `StripPrefix` clones the request and URL
before stripping, while retaining the original `RequestURI`. A nested compiled
router, JSON-RPC endpoint, webhook, health probe, metrics endpoint, or debug
endpoint remains an ordinary handler with caller-owned lifecycle.
Default mounts omit TRACE; include it in `MountOptions.Methods` only when the
mounted handler deliberately supports it.

## Named path and URL

Name a route and call `Path` with `Param` values. Use `Remainder` for a
`{path...}` wildcard. Absolute generation additionally uses a `BaseURL`
created by `NewBaseURL`; it never reads a request or forwarding header. Query
values use `url.Values`.

All quickstarts are compiled and executed in `example_test.go`.
