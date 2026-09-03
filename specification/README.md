# Specification conformance matrix

`manifest.tsv` pins the normative sources used by the router decision register.
RFC sources use immutable RFC Editor text. The matcher and request contract use
the official Go 1.26.6 source archive. SHA-256 digests make source drift
explicit.

The module claims only its documented routing and URL-generation surface. HTTP
framing remains delegated to Go's server, and deliberate ServeMux differences
remain visible in the decision register and differential tests.

The canonical
[`docs/specification-decisions.md`](../docs/specification-decisions.md)
records every material interpretation, consequence, and condition for
reconsideration behind this conformance matrix.

## Decision conformance matrix

| Decision | Authority | Executable evidence | Differential evidence |
| --- | --- | --- | --- |
| ROUTER-DEC-001 | Go 1.26.6 `http.ServeMux` | `TestSupportedMatchingIsDifferentialWithServeMux`, `TestSupportedMethodsAndLiteralHostsAreDifferentialWithServeMux`, `TestPinnedServeMuxRedirectsPreserveGo126EscapedPathBehavior` | `specification/interoperability.tsv` |
| ROUTER-DEC-002 | Go 1.26.6 `http.ServeMux` registration | `TestCompileReturnsTypedConflictsAndFreezesOnlyOnSuccess`, `TestConflictFailureDoesNotConstructAnyMiddleware`, `TestRegistrationOrderDoesNotChangeDispatchOrIntrospection`, `TestPatternValidationPropagatesUncontrolledPanics` | `specification/interoperability.tsv` |
| ROUTER-DEC-003 | RFC 9110 methods and `Allow` | `TestCompiledRouterPreservesHTTPMethodSemantics`, `TestExplicitOptionsAndHeadRoutesWin`, `TestDefaultNotFoundAndMethodNotAllowedMatchServeMux`, `TestDocumentedServeMuxDispatchDivergences` | `specification/interoperability.tsv` |
| ROUTER-DEC-004 | RFC 9112 request-target forms | `TestAsteriskOptionsAndMalformedAuthority`, `TestMalformedRequestsBypassCustomHandlersAndRouteMiddleware`, `TestUnsupportedConnectRouteFailsAtStartup` | `specification/interoperability.tsv` |
| ROUTER-DEC-005 | Go 1.26.6 `http.ServeMux` redirects | `TestCanonicalRedirectsPrecedeRouteAndMethodSelection`, `TestPinnedServeMuxRedirectsPreserveGo126EscapedPathBehavior`, `TestPinnedSubtreeRedirectDoesNotOverrideExplicitRoot`, `TestRejectRedirectPolicyTreatsEncodedSeparatorsAsWildcardData`, `TestRejectRedirectPolicyRejectsSemanticSubtreeRoots` | `specification/interoperability.tsv` |
| ROUTER-DEC-006 | RFC 9110 authority handling | `TestHostPatternsMatchPortsAndSingleLabels`, `TestHostSpecificityFallbackAndEquivalentPatterns`, `TestAmbiguousHostPatternsAndUnsafeAuthoritiesAreRejected` | `specification/interoperability.tsv` |
| ROUTER-DEC-007 | Go 1.26.6 handler composition | `TestNestedGroupsFlattenComposition`, `TestGroupCompositionRejectsInvalidAndPartialState`, `TestNestedGroupPrefixesUseComposedBudgets` | Not assessed |
| ROUTER-DEC-008 | Go 1.26.6 handlers and response writers | `TestMiddlewareOrderAndIntrospectionAreStableAndImmutable`, `TestRouteMayExcludeNamedGroupMiddleware`, `TestMiddlewareMayShortCircuitPanicCancelAndReenter`, `TestRouterPreservesResponseWriterOptionalInterfaces` | Not assessed |
| ROUTER-DEC-009 | Go 1.26.6 request and URL ownership | `TestMountStripsPathOnCloneAndPreservesRequestTarget`, `TestMountStripsEscapedLiteralPrefixWithoutLosingRawPath`, `TestCompiledRouterMountIsAnOrdinaryHandler` | `specification/interoperability.tsv` |
| ROUTER-DEC-010 | Go 1.26.6 `Request.PathValue` | `TestCompiledRouterDispatchesWithPathValuesAndMatchedRoute`, `TestMiddlewareOrderAndIntrospectionAreStableAndImmutable`, `TestConcurrentDispatchIntrospectionAndGeneration` | `specification/interoperability.tsv` |
| ROUTER-DEC-011 | RFC 3986 path segments and encoding | `TestNamedPathGenerationEscapesSegmentsAndRoundTrips`, `TestRemainderGenerationRequiresExplicitSafeSegments`, `TestGenerationRejectsParameterSetErrors` | `specification/interoperability.tsv` |
| ROUTER-DEC-012 | RFC 3986 authority syntax | `TestAbsoluteURLGenerationValidatesBaseHostAndQuery`, `TestGenerationEnforcesOutputAndQueryLimits`, `TestTrustedBaseAuthorityIsBounded` | `specification/interoperability.tsv` |
| ROUTER-DEC-013 | RFC 9110 status and resource policy | `TestRegisterEnforcesLimitsAndBoundsDiagnostics`, `TestFineGrainedInputByteBudgets`, `TestDiagnosticsAreSingleLineValidUTF8`, `TestCustomErrorHandlerOwnsPartialResponsesAndPanics`, `TestEveryLimitMustBePositive` | Not assessed |

Run the focused map and evidence check with `make conformance`. For an update,
verify provenance and digest, review Go routing changes and RFC errata, update
the structured and human decisions, preserve decision history, update tests,
then change the manifest. A digest change alone MUST NOT silently change
behavior.
