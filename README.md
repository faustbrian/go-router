# router

[![CI](https://github.com/faustbrian/go-router/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-router/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-router/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-router.svg)](https://pkg.go.dev/github.com/faustbrian/go-router)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-router?sort=semver)](https://github.com/faustbrian/go-router/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`router` is an explicit, immutable HTTP router built on Go's `net/http`
programming model. It adds deterministic composition, groups, names, safe URL
generation, metadata, introspection, mounts, and route-scoped middleware while
keeping handlers as ordinary `http.Handler` values.

The minimum supported toolchain is Go 1.26.6. The package has no runtime
dependencies and no global router, reflection discovery, controller resolver,
container, session, template, or application lifecycle.

## Five-minute start

```go
builder := router.New()
err := builder.Register(router.Route{
    Name:    "users.show",
    Methods: []string{http.MethodGet},
    Path:    "/users/{id}",
    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, r.PathValue("id"))
    }),
})
if err != nil {
    return err
}

handler, err := builder.Compile()
if err != nil {
    return err
}
return http.ListenAndServe(":8080", handler)
```

Registration is single-owner and startup-time. The compiled router is an
immutable `http.Handler` safe for concurrent serving and introspection.

## Contracts

- [Semantics](docs/semantics.md)
- [Specification decisions](docs/specification-decisions.md)
- [Pinned conformance matrix](specification/README.md)
- [Compatibility](docs/compatibility.md)
- [Behavior matrices](docs/matrices.md)
- [Resource limits](docs/limits.md)
- [Security](docs/security.md)
- [Architecture](docs/architecture.md)
- [Five-minute quickstarts](docs/quickstart.md)
- [API reference](docs/api.md)
- [Adoption guides](docs/adoption.md)
- [Migration guides](docs/migration.md)
- [Cookbook](docs/cookbook.md)
- [Performance](docs/performance.md)
- [FAQ](docs/faq.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Release process](docs/release.md)

## Development

Run `make check` for the blocking local checks. NilAway remains visible as an
advisory result in the shared contract. Each target is independently
reproducible.

## License

MIT. See [LICENSE](LICENSE) and [NOTICE](NOTICE).

## Documentation

Use the [documentation index](docs/README.md) for package-owned guides,
operational contracts, examples, and maintainer references.
