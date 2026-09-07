# Compatibility Policy

This repository publishes one independently releasable root Go module. It uses
semantic versioning and root tags of the form `v<version>`. The stable v1 line
keeps patch and minor releases backward compatible; incompatible exported API
or documented behavior changes require a new major version.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. Deprecated APIs
follow [`DEPRECATION.md`](DEPRECATION.md).
