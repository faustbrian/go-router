# Adoption Guides

Adopt `github.com/faustbrian/go-router/v2` at v2.0.0 once its public tag and
module artifacts are available. Until then, existing applications can retain
released v1 without local replacements or pseudo-versions. A prepared source
tree or changelog date alone does not establish publication.

The owned `go-http-middleware/integration/siblings` and
`go-service/integration/reference-http` consumers must migrate and verify
compatibility against the public v2 module before their migration is complete.
Candidate-source integration checks are not public-consumer evidence. Follow
[migration](migration.md) for inherited middleware and mount-default changes.

## REST-like APIs

Use explicit method sets and `{id}` path wildcards. Group stable API versions,
attach generic middleware by concern, publish names only for URLs that are a
SemVer contract, and inspect `Routes` during startup documentation checks.

## JSON-RPC

Create the dispatcher and HTTP adapter in `jsonrpc`, then mount that
`http.Handler` at `/rpc` with `StripPrefix` only when the adapter expects `/`.
JSON-RPC method names remain inside the dispatcher and never become HTTP routes.

## Webhooks

Mount each provider handler at an explicit boundary or register one POST route.
Attach provider authentication as route middleware. Do not infer authentication
from handler type or metadata.

## Health, metrics, and debug

Mount ordinary handlers at explicit paths and hosts. Keep operational endpoints
on a dedicated host or group when policy differs. A mount does not inherit any
authorization beyond middleware visible in its descriptor.

## Mixed services

REST routes, an RPC mount, webhooks, and probes can share one builder. Compile
once, inspect the flattened table, then give the immutable handler to
`service`. Track, Postal, and Location retain ownership of handlers and
middleware lifecycle.
