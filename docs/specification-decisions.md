# Router Specification Decisions

This register records observable routing choices where HTTP, URI syntax, Go's
`net/http`, and package policy permit different outcomes. The package claims
only the scoped behavior below; it does not implement an HTTP wire stack,
proxy, application framework, or alternate copy of Go's matcher.

Each resolved entry names executable evidence. Changing one requires
compatibility, security, resource, API, conformance, and changelog review.
Superseded decisions remain linked from their replacements.

## ROUTER-DEC-001: ServeMux as the matching authority

- **Status, owner, and classification:** `resolved`; router maintainers; Go
  interoperability policy.
- **Source and issue:** Go 1.26.6 [`ServeMux`](https://cs.opensource.google/go/go/+/refs/tags/go1.26.6:src/net/http/server.go)
  defines patterns, specificity, wildcards, path values, and redirects, but a
  router could copy or reinterpret those rules.
- **Interpretations and peer behavior:** Delegate, copy internals, implement a
  separate trie, or support only literal paths. Third-party routers commonly
  differ on precedence and escaped segments.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Delegate supported literal,
  `{name}`, `{name...}`, and `{$}` path parsing, specificity, conflict,
  extraction, GET-to-HEAD, and redirects to the pinned `ServeMux`. Package
  extensions operate outside that matcher; no copied internal source, unsafe,
  or registration-order tie-break exists.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestSupportedMatchingIsDifferentialWithServeMux`,
  `TestSupportedMethodsAndLiteralHostsAreDifferentialWithServeMux`, and
  `FuzzRoutePatternCompilation` cover `Builder.Register`, `Compile`, and
  `Router.ServeHTTP`. Reconsider when the minimum Go routing contract changes.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-001",
  "title": "ServeMux as the matching authority",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "interoperability policy",
  "decision_scope": "normative",
  "specification": "Go 1.26.6 net/http and net/url contracts",
  "version": "Go 1.26.6",
  "source_authority": "go-net-http",
  "section": "net/http ServeMux patterns, precedence, redirects, and path values",
  "requirement_strength": "not specified",
  "issue": "The standard library defines matching behavior, but a router can delegate to it, copy it, or define incompatible precedence and escaping rules.",
  "interpretations": [
    "Delegate supported behavior to ServeMux",
    "Copy internal matching behavior",
    "Implement an independent matcher",
    "Support literal paths only"
  ],
  "peer_behavior": "Go 1.26.6 ServeMux is the maintained differential peer for supported precedence, escaped segments, redirects, and path-value extraction.",
  "selected_behavior": "Delegate supported literal, named wildcard, remainder wildcard, and end-marker parsing, specificity, conflicts, path values, GET-to-HEAD matching, and redirects to Go 1.26.6 ServeMux.",
  "rationale": "Delegation preserves the declared Go contract and avoids an unreviewed parallel matcher.",
  "security_consequences": "No copied unsafe matcher or registration-order tie-break can reinterpret escaped path structure.",
  "resource_consequences": "ServeMux owns matcher construction and dispatch within the router's finite registration and request limits.",
  "compatibility_consequences": "Supported patterns track the pinned Go 1.26.6 contract; package extensions remain explicitly separate.",
  "wire_consequences": "Dispatch, redirects, and extracted path values match ServeMux for the supported surface.",
  "executable_evidence": [
    "TestSupportedMatchingIsDifferentialWithServeMux",
    "TestSupportedMethodsAndLiteralHostsAreDifferentialWithServeMux"
  ],
  "fixture_evidence": [
    "compatibility_security_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRoutePatternCompilation"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "Builder.Register",
    "Builder.Compile",
    "Router.ServeHTTP"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "Go 1.26.6 source is the selected upstream authority.",
  "reconsider_when": "The minimum supported Go routing contract changes."
}
```

Authority URL: https://go.dev/src/net/http/server.go?m=text

</details>

## ROUTER-DEC-002: Registration panics, conflicts, and publication

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  startup policy.
- **Source and issue:** Go 1.26.6
  [`ServeMux.Handle`](https://cs.opensource.google/go/go/+/refs/tags/go1.26.6:src/net/http/server.go)
  panics for malformed or conflicting patterns. A reusable startup API must
  decide which panics become errors and whether failed compilation mutates
  publication state.
- **Interpretations and peer behavior:** Preserve every panic, recover every
  panic, preimplement conflict logic, or convert only controlled synchronous
  registration failures. Frameworks vary between panic and error APIs.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Validate and sort the complete table
  before constructing middleware; convert only non-runtime panics produced by
  package-owned `ServeMux.Handle` calls into bounded typed errors; propagate
  unrelated panics. Failed compile publishes nothing and leaves the builder
  repairable; successful compile freezes it.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestCompileReturnsTypedConflictsAndFreezesOnlyOnSuccess`,
  `TestConflictFailureDoesNotConstructAnyMiddleware`,
  `TestRegistrationOrderDoesNotChangeDispatchOrIntrospection`, and
  `TestPatternValidationPropagatesUncontrolledPanics` cover `Compile` and
  `Error`. Reconsider if Go replaces registration panics with a typed API.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-002",
  "title": "Registration panics, conflicts, and publication",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "defensive",
  "specification": "Go 1.26.6 net/http and net/url contracts",
  "version": "Go 1.26.6",
  "source_authority": "go-net-http",
  "section": "net/http ServeMux.Handle registration failures",
  "requirement_strength": "not specified",
  "issue": "ServeMux registration panics, but a reusable startup API must decide which panics become errors and whether failed compilation publishes partial state.",
  "interpretations": [
    "Preserve every panic",
    "Recover every panic",
    "Reimplement conflict detection",
    "Convert only controlled ServeMux registration failures"
  ],
  "peer_behavior": "Go 1.26.6 ServeMux is the maintained peer and panics synchronously for malformed or conflicting registrations.",
  "selected_behavior": "Validate and sort the complete route table before middleware construction, convert only non-runtime panics from package-owned ServeMux.Handle calls into bounded typed errors, publish nothing after failure, and freeze only after success.",
  "rationale": "Controlled conversion provides a repairable startup API without hiding runtime faults or publishing a partial handler graph.",
  "security_consequences": "Unrelated runtime panics remain visible and failed validation cannot expose a partially constructed router.",
  "resource_consequences": "Whole-table validation completes before middleware allocation and publication.",
  "compatibility_consequences": "Callers receive stable typed compile errors for controlled conflicts and unchanged panics for unrelated faults.",
  "wire_consequences": "Failed compilation produces no serving behavior; the previous compiled router remains unchanged.",
  "executable_evidence": [
    "TestCompileReturnsTypedConflictsAndFreezesOnlyOnSuccess",
    "TestConflictFailureDoesNotConstructAnyMiddleware",
    "TestRegistrationOrderDoesNotChangeDispatchOrIntrospection",
    "TestPatternValidationPropagatesUncontrolledPanics"
  ],
  "fixture_evidence": [
    "router_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRoutePatternCompilation"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "Builder.Compile",
    "Error"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "Go 1.26.6 retains panic-based ServeMux registration.",
  "reconsider_when": "Go replaces registration panics with a typed public API."
}
```

Authority URL: https://go.dev/src/net/http/server.go?m=text

</details>

## ROUTER-DEC-003: Method, HEAD, OPTIONS, Allow, and miss outcomes

- **Status, owner, and classification:** `resolved`; maintainers; RFC 9110
  semantics with explicit `ServeMux` divergences.
- **Source and issue:** RFC 9110 [methods](https://www.rfc-editor.org/rfc/rfc9110.html#section-9),
  [HEAD](https://www.rfc-editor.org/rfc/rfc9110.html#section-9.3.2),
  [OPTIONS](https://www.rfc-editor.org/rfc/rfc9110.html#section-9.3.7), and
  [`Allow`](https://www.rfc-editor.org/rfc/rfc9110.html#section-10.2.1) do not
  dictate a router's automatic OPTIONS or unsupported-method policy.
- **Interpretations and peer behavior:** Mirror `ServeMux`, synthesize OPTIONS,
  treat every miss as 404, or distinguish known path 405 from unsupported 501.
  Routers disagree on implied HEAD and OPTIONS in `Allow`.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Require explicit uppercase method
  tokens. GET implies HEAD unless explicit HEAD wins. Explicit OPTIONS wins;
  otherwise enabled automation emits 204 and a sorted, deduplicated `Allow`
  including implied HEAD and OPTIONS. A known path with another method is 405;
  a valid method absent from the compiled table is 501 only when no known path
  establishes 405. Disabling automation restores default ServeMux 404/405.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestCompiledRouterPreservesHTTPMethodSemantics`,
  `TestExplicitOptionsAndHeadRoutesWin`,
  `TestDefaultNotFoundAndMethodNotAllowedMatchServeMux`, and
  `TestDocumentedServeMuxDispatchDivergences` cover router options and dispatch.
  Reconsider only through a versioned compatibility change.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-003",
  "title": "Method, HEAD, OPTIONS, Allow, and miss outcomes",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "optional behavior",
  "decision_scope": "transport-specific",
  "specification": "RFC 9110 HTTP Semantics",
  "version": "RFC 9110",
  "source_authority": "rfc9110",
  "section": "Sections 9, 9.3.2, 9.3.7, and 10.2.1",
  "requirement_strength": "not specified",
  "issue": "HTTP defines methods, HEAD, OPTIONS, and Allow semantics without selecting a router's automatic OPTIONS or unsupported-method policy.",
  "interpretations": [
    "Mirror ServeMux misses",
    "Synthesize OPTIONS",
    "Treat every miss as 404",
    "Distinguish known-path 405 from unsupported-method 501"
  ],
  "peer_behavior": "Go 1.26.6 ServeMux is the maintained peer for default method misses and implied HEAD; this package separately assesses its documented OPTIONS and 501 extensions.",
  "selected_behavior": "Require uppercase method tokens; let GET imply HEAD unless explicit HEAD wins; let explicit OPTIONS win; optionally synthesize 204 with sorted Allow; return 405 for a known path and 501 only for a valid method absent from the table when no known path establishes 405.",
  "rationale": "The policy preserves ServeMux defaults unless automation is explicitly enabled and keeps every deliberate divergence observable.",
  "security_consequences": "Invalid methods are rejected before route middleware and Allow discloses only the bounded compiled method set.",
  "resource_consequences": "Method validation and Allow construction operate within configured method and route limits.",
  "compatibility_consequences": "Disabling automation restores default ServeMux outcomes; enabling it is a documented behavioral choice.",
  "wire_consequences": "GET can serve HEAD, automatic OPTIONS returns 204, known-path misses return 405, and unsupported-method misses can return 501.",
  "executable_evidence": [
    "TestCompiledRouterPreservesHTTPMethodSemantics",
    "TestExplicitOptionsAndHeadRoutesWin",
    "TestDefaultNotFoundAndMethodNotAllowedMatchServeMux",
    "TestDocumentedServeMuxDispatchDivergences"
  ],
  "fixture_evidence": [
    "router_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRequestTargets"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "Router.ServeHTTP",
    "WithAutomaticOPTIONS"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "RFC 9110 intentionally leaves router automation policy unspecified.",
  "reconsider_when": "A versioned compatibility change selects different method or OPTIONS behavior."
}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

</details>

## ROUTER-DEC-004: Request-target forms and malformed input

- **Status, owner, and classification:** `resolved`; maintainers; RFC 9110/9112
  transport-boundary policy.
- **Source and issue:** RFC 9112 [request-target](https://www.rfc-editor.org/rfc/rfc9112.html#section-3.2)
  permits origin, absolute, authority, and asterisk forms. Go supplies parsed
  requests, but routing support for CONNECT authority form and `OPTIONS *`
  remains a package choice.
- **Interpretations and peer behavior:** Accept every parsed target, pass all
  forms to ServeMux, reject non-origin forms, or support only the server-wide
  asterisk case. Routers differ on malformed authority and CONNECT.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Accept valid origin and absolute
  forms supplied by `net/http`; support only `OPTIONS *`, returning automatic
  204 or the explicit not-found policy when automation is disabled; reject
  other asterisk use and CONNECT authority-form dispatch with 400. CONNECT
  registration fails at startup. Invalid methods, URLs, and authorities bypass
  route middleware and custom miss handlers.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestAsteriskOptionsAndMalformedAuthority`,
  `TestMalformedRequestsBypassCustomHandlersAndRouteMiddleware`,
  `TestUnsupportedConnectRouteFailsAtStartup`, and `FuzzRequestTargets` cover
  dispatch and registration. Reconsider when a concrete CONNECT service needs
  a separately threat-modeled contract.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-004",
  "title": "Request-target forms and malformed input",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "transport-specific",
  "specification": "RFC 9112 HTTP/1.1 request-target forms",
  "version": "RFC 9112",
  "source_authority": "rfc9112",
  "section": "Section 3.2",
  "requirement_strength": "not specified",
  "issue": "HTTP/1.1 permits origin, absolute, authority, and asterisk request-target forms, while routing support for CONNECT authority form and non-OPTIONS asterisk form remains a package choice.",
  "interpretations": [
    "Accept every parsed target",
    "Pass every form to ServeMux",
    "Reject every non-origin form",
    "Support only the server-wide OPTIONS asterisk case"
  ],
  "peer_behavior": "Go 1.26.6 net/http is the maintained peer for parsed request forms; this package deliberately narrows CONNECT and asterisk routing.",
  "selected_behavior": "Accept parsed origin and absolute forms, support only OPTIONS asterisk with automatic 204 or the configured not-found policy, reject other asterisk and CONNECT authority-form dispatch with 400, and reject CONNECT registration at startup.",
  "rationale": "The router supports ordinary routing forms without claiming ownership of CONNECT tunneling or server-wide protocol behavior beyond OPTIONS.",
  "security_consequences": "Invalid methods, URLs, and authorities bypass custom miss handlers and route middleware.",
  "resource_consequences": "Request targets and authorities are bounded before matching.",
  "compatibility_consequences": "Origin and absolute requests interoperate with net/http while unsupported forms fail deterministically.",
  "wire_consequences": "Unsupported targets return 400 and automatic OPTIONS asterisk returns 204 when enabled.",
  "executable_evidence": [
    "TestAsteriskOptionsAndMalformedAuthority",
    "TestMalformedRequestsBypassCustomHandlersAndRouteMiddleware",
    "TestUnsupportedConnectRouteFailsAtStartup"
  ],
  "fixture_evidence": [
    "compatibility_security_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRequestTargets"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "Builder.Register",
    "Router.ServeHTTP"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "RFC 9112 defines request-target forms but not this router's supported subset.",
  "reconsider_when": "A concrete CONNECT service receives a separately threat-modeled contract."
}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9112.txt

</details>

## ROUTER-DEC-005: Canonical redirects and escaped path structure

- **Status, owner, and classification:** `resolved`; maintainers; Go-compatible
  redirect behavior with explicit defensive override.
- **Source and issue:** Go 1.26.6 ServeMux canonicalizes paths and subtree roots;
  RFC 3986 [path syntax](https://www.rfc-editor.org/rfc/rfc3986.html#section-3.3)
  distinguishes separators from percent-encoded data. Rejecting redirects by
  decoded path can misclassify encoded slash or dot text.
- **Interpretations and peer behavior:** Always follow ServeMux redirects,
  disable all canonicalization, clean decoded paths, or classify structural
  changes using escaped paths. Routers disagree on `%2F` and trailing slash.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Follow ServeMux canonical and subtree
  redirects by default before route/method miss selection. `RejectRedirects`
  converts structural redirects to 404 using escaped-path semantics and
  standard patterns. Encoded separators and dot text inside a wildcard remain
  data; literal and percent-encoded dot segments in registered patterns are
  rejected.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestCanonicalRedirectsPrecedeRouteAndMethodSelection`,
  `TestRejectRedirectPolicyTreatsEncodedSeparatorsAsWildcardData`, and
  `TestRejectRedirectPolicyRejectsSemanticSubtreeRoots` cover `RedirectPolicy`.
  Reconsider when Go canonicalization behavior changes.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-005",
  "title": "Canonical redirects and escaped path structure",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "interoperability policy",
  "decision_scope": "defensive",
  "specification": "Go 1.26.6 net/http and net/url contracts",
  "version": "Go 1.26.6",
  "source_authority": "go-net-http",
  "section": "net/http ServeMux path canonicalization and subtree redirects",
  "requirement_strength": "not specified",
  "issue": "ServeMux canonicalizes paths and subtree roots, while a defensive redirect policy must avoid treating percent-encoded separators or dot text as decoded structure.",
  "interpretations": [
    "Always follow ServeMux redirects",
    "Disable canonicalization",
    "Clean decoded paths",
    "Classify structural changes using escaped paths"
  ],
  "peer_behavior": "Go 1.26.6 ServeMux is the maintained differential peer for canonical redirects, escaped slashes, dot segments, and subtree roots.",
  "selected_behavior": "Follow ServeMux canonical and subtree redirects by default before miss selection; when RejectRedirects is selected, convert structural redirects to 404 using escaped-path semantics while keeping encoded separators and dot text inside wildcards as data.",
  "rationale": "Default delegation preserves Go behavior and escaped-path classification prevents encoded data from becoming false structure.",
  "security_consequences": "Literal and percent-encoded dot segments in registered patterns are rejected and encoded separators cannot bypass structural checks.",
  "resource_consequences": "Redirect classification uses bounded request and pattern data.",
  "compatibility_consequences": "FollowRedirects remains ServeMux-compatible; RejectRedirects is an explicit defensive divergence.",
  "wire_consequences": "Canonical requests redirect by default or return 404 under the explicit rejection policy.",
  "executable_evidence": [
    "TestCanonicalRedirectsPrecedeRouteAndMethodSelection",
    "TestRejectRedirectPolicyTreatsEncodedSeparatorsAsWildcardData",
    "TestRejectRedirectPolicyRejectsSemanticSubtreeRoots"
  ],
  "fixture_evidence": [
    "compatibility_security_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRequestTargets"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "RedirectPolicy",
    "Router.ServeHTTP"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "Go 1.26.6 ServeMux is the redirect-behavior authority.",
  "reconsider_when": "Go canonicalization behavior changes."
}
```

Authority URL: https://go.dev/src/net/http/server.go?m=text

</details>

## ROUTER-DEC-006: Host patterns, ports, wildcard labels, and IDNA

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  authority policy and documented ServeMux extension.
- **Source and issue:** RFC 9110 [authority](https://www.rfc-editor.org/rfc/rfc9110.html#section-7.2)
  and RFC 3986 [host](https://www.rfc-editor.org/rfc/rfc3986.html#section-3.2.2)
  allow forms unsuitable for implicit tenant routing. ServeMux supports literal
  hosts but not this package's single-label wildcard precedence.
- **Interpretations and peer behavior:** Match raw Host, normalize Unicode,
  include ports, allow arbitrary wildcard suffixes, or require validated ASCII
  DNS labels. Host-routing libraries vary on proxy and IDNA trust.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Route patterns accept bounded ASCII
  DNS labels and one-label `{name}` wildcards; ports and IP literals are not
  patterns. Request ports are removed for matching; exact hosts precede wildcard
  hosts, then hostless fallback. Malformed, non-ASCII, user-info, and invalid
  port authorities fail 400. No forwarding header or implicit IDNA conversion
  is used; callers normalize trusted IDNA at their boundary.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestHostPatternsMatchPortsAndSingleLabels`,
  `TestHostSpecificityFallbackAndEquivalentPatterns`,
  `TestAmbiguousHostPatternsAndUnsafeAuthoritiesAreRejected`, and
  `FuzzHostMatching` cover host registration, dispatch, and generation.
  Reconsider for a versioned IDNA policy or multi-label wildcard design.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-006",
  "title": "Host patterns, ports, wildcard labels, and IDNA",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "defensive",
  "specification": "RFC 9110 HTTP Semantics",
  "version": "RFC 9110",
  "source_authority": "rfc9110",
  "section": "Section 7.2 authority handling",
  "requirement_strength": "not specified",
  "issue": "HTTP authority syntax permits forms that are unsafe for implicit tenant routing, and ServeMux does not define this package's wildcard-host precedence or IDNA trust policy.",
  "interpretations": [
    "Match raw Host",
    "Normalize Unicode implicitly",
    "Include ports in route patterns",
    "Allow arbitrary wildcard suffixes",
    "Require validated ASCII DNS labels"
  ],
  "peer_behavior": "Go 1.26.6 ServeMux is the maintained peer for literal host matching; the single-label wildcard and defensive authority policy are assessed as deliberate differences.",
  "selected_behavior": "Accept bounded ASCII DNS labels and one-label named wildcards, remove request ports for matching, prefer exact hosts then wildcard hosts then hostless routes, reject malformed or non-ASCII authorities, and perform no forwarding-header or implicit IDNA conversion.",
  "rationale": "Explicit ASCII authority policy avoids hidden proxy and normalization trust while preserving deterministic host precedence.",
  "security_consequences": "User-info, malformed ports, non-ASCII hosts, IP pattern ambiguity, and implicit forwarding trust are rejected.",
  "resource_consequences": "Authority and label lengths are bounded before matching or generation.",
  "compatibility_consequences": "Literal host matching remains ServeMux-compatible and single-label wildcards are an explicit extension.",
  "wire_consequences": "Request ports do not affect host matching and unsafe authorities return 400.",
  "executable_evidence": [
    "TestHostPatternsMatchPortsAndSingleLabels",
    "TestHostSpecificityFallbackAndEquivalentPatterns",
    "TestAmbiguousHostPatternsAndUnsafeAuthoritiesAreRejected"
  ],
  "fixture_evidence": [
    "compatibility_security_test.go"
  ],
  "fuzz_evidence": [
    "FuzzHostMatching"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "Route.Host",
    "Router.ServeHTTP",
    "Router.URL"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "RFC 9110 defines authority semantics but not application host-routing trust.",
  "reconsider_when": "A versioned IDNA policy or multi-label wildcard design is adopted."
}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

</details>

## ROUTER-DEC-007: Group composition and transactional callbacks

- **Status, owner, and classification:** `resolved`; maintainers; application
  composition policy.
- **Source and issue:** Go's
  [`net/http.Handler`](https://cs.opensource.google/go/go/+/refs/tags/go1.26.6:src/net/http/server.go)
  and HTTP do not define route groups, prefix joining, metadata merge, naming,
  or callback failure semantics.
- **Interpretations and peer behavior:** Mutate a shared group, clean paths,
  silently override metadata, retain partial callback routes, or flatten one
  validated transaction. Framework groups commonly hide mutation and merge
  precedence.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Flatten nested host, path, name,
  metadata, and middleware state at registration. Join prefixes without
  `path.Clean`; reject empty/dot/wildcard-invalid segments; require equal
  repeated hosts; reject metadata collisions; and publish no child route when
  callback or validation fails. Group and nesting budgets apply even to empty
  groups.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestNestedGroupsFlattenComposition`,
  `TestGroupCompositionRejectsInvalidAndPartialState`,
  `TestNestedGroupPrefixesUseComposedBudgets`, and `FuzzGroupComposition` cover
  `Builder.Group` and `GroupOptions`. No upstream issue exists; reconsider only
  for a new explicit merge policy.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-007",
  "title": "Group composition and transactional callbacks",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "application-policy",
  "specification": "Go 1.26.6 net/http and net/url contracts",
  "version": "Go 1.26.6",
  "source_authority": "go-net-http",
  "section": "net/http Handler composition",
  "requirement_strength": "not specified",
  "issue": "Go handlers and HTTP do not define route groups, prefix joining, metadata merging, naming, or callback failure publication.",
  "interpretations": [
    "Mutate a shared group",
    "Clean joined paths",
    "Silently override metadata",
    "Retain partial callback routes",
    "Flatten one validated transaction"
  ],
  "peer_behavior": "No maintained peer comparison is currently assessed for transactional group composition.",
  "selected_behavior": "Flatten nested host, path, name, metadata, and middleware state at registration; join without path cleaning; reject invalid segments, differing repeated hosts, and metadata collisions; publish no child route when callback or validation fails.",
  "rationale": "Transactional flattening makes the compiled route table deterministic and prevents hidden partial state.",
  "security_consequences": "Invalid composition cannot silently change paths, hosts, or metadata ownership.",
  "resource_consequences": "Group, nesting, route, prefix, name, metadata, and middleware budgets apply before publication, including empty groups.",
  "compatibility_consequences": "Composition and failure behavior remain deterministic and independent of callback partial progress.",
  "wire_consequences": "Only fully validated grouped routes can dispatch.",
  "executable_evidence": [
    "TestNestedGroupsFlattenComposition",
    "TestGroupCompositionRejectsInvalidAndPartialState",
    "TestNestedGroupPrefixesUseComposedBudgets"
  ],
  "fixture_evidence": [
    "group_test.go"
  ],
  "fuzz_evidence": [
    "FuzzGroupComposition"
  ],
  "interoperability_evidence": [],
  "public_apis": [
    "Builder.Group",
    "GroupOptions"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "No upstream route-group contract exists.",
  "reconsider_when": "A new explicit group merge policy is versioned."
}
```

Authority URL: https://go.dev/src/net/http/server.go?m=text

</details>

## ROUTER-DEC-008: Middleware order, exclusions, and runtime ownership

- **Status, owner, and classification:** `resolved`; maintainers; explicit
  composition and lifecycle policy.
- **Source and issue:** Go's
  [`http.Handler`](https://cs.opensource.google/go/go/+/refs/tags/go1.26.6:src/net/http/server.go)
  permits decorators but defines no router, group, mount, route order,
  exclusion, duplicate-name, panic, or recovery behavior.
- **Interpretations and peer behavior:** Resolve aliases globally, instantiate
  by reflection, apply route-first, silently deduplicate, or freeze explicit
  values at compile. Frameworks differ and often own recovery implicitly.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Execute router, outer group, inner
  group, then route middleware; unwind in reverse. Named exclusions remove only
  inherited layers and are explicit on the route. Nil, duplicate resolved names,
  and nil constructed handlers fail before publication. Serving panics,
  cancellation, short circuits, re-entry, writer capabilities, and handler
  lifecycle remain caller-owned; no recovery or wrapper is injected.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestMiddlewareOrderAndIntrospectionAreStableAndImmutable`,
  `TestRouteMayExcludeNamedGroupMiddleware`,
  `TestMiddlewareMayShortCircuitPanicCancelAndReenter`, and
  `TestRouterPreservesResponseWriterOptionalInterfaces` cover middleware APIs.
  Reconsider only through explicit versioned ordering policy.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-008",
  "title": "Middleware order, exclusions, and runtime ownership",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "application-policy",
  "specification": "Go 1.26.6 net/http and net/url contracts",
  "version": "Go 1.26.6",
  "source_authority": "go-net-http",
  "section": "net/http Handler and ResponseWriter contracts",
  "requirement_strength": "not specified",
  "issue": "Go handlers define decorator composition but not router, group, mount, and route order, exclusions, duplicate names, recovery, or middleware construction failure.",
  "interpretations": [
    "Resolve aliases globally",
    "Instantiate by reflection",
    "Apply route middleware first",
    "Silently deduplicate middleware",
    "Freeze explicit values during compilation"
  ],
  "peer_behavior": "No maintained peer comparison is currently assessed for the complete middleware ownership contract.",
  "selected_behavior": "Execute router, outer group, inner group, then route middleware and unwind in reverse; named exclusions remove inherited layers only; fail compilation on nil or duplicate resolved middleware; leave serving panics, cancellation, short circuits, re-entry, writer capabilities, and handler lifetime caller-owned.",
  "rationale": "Explicit deterministic composition preserves ordinary net/http ownership without hidden recovery or service location.",
  "security_consequences": "No implicit recovery, reflection, or alias lookup can hide faults or select unexpected middleware.",
  "resource_consequences": "Middleware counts and identifiers are bounded before handler construction.",
  "compatibility_consequences": "Order, exclusion, panic, cancellation, and writer behavior are stable public contracts.",
  "wire_consequences": "Middleware response order is the reverse of request order and short circuits remain middleware-owned.",
  "executable_evidence": [
    "TestMiddlewareOrderAndIntrospectionAreStableAndImmutable",
    "TestRouteMayExcludeNamedGroupMiddleware",
    "TestMiddlewareMayShortCircuitPanicCancelAndReenter",
    "TestRouterPreservesResponseWriterOptionalInterfaces"
  ],
  "fixture_evidence": [
    "middleware_hardening_test.go"
  ],
  "fuzz_evidence": [
    "FuzzGroupComposition"
  ],
  "interoperability_evidence": [],
  "public_apis": [
    "Middleware",
    "NamedMiddleware",
    "Route.ExcludeMiddleware",
    "Router.ServeHTTP"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "Go deliberately leaves handler composition policy to applications.",
  "reconsider_when": "A versioned middleware ordering or ownership policy is adopted."
}
```

Authority URL: https://go.dev/src/net/http/server.go?m=text

</details>

## ROUTER-DEC-009: Mount stripping and request identity

- **Status, owner, and classification:** `resolved`; maintainers; Go request
  interoperability and defensive path policy.
- **Source and issue:** Go's
  [`StripPrefix`](https://cs.opensource.google/go/go/+/refs/tags/go1.26.6:src/net/http/server.go)
  and HTTP do not define this package's mounted-handler prefix stripping,
  `RawPath`, `RequestURI`, inherited path values, or mutation ownership.
- **Interpretations and peer behavior:** Mutate the original request, use
  `http.StripPrefix`, discard escaped form, copy routers internally, or clone
  only the URL view supplied to the mounted handler.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Mount one ordinary remainder route.
  Optional stripping clones request and URL, compares encoded literal prefixes
  in decoded path space, preserves the escaped suffix in `RawPath`, and leaves
  caller URL and `RequestURI` untouched. A compiled router mounts as an ordinary
  handler. Outer path values survive unless an inner wildcard reuses the name.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestMountStripsPathOnCloneAndPreservesRequestTarget`,
  `TestMountStripsEscapedLiteralPrefixWithoutLosingRawPath`, and
  `TestCompiledRouterMountIsAnOrdinaryHandler` cover `Mount` and `MountOptions`.
  Reconsider if Go publishes richer mount semantics.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-009",
  "title": "Mount stripping and request identity",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "application-policy",
  "specification": "Go 1.26.6 net/http and net/url contracts",
  "version": "Go 1.26.6",
  "source_authority": "go-net-http",
  "section": "net/http StripPrefix and Request URL ownership",
  "requirement_strength": "not specified",
  "issue": "Go does not define this package's mounted-handler prefix stripping, RawPath preservation, RequestURI ownership, inherited path values, or nested-router behavior.",
  "interpretations": [
    "Mutate the original request",
    "Use http.StripPrefix directly",
    "Discard escaped form",
    "Copy routers internally",
    "Clone only the URL view supplied to the mounted handler"
  ],
  "peer_behavior": "Go 1.26.6 http.StripPrefix is the maintained peer; this package preserves caller request identity and escaped suffixes as deliberate differences.",
  "selected_behavior": "Mount one ordinary remainder route; optional stripping clones the request and URL, compares encoded literal prefixes in decoded path space, preserves the escaped suffix in RawPath, leaves caller URL and RequestURI untouched, and preserves outer path values unless an inner wildcard reuses the name.",
  "rationale": "Cloned URL mutation preserves caller request identity while supplying an accurate mounted-handler view.",
  "security_consequences": "Encoded literal prefixes cannot be confused with structural separators during mount stripping.",
  "resource_consequences": "Mount prefixes and cloned URL state remain within configured limits.",
  "compatibility_consequences": "Compiled routers remain ordinary handlers and nested path-value collisions have deterministic ownership.",
  "wire_consequences": "Mounted handlers observe a stripped URL without changes to the caller's original RequestURI.",
  "executable_evidence": [
    "TestMountStripsPathOnCloneAndPreservesRequestTarget",
    "TestMountStripsEscapedLiteralPrefixWithoutLosingRawPath",
    "TestCompiledRouterMountIsAnOrdinaryHandler"
  ],
  "fixture_evidence": [
    "mount_policy_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRequestTargets"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "Builder.Mount",
    "MountOptions"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "Go exposes StripPrefix but no router mount contract.",
  "reconsider_when": "Go publishes richer mount and path-value semantics."
}
```

Authority URL: https://go.dev/src/net/http/server.go?m=text

</details>

## ROUTER-DEC-010: Path values, matched metadata, and introspection

- **Status, owner, and classification:** `resolved`; maintainers; Go-compatible
  parameter behavior plus bounded application metadata policy.
- **Source and issue:** Go 1.26.6
  [`Request.PathValue`](https://cs.opensource.google/go/go/+/refs/tags/go1.26.6:src/net/http/request.go)
  owns route parameters but does not expose a public immutable route table or
  matched-route descriptor.
  A parallel parameter API risks divergence; raw metadata risks disclosure and
  telemetry cardinality.
- **Interpretations and peer behavior:** Copy parameters into a custom context,
  expose handler pointers, return internal maps, or retain standard path values
  and add a defensive metadata snapshot.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Handlers use `Request.PathValue`
  directly. `MatchedRoute` adds a copied bounded `RouteInfo` only inside the
  selected chain; route tables are deterministically sorted and copied, omit
  handlers/function names, and retain bounded caller metadata. Misses install
  neither route metadata nor path values.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestCompiledRouterDispatchesWithPathValuesAndMatchedRoute`,
  `TestMiddlewareOrderAndIntrospectionAreStableAndImmutable`, and
  `TestConcurrentDispatchIntrospectionAndGeneration` cover `MatchedRoute` and
  `Routes`. Reconsider if Go exposes equivalent route metadata.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-010",
  "title": "Path values, matched metadata, and introspection",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "application-policy",
  "specification": "Go 1.26.6 net/http and net/url contracts",
  "version": "Go 1.26.6",
  "source_authority": "go-net-http-request",
  "section": "net/http Request.PathValue and ServeMux dispatch",
  "requirement_strength": "not specified",
  "issue": "Go owns request path values but does not expose immutable route introspection or matched-route metadata.",
  "interpretations": [
    "Copy parameters into a custom context",
    "Expose handler pointers",
    "Return internal metadata maps",
    "Keep standard path values and add a bounded defensive metadata snapshot"
  ],
  "peer_behavior": "Go 1.26.6 Request.PathValue is the maintained peer for parameters; bounded matched-route metadata is a deliberate package extension.",
  "selected_behavior": "Use Request.PathValue directly; expose a copied bounded RouteInfo only inside the selected chain; sort and copy route tables deterministically; omit handlers and function names; install neither route metadata nor path values on misses.",
  "rationale": "Standard path values avoid a divergent parameter API and defensive copies prevent mutation or handler disclosure.",
  "security_consequences": "Introspection omits handler identity and bounds caller metadata to control disclosure and cardinality.",
  "resource_consequences": "Route metadata and introspection copies obey finite table and metadata limits.",
  "compatibility_consequences": "Handlers retain Go path-value semantics while route metadata has explicit immutable ownership.",
  "wire_consequences": "Only selected handlers observe matched metadata; misses do not receive route state.",
  "executable_evidence": [
    "TestCompiledRouterDispatchesWithPathValuesAndMatchedRoute",
    "TestMiddlewareOrderAndIntrospectionAreStableAndImmutable",
    "TestConcurrentDispatchIntrospectionAndGeneration"
  ],
  "fixture_evidence": [
    "router_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRoutePatternCompilation"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "MatchedRoute",
    "RouteInfo",
    "Router.Routes",
    "Request.PathValue"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "Go exposes PathValue but no equivalent route metadata API.",
  "reconsider_when": "Go exposes equivalent immutable route metadata."
}
```

Authority URL: https://go.dev/src/net/http/request.go?m=text

</details>

## ROUTER-DEC-011: Named path generation and percent encoding

- **Status, owner, and classification:** `resolved`; maintainers; RFC 3986
  component encoding with Go URL interoperability.
- **Source and issue:** RFC 3986 [percent encoding](https://www.rfc-editor.org/rfc/rfc3986.html#section-2.1)
  and [path segments](https://www.rfc-editor.org/rfc/rfc3986.html#section-3.3)
  permit many strings, but interpolating raw values can inject separators,
  traversal, or double decoding.
- **Interpretations and peer behavior:** Raw substitution, escape the entire
  path, stringify arbitrary models, accept one raw remainder, or require typed
  explicit segment inputs. Reverse routers differ on extra parameters.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Require every wildcard exactly once;
  reject missing, duplicate, unknown, unused, or wrong-kind parameters. Escape
  each segment once with `url.PathEscape`; a remainder is an explicit non-empty
  segment list and rejects empty/dot segments. Query follows bounded
  deterministic `url.Values.Encode`. Generated paths must round-trip to the
  intended ServeMux path values; no reflection or model binding occurs.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestNamedPathGenerationEscapesSegmentsAndRoundTrips`,
  `TestRemainderGenerationRequiresExplicitSafeSegments`,
  `TestGenerationRejectsParameterSetErrors`, and `FuzzNamedRouteRoundTrip`
  cover `Param`, `Remainder`, and `Router.Path`. Reconsider for new typed
  component kinds, not raw interpolation.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-011",
  "title": "Named path generation and percent encoding",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "defensive",
  "specification": "RFC 3986 URI Generic Syntax",
  "version": "RFC 3986",
  "source_authority": "rfc3986",
  "section": "Sections 2.1 and 3.3",
  "requirement_strength": "not specified",
  "issue": "URI syntax permits many strings, but a named-route generator must select parameter ownership, escaping, remainder, query, and error behavior.",
  "interpretations": [
    "Substitute raw values",
    "Escape the complete path",
    "Stringify arbitrary models",
    "Accept one raw remainder",
    "Require typed explicit segment inputs"
  ],
  "peer_behavior": "Go 1.26.6 net/url and ServeMux are the maintained peers for component escaping and generated-path round trips.",
  "selected_behavior": "Require every wildcard exactly once; reject missing, duplicate, unknown, unused, or wrong-kind parameters; escape each segment once with url.PathEscape; require explicit non-empty safe remainder segments; encode bounded queries deterministically with url.Values.Encode.",
  "rationale": "Component-level typed generation preserves URI structure and makes generated values round-trip through ServeMux.",
  "security_consequences": "Segment, traversal, separator, and double-decoding injection are rejected or escaped as data.",
  "resource_consequences": "Parameter counts, segment counts, query data, and generated output are bounded before expensive work.",
  "compatibility_consequences": "Generated routes use deterministic typed parameters rather than reflection or model binding.",
  "wire_consequences": "Generated paths and queries are percent-encoded once and dispatch to the intended path values.",
  "executable_evidence": [
    "TestNamedPathGenerationEscapesSegmentsAndRoundTrips",
    "TestRemainderGenerationRequiresExplicitSafeSegments",
    "TestGenerationRejectsParameterSetErrors"
  ],
  "fixture_evidence": [
    "generation_test.go"
  ],
  "fuzz_evidence": [
    "FuzzNamedRouteRoundTrip"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "Param",
    "Remainder",
    "Router.Path"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "RFC 3986 defines URI components but no named-route API.",
  "reconsider_when": "New typed component kinds are introduced."
}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc3986.txt

</details>

## ROUTER-DEC-012: Absolute URL bases and route hosts

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  origin and host trust policy.
- **Source and issue:** RFC 3986 [authority](https://www.rfc-editor.org/rfc/rfc3986.html#section-3.2)
  allows userinfo and broad host syntax. Inferring a base from request or proxy
  fields can produce open redirects or cross-tenant URLs.
- **Interpretations and peer behavior:** Infer request scheme/host, trust
  forwarding fields, accept arbitrary schemes, or require one validated
  immutable base. Framework route helpers commonly depend on ambient requests.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  `NewBaseURL` accepts only explicit
  bounded HTTP(S) scheme and validated authority without userinfo, controls, or
  malformed ports. Generation never reads a request or forwarding field. A
  literal or rendered route host replaces the base hostname while preserving
  its explicit trusted port. Unicode hosts require caller-selected ASCII IDNA.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestAbsoluteURLGenerationValidatesBaseHostAndQuery`,
  `TestGenerationEnforcesOutputAndQueryLimits`,
  `TestTrustedBaseAuthorityIsBounded`, and `FuzzURLGenerationInputs` cover
  `BaseURL`, `NewBaseURL`, and `Router.URL`. Reconsider only with an explicit
  trusted-proxy integration contract.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-012",
  "title": "Absolute URL bases and route hosts",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "defensive",
  "specification": "RFC 3986 URI Generic Syntax",
  "version": "RFC 3986",
  "source_authority": "rfc3986",
  "section": "Section 3.2",
  "requirement_strength": "not specified",
  "issue": "URI authority permits userinfo and broad host syntax, while ambient request and proxy fields can produce untrusted absolute URLs.",
  "interpretations": [
    "Infer scheme and host from requests",
    "Trust forwarding fields",
    "Accept arbitrary schemes",
    "Require one explicit validated immutable base"
  ],
  "peer_behavior": "Go 1.26.6 net/url is the maintained peer for URL construction; the package deliberately requires an explicit trusted base.",
  "selected_behavior": "Accept only an explicit bounded HTTP or HTTPS base with validated authority and no userinfo, controls, or malformed ports; never read request or forwarding fields; replace the base hostname with a rendered route host while preserving its trusted explicit port.",
  "rationale": "Explicit trusted bases prevent ambient authority from controlling generated origins.",
  "security_consequences": "Open redirects, credential-bearing authorities, malformed ports, controls, and implicit proxy trust are rejected.",
  "resource_consequences": "Schemes, authorities, query data, parameters, and output length are bounded.",
  "compatibility_consequences": "Callers provide stable explicit origins and select any IDNA normalization before construction.",
  "wire_consequences": "Generated absolute URLs use only the validated base and rendered route host.",
  "executable_evidence": [
    "TestAbsoluteURLGenerationValidatesBaseHostAndQuery",
    "TestGenerationEnforcesOutputAndQueryLimits",
    "TestTrustedBaseAuthorityIsBounded"
  ],
  "fixture_evidence": [
    "generation_test.go"
  ],
  "fuzz_evidence": [
    "FuzzURLGenerationInputs"
  ],
  "interoperability_evidence": [
    "specification/interoperability.tsv"
  ],
  "public_apis": [
    "BaseURL",
    "NewBaseURL",
    "Router.URL"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "RFC 3986 defines authority syntax but not trusted-origin selection.",
  "reconsider_when": "An explicit trusted-proxy integration contract is introduced."
}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc3986.txt

</details>

## ROUTER-DEC-013: Finite limits, errors, and custom miss handlers

- **Status, owner, and classification:** `resolved`; maintainers; defensive
  resource, disclosure, and handler-ownership policy.
- **Source and issue:** RFC 9110
  [`414 URI Too Long`](https://www.rfc-editor.org/rfc/rfc9110.html#section-15.5.15)
  and Go's `ServeMux` do not bound route tables, metadata, patterns, targets,
  diagnostics, middleware, parameters, or generated output, and do not define
  custom miss-handler panic/partial-write behavior.
- **Interpretations and peer behavior:** Leave all input unbounded, truncate
  semantic values, expose route inventories in errors, recover custom handlers,
  or reject at exact finite boundaries with safe diagnostics.
- **Selected behavior, security and resource consequences, compatibility and wire consequences:**
  Validate positive immutable limits
  and reject syntactically valid excess with `ErrLimitExceeded`; reject runtime
  target excess with 414 before matching. Errors are typed, deterministic,
  bounded, single-line valid UTF-8, and omit route inventories and values.
  Custom 404/405 handlers own cancellation, partial responses, and panics;
  malformed input bypasses them. No hidden I/O, goroutine, global registration,
  or mutable compiled state exists.
- **Evidence, public surface, upstream, and reconsideration:**
  `TestRegisterEnforcesLimitsAndBoundsDiagnostics`,
  `TestFineGrainedInputByteBudgets`,
  `TestDiagnosticsAreSingleLineValidUTF8`,
  `TestCustomErrorHandlerOwnsPartialResponsesAndPanics`, and
  `TestEveryLimitMustBePositive` cover `Limits`, `Error`, and custom options.
  Reconsider limits only with measured compatibility and resource evidence.

<details>
<summary>Machine-auditable decision record</summary>

```json
{
  "id": "ROUTER-DEC-013",
  "title": "Finite limits, errors, and custom miss handlers",
  "status": "resolved",
  "owner": "router maintainers",
  "classification": "omission",
  "decision_scope": "defensive",
  "specification": "RFC 9110 HTTP Semantics",
  "version": "RFC 9110",
  "source_authority": "rfc9110",
  "section": "Section 15.5.15 and application resource policy",
  "requirement_strength": "not specified",
  "issue": "HTTP and ServeMux do not bound router inputs, diagnostics, metadata, target length, or generated output and do not define custom miss-handler ownership after panic or partial write.",
  "interpretations": [
    "Leave every input unbounded",
    "Truncate semantic values",
    "Expose route inventories in errors",
    "Recover custom handlers",
    "Reject at exact finite boundaries with safe diagnostics"
  ],
  "peer_behavior": "No maintained peer comparison is currently assessed for the complete limits and custom-handler contract.",
  "selected_behavior": "Require positive immutable limits; reject syntactically valid excess with ErrLimitExceeded; return 414 for runtime target excess before matching; produce typed deterministic bounded single-line UTF-8 errors; leave cancellation, partial responses, and panics to custom handlers; perform no hidden I/O, goroutines, or global registration.",
  "rationale": "Finite rejection boundaries and caller-owned handlers make resource and failure behavior explicit without semantic truncation.",
  "security_consequences": "Errors omit route inventories and untrusted values, malformed input bypasses custom handlers, and excessive input is rejected before amplification.",
  "resource_consequences": "Every configurable and runtime dimension has a positive finite bound enforced before allocation-heavy work.",
  "compatibility_consequences": "Exact boundary outcomes, error categories, and custom-handler ownership remain stable.",
  "wire_consequences": "Overlong targets return 414; malformed targets return 400; custom 404 and 405 handlers own any partial response or panic.",
  "executable_evidence": [
    "TestRegisterEnforcesLimitsAndBoundsDiagnostics",
    "TestFineGrainedInputByteBudgets",
    "TestDiagnosticsAreSingleLineValidUTF8",
    "TestCustomErrorHandlerOwnsPartialResponsesAndPanics",
    "TestEveryLimitMustBePositive"
  ],
  "fixture_evidence": [
    "builder_test.go"
  ],
  "fuzz_evidence": [
    "FuzzRoutePatternCompilation",
    "FuzzURLGenerationInputs"
  ],
  "interoperability_evidence": [],
  "public_apis": [
    "Limits",
    "Error",
    "WithNotFound",
    "WithMethodNotAllowed"
  ],
  "documentation": [
    "docs/specification-decisions.md"
  ],
  "upstream_status": "RFC 9110 defines 414 but leaves router resource and custom-handler policy unspecified.",
  "reconsider_when": "Measured compatibility or resource evidence justifies a versioned limit change."
}
```

Authority URL: https://www.rfc-editor.org/rfc/rfc9110.txt

</details>

## Unresolved and excluded behavior

No known material ambiguity in the current public surface is unresolved.
Dynamic routes, regex matching, CONNECT tunneling, proxy trust, implicit IDNA,
controller resolution, model binding, sessions, CSRF, templates, dependency
injection, authentication, authorization, RPC dispatch, and server lifecycle
are outside the v1 claim. Adding one requires a new decision before runtime
implementation.
