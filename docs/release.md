# Release Process

The root module uses `github.com/faustbrian/go-router/v2` and root `v2.0.0`
tags on main. Eligibility does not establish publication. The owned
`go-http-middleware/integration/siblings` and
`go-service/integration/reference-http` consumers must migrate and verify
against the public v2 module without local replacements after publication.

1. Update `CHANGELOG.md`, the API baseline, compatibility notes, and docs.
2. Run the affected proportional-assurance gates with pinned Go and tooling;
   reuse valid evidence for unchanged source rather than repeating every gate.
3. Independently review the complete final diff and require exact-source CI
   green on remote main before creating a signed `vMAJOR.MINOR.PATCH` tag.
4. Verify the tag, public module source, and a fresh public consumer without
   replacements. Migrate the directly owned consumers reached by the change.
5. Publish the release notes, actual SBOM, honest provenance, and signed
   checksums; refetch and verify the published assets. A locally prepared
   payload or signed tag alone is not completed release evidence.

Route names and exported APIs are SemVer contracts after v1. Security releases
must not disclose a private advisory before coordinated publication.
