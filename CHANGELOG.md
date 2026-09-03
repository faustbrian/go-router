# Changelog

All notable changes to this project are documented in this file.

The format is based on Keep a Changelog, and this project adheres to Semantic
Versioning.

## [Unreleased]

### Changed

- Adopt the checksum-verified `go-library-tools` v1.4.0 CLI and immutable
  shared workflow while preserving every existing repository gate.
- Include authoritative online specification validation in the complete local
  CI contract enforced by `make ci`.
- Replace copied repository tooling with the checksum-pinned
  `go-library-tools` v1.2.0 specification-governance contract while retaining
  package-owned policy and verification evidence.
- Pin CI to the immutable shared-tooling revision that enforces the structured
  specification decision, conformance, history, and monitoring contract.
- Expose the owned middleware, JSON-RPC, and service integrations through the
  typed interoperability gate required by the repository contract.
- Adopt the checksum-verified `go-library-tools` v1.3.0 CLI, schema-v2
  cohesion metadata, and repository-local cohesion gate while retaining the
  router's public API and package-owned verification evidence.
- Pin reusable CI to the immutable v1.3.0 workflow so hosted checks enforce
  the same cohesion contract as local validation.
- Preserve the pinned Go 1.26.6 escaped-path redirect serialization across
  later Go toolchains instead of silently adopting Go 1.27 wire changes.

### Documentation

- Preserve the reviewed Go release and RFC 9110 errata dispositions in an
  append-only upstream authority history.
- Make the [specification decision register](docs/specification-decisions.md)
  machine-auditable with monitored source and errata authorities, attributable
  conformance evidence, classified interoperability results, and durable
  decision history.

  - ROUTER-DEC-001 sha256:56f4ea57efd1fb24b8650c883bea47f05872c63d6fb4f4fb43f1f2bdafe9ff8c
  - ROUTER-DEC-001 sha256:f446182ad093b6e3dc1f40e180fe8ad03f88733233cf8a3f2722740c3b804585
  - ROUTER-DEC-002 sha256:f8dbfe4fe493cc3b71c400f985206aef6654f266f705ff0a6f8fdbb9fd979811
  - ROUTER-DEC-003 sha256:e57444cf3c6e408a366c339ed4f622897c655ccf85dfa6545628d07e04cd5964
  - ROUTER-DEC-003 sha256:131ba5fada36d81eb47b1d595b624ee2cf0b95e6cfdfde47409b8492e156d2f3
  - ROUTER-DEC-004 sha256:d06af40fb0a74ae46f2f22d843c09effd2c397777bf07d2a9f931fdec3e9efb2
  - ROUTER-DEC-005 sha256:383d99f58bc5bf65b234941feb386c096b673e5276e45cba86aa31e5d80971e9
  - ROUTER-DEC-005 sha256:35abe45420556f55ed3d65b43fb6f2d172710ada3949468b4a0526aac74db0d4
  - ROUTER-DEC-005 sha256:adf00e74ff1fee7ab18248fef929a805e30d4dfd14ddaa69d67859d70af701f0
  - ROUTER-DEC-006 sha256:9fa23204c8ff2ae12d5631c7947091b32436be3bf075a29f2700037a7bce27c1
  - ROUTER-DEC-007 sha256:1eb13de30f37d6bc1e8e7924ace29a967364d88f22570ca7cc84f2c5659e9139
  - ROUTER-DEC-008 sha256:bcf6f5c4c6b7b3960e028c33645cedbb22dcfab2d7460d2acee73077edee2c2b
  - ROUTER-DEC-009 sha256:d9b4e24c30bcd894c78acf8b92d8bea117053a186f64b34cdb98bd6f24a5cc50
  - ROUTER-DEC-010 sha256:b5d8f22178d8134c9f8d51926982579efddbd5f06abe7e73ed3bbb323718d4a7
  - ROUTER-DEC-011 sha256:344c1ed6825f17af34fc9ed9be101d2bf730a2ce7efd316091f28cc505b14bc1
  - ROUTER-DEC-012 sha256:061923f106754b9924c7895baeb6aaba09b137e011f365f227d1be701641c513
  - ROUTER-DEC-013 sha256:a016415608cdc167fa9420223fd33a537d5594acef8d03d5d0dc90fe658d42ce
  - ROUTER-DEC-013 sha256:2c3dc12f61f2cfd486dcf93fc8a264758733eb61b13a0249920c06778fe81317

- Replace archived monorepo links and completed execution artifacts with a
  standalone, human-oriented documentation structure.
- Link the router to the immutable v1.3.0 Golib ecosystem index and
  service-edge package-selection guidance.
- Bind monitored Go source files to immutable Go 1.26.6 URLs and record the
  reviewed Go 1.27 redirect change without changing the support baseline.
