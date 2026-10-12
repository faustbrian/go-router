# Source documentation audit

## Audited boundary

Baseline `fde5afea62dd6de3878414ffcc6c91cc7e96a719` contains one root
module and 25 tracked Go files across `router` and `routertest`, including
internal/external tests, six fuzz targets, five benchmarks and five examples.
There are no nested modules, vendor files, generated Go sources, embedded
executable templates, build constraints or platform-specific variants.

Reproduce the inventory with `git ls-files '*.go' 'go.mod' '**/go.mod'` and
`golib source-docs inventory --json`. Both package comments existed; the
67 objective baseline findings were undocumented exported fields and
`routertest.TestingT` methods. Semantic inspection additionally found terse
builder/callback/zero-value contracts, a stale v1 reference on `ErrUnsupported`,
and no adjacent explanation for the deliberately pinned redirect encoding.
No debt, deprecation or generation directives require resolution.

| Source boundary | Contract reviewed and documentation outcome |
| --- | --- |
| `builder.go`, `compiler.go` | Single-owner startup state; copied collections versus retained callbacks; deferred option errors; success-only freezing; validation versus constructor side effects; bounded controlled panic conversion. |
| `route.go`, `group.go`, `mount.go` | Exported members, transactional descriptor publication, child-builder lifetime, non-excludable inherited security layers, metadata collisions and cloned mounted URL views. |
| `router.go` | Immutable internal indexes versus mutable caller-owned handlers; defensive introspection; miss/redirect order; nested path-value ownership. |
| `generation.go`, `limits.go`, `errors.go` | Explicit trusted origins, segment escaping, parameter/query counting, positive caller-selected maxima and diagnostic sanitization versus redaction. |
| `routertest` | Test-double obligations, synchronous bodyless requests, nil/error outcomes and caller-owned cancellation/handler lifetime. |
| All tracked tests and corpus | Registration/conflict/copy boundaries, private doubles, pinned redirects, deterministic hostile seeds, six bounded fuzz input domains, benchmark setup and request-clone costs, unchanged executable examples. |

The actual policy authority remains the existing
[specification decision register](specification-decisions.md), including the
Go 1.26.6 matching/redirect contract, RFC 9110 method policy and RFC 3986
component generation. This audit does not reopen or claim new conformance to
those authorities. Go 1.27.2 remains the producer toolchain.

## Preservation and native checks

This pass changes prose, not executable tokens, literals, test assertions,
example output directives, dependency selections or specification decisions.
Handlers and middleware remain caller-owned; comments do not certify arbitrary
callbacks as safe for concurrent use or make their effects transactional.

The configured API baseline uses full `go-doc` mode for `.` and `./routertest`.
Refresh `api/v2.txt` through `golib api update`, then use `golib api check`;
`api/baseline.txt` is historical v1 material, not this generation target.
The native `scripts/check-docs.sh` checks local documentation links and runs
existing examples and vet. Objective source checks and changed-content secret
scanning complement ordinary complete-diff and rendered-document inspection.
No unchanged broad race, active fuzz, mutation, performance, external-service,
conformance or integration campaign is required by this prose-only pass.

Local source checks do not establish remote-main delivery, hosted Required CI
or adoption of the new shared source-documentation enforcement. The coordinator
owns those delivery boundaries. This documentation pass requires no module
release or downstream repinning.

## Objective Enforcement

The qualified public Tools v2.4.0 binary and immutable reusable CI source
provide the strict tracked-source check through `make source-docs-check`,
the existing local check/CI routes and Required aggregation. Release and
selected security/assurance routes retain their existing controls. This
structural status supplements the rendered semantic audit; it does not
certify prose meaning, generator execution or runtime conformance.
