# Release Process

The current source module prepares v2 and is not published. The public badge
and installation guidance remain bound to released v1. Before v2 can be
declared releasable, publish the first v2 tag and migrate the owned
`go-http-middleware/integration/siblings` and
`go-service/integration/reference-http` consumers from released v1 without
local replacements. Until then, aggregate shared-tool reconciliation and those
published-module consumer migrations remain release blockers.

1. Update `CHANGELOG.md`, the API baseline, compatibility notes, and docs.
2. Run `make check` with the pinned Go and shared-tool versions.
3. Review coverage, mutation output, vulnerability results, fuzz smoke,
   benchmarks, and integration evidence.
4. Open and merge a reviewed pull request with every blocking workflow green.
5. Create a signed `vMAJOR.MINOR.PATCH` tag. The release workflow verifies the
   tag commit, reruns the blocking gate, builds provenance, and publishes the
   changelog entry.

Route names and exported APIs are SemVer contracts after v1. Security releases
must not disclose a private advisory before coordinated publication.