- Record RFC 9110 Erratum 9162 as behavior-neutral because generated `Allow`
  field values already use comma-space separation.

## [1.0.0] - 2026-08-25

### Fixed

- Exercise router, middleware, JSON-RPC, and service composition through
  versioned module archives without sibling checkouts or Git fetches.

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Replace obsolete standalone-repository links and workflow claims with
  monorepo-canonical targets and current release guidance.

- Link the package README to package-owned documentation.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-router` identity while preserving its documented API and behavior.
- Link conformance and contribution guidance directly to the canonical
  specification decision register.
- Delegate local mutation checks to the canonical exact-100 repository runner
  instead of accepting package-local survivors.

### Added

- A pinned specification decision register and focused conformance gate for
  Go routing, HTTP semantics, request-target forms, and URI generation.

- An explicit request-target byte budget enforced before dispatch matching.
- Blocking architecture checks for production goroutines and process-global
  HTTP registration.
- Explicit validated route descriptors with typed bounded errors and limits.
- Immutable `ServeMux`-backed compilation, standard path values, host patterns,
  deterministic method handling, and safe route introspection.
- Transactional nested groups and visible router, group, route, and mount
  middleware composition.
- Safe named path and absolute URL generation with typed remainder segments.
- Explicit handler mounts, automatic OPTIONS control, redirect policy, and
  customizable 404 and 405 handlers.
- `routertest` consumer helpers, differential compatibility fixtures, fuzzing,
  race tests, mutation checks, benchmarks, and full adoption documentation.
- Pinned local and GitHub Actions release gates with signed-tag verification and
  build provenance attestations.
- Expanded conflict, path-value, middleware, mount, and URL-security truth
  tables with executable panic, cancellation, and partial-write evidence.
- Expanded standard-library differential coverage across all standard method
  classes and nested group fuzz properties.
- Documented migration guidance for every unsupported `ServeMux` pattern and
  resource-boundary difference.
- Recorded the complete local release-gate results and refreshed the measured
  performance baseline.
- Proved malformed requests bypass custom miss handlers and route middleware.

### Fixed

- Bound middleware identifiers, exclusions, mount prefixes, schemes, route
  name lookups, and request methods before expensive work.
- Bound diagnostic input before UTF-8 normalization and control sanitization.
- Accept the complete `ServeMux` wildcard identifier set during named-route
  generation.
- Preserve `ServeMux` canonical redirects before route and method miss
  classification.
- Strip encoded literal mount prefixes in decoded path space while preserving
  the escaped suffix for mounted handlers.
- Convert only controlled `ServeMux` registration errors while allowing
  runtime faults to propagate.
- Reject slash-only generated wildcard values that cannot round-trip through
  `ServeMux` as one segment.
- Reject middleware chains that resolve to a nil handler during compilation.
- Reject unsupported CONNECT routes and bound documentation, empty-group
  metadata, and trusted absolute-URL authorities during startup.
- Bound method tokens, wildcard identifiers, URL parameter input, and raw
  query input before allocation-heavy parsing or encoding.
- Count host and path wildcards together against the per-route budget.
- Reject oversized route collection fields before copying caller-owned values,
  including middleware exclusion lists.
- Bound the total segment values supplied to remainder URL parameters.
- Reject remainder constructors above a fixed segment ceiling without copying
  the oversized caller slice.
- Validate router-wide middleware against the final option limits before
  registration while retaining defensive copies and option-order independence.
- Validate the complete `ServeMux` pattern set before constructing any
  middleware, keeping conflict failures free of partial handler graphs.
- Apply named inherited middleware exclusions to group layers as well as
  router-wide layers.
- Enforce path and name budgets on composed nested groups even when a group
  callback registers no routes.
- Propagate remaining route and group capacity into child builders and reject
  exhausted group counts before invoking another callback.
- Return `ErrLimitExceeded` consistently for syntactically valid route names,
  hosts, paths, and group prefixes that exceed configured byte budgets.
- Evaluate rejected redirects against escaped paths so encoded separators and
  dot text remain wildcard data instead of false canonicalization misses.
- Match rejected subtree roots with standard patterns across wildcard, Unicode,
  percent-escape, exact-root, and implied HEAD semantics.
- Sanitize rendered diagnostics to bounded single-line valid UTF-8 without
  splitting multibyte characters.
- Align default 404 and non-automatic 405 responses with `http.ServeMux`.
- Document and freeze every dispatch difference caused by automatic OPTIONS,
  unsupported method misses, host extensions, redirect policy, and CONNECT.

[Unreleased]: https://github.com/faustbrian/go-router/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-router/releases/tag/v1.0.0
