# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's documentation.
2. Run ordinary Go tests for the affected package before editing.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected entries in the
[specification decision register](docs/specification-decisions.md), and complete
the Specification Decisions section of the pull request template. An
unresolved interpretation or stale source pin is release-blocking; peer
behavior cannot silently select policy.

Selected mutation gates must finish with zero surviving viable mutants.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and ordinary CI; the public minimum remains
Go 1.27.0. An ordinary standalone checkout is sufficient for tests and a
focused active fuzz smoke, without sibling repositories or goal ledgers:

```bash
go test ./...
go test . -run '^$' -fuzz '^FuzzRequestTargets$' -fuzztime=2s
```

Select the native fuzz target affected by your change; ordinary tests already
exercise deterministic fuzz seeds. For shared checks, install the public
`golib` v2.1.1 release declared in `.golib.yaml` and verify its checksums:

```bash
make check
make ci
```

Both Make routes use bounded local assurance. `make ci` also checks the
specification register offline; it does not fetch authoritative sources. The
module declares no services or test tags. Optional mutation analysis can
select the root package:

```bash
golib mutation --module . --package .
```

For matching, HTTP composition or public-contract changes, select the relevant
existing conformance and owned integration checks:

```bash
make -f verification/package.mk conformance
make -f verification/package.mk integration
```

The integration check compares released v1 and the local v2 candidate with
pinned HTTP middleware, JSON-RPC and service modules in a disposable consumer;
it does not require sibling checkouts or prove published v2 availability.
Use `make check-full` for explicit online authoritative-source validation and
full assurance with the unchanged native fuzz budgets, v2 API baseline,
conformance, integration and benchmarks. The pinned Go 1.26.6 routing authority
is distinct from the Go 1.27.2 producer and must not drift silently.
Full risk campaigns remain available for explicitly selected scheduled or
release qualification, including [release rehearsal](docs/release.md).
The ordinary scheduled CI route remains bounded and separately monitors
authoritative specification sources online. Report unavailable or failing checks
without claiming release readiness. CI/configuration-only changes do not
require a library release.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
