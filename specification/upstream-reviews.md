# Upstream authority review history

This append-only record preserves reviewed changes to the authorities monitored
by [`monitoring.json`](monitoring.json). A monitoring digest changes only after
the corresponding upstream delta has been classified against the applicable
specification decisions.

## 2026-09-03: Go releases feed

- **Authority:** `go-releases`
- **URL:** https://go.dev/dl/?mode=json&include=all
- **Previous SHA-256:**
  `c0696d6e5708cce9644204513ee4ecb994f3fd748065c54b0890158c34631207`
- **Reviewed SHA-256:**
  `638127a053a86576fc235aa196b26145c6f2fce8ce839ded767212a18d1c9415`
- **Retrieved and reviewed:** 2026-09-03
- **Applicability:** `ROUTER-DEC-001` and `ROUTER-DEC-005` for the newer
  escaped-path redirect behavior; no current decision for the Go 1.26.7 h2c
  deadline fix
- **Disposition:** Behavior-neutral for the selected Go 1.26.6 contracts.

The feed added Go 1.26.8 at commit
`c293dd49cbe25e1fe8d97d94a5cb618e7b6d831e` and Go 1.27.1 at commit
`862c888e612ac346c7c4d99c9392bdfd265f33b0` after the previous review.
Comparing Go 1.26.6 with Go 1.26.8 shows that `request.go` and `url.go` are
byte-identical. The only `server.go` change clears connection deadlines when an
HTTP/1 connection upgrades to unencrypted HTTP/2; the router does not own
server connections, deadlines, or h2c lifecycle.

Go 1.27 changed escaped-path `ServeMux` redirect serialization, which directly
intersects `ROUTER-DEC-001` and `ROUTER-DEC-005`. That change was separately
reviewed and the router now preserves the selected Go 1.26.6 redirect wire
behavior on later toolchains through focused executable evidence. No other
candidate decision is affected by the reviewed release delta.

The monitored `server.go`, `request.go`, and `url.go` URLs and SHA-256 values
were independently verified against the immutable `go1.26.6` tag. The selected
source and minimum supported Go version remain unchanged. Reconsider these
decisions if the minimum supported Go version changes or a versioned
compatibility decision adopts newer redirect serialization.

## 2026-09-03: RFC 9110 errata

- **Authority:** `rfc9110-errata`
- **URL:** https://errata.rfc-editor.org/search/?rfc_number=9110&presentation=records
- **Previous SHA-256:**
  `38bd006c96f8963d58573f704c5313a5f81968b90738c03ade0b036ec7bbdf4b`
- **Reviewed SHA-256:**
  `1f6790054c0cdb2f2a70a94fa2b9c73b09a4ee0578a32b4a3006ed0ecfaac86d`
- **Retrieved and reviewed:** 2026-09-03
- **Applicability:** `ROUTER-DEC-003`
- **Disposition:** Behavior-neutral because generated `Allow` field values
  already use comma-space separation.

[Errata ID 9162](https://errata.rfc-editor.org/eid9162/) was reported on
2026-09-01 as a Technical erratum against RFC 9110 Section 5.2. It proposes
changing the repeated-field combination wording from values separated by a
comma to values separated by comma plus space so the rule matches its example.
The erratum remains Reported, not Verified, and does not revise the immutable
RFC 9110 source.

`ROUTER-DEC-003` owns `Allow` construction, and the router already joins its
sorted method values with comma plus space into one field value. No selected
behavior, decision, source binding, or executable evidence changes. Reconsider
this disposition if Errata ID 9162 becomes Verified or the router later accepts
or combines repeated list-valued HTTP field lines.
