# router

[![CI](https://github.com/faustbrian/go-router/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-router/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-router/actions/workflows/ci.yml)
[![Mutation](https://img.shields.io/badge/mutation-local_or_scheduled-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-structural_checks_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-router/v2.svg)](https://pkg.go.dev/github.com/faustbrian/go-router/v2)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-router?sort=semver)](https://github.com/faustbrian/go-router/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`router` is an explicit, immutable HTTP router built on Go's `net/http`
programming model. It adds deterministic composition, groups, names, safe URL
generation, metadata, introspection, mounts, and route-scoped middleware while
keeping handlers as ordinary `http.Handler` values.

The published v2.0.1 release retains the fail-closed inherited middleware
and mount defaults introduced in v2.0.0. The minimum supported toolchain
is Go 1.27.0. The package has no runtime dependencies and no global router,
reflection discovery, controller resolver, container, session, template, or
application lifecycle.

See the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/README.md)
and its [service-edge package guidance](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/design-language.md#package-families-and-selection)
for the shared design language and related packages.

## Installation

Install the published v2.0.1 release:

```sh
go get github.com/faustbrian/go-router/v2@v2.0.1
```

Existing v2 applications can retain released v2.0.0 while planning an
upgrade. Legacy consumers can retain
`github.com/faustbrian/go-router@v1.0.0`; do not bypass publication with a
local replacement or pseudo-version.

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
- [Support](SUPPORT.md)
- [Vulnerability reporting](SECURITY.md)
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

Use ordinary Go tests and a focused two-second fuzz smoke as described in
[CONTRIBUTING.md](CONTRIBUTING.md#verification). `make check` and `make ci`
use bounded local assurance with the public verifier; `make ci` checks the
specification register offline. `make check-full` retains full assurance and
explicit online authoritative-source validation. NilAway remains advisory.

## License

MIT. See [LICENSE](LICENSE) and [NOTICE](NOTICE).

## Documentation

Use the [documentation index](docs/README.md) for package-owned guides,
operational contracts, examples, and maintainer references.
