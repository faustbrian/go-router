#!/usr/bin/env bash
set -euo pipefail

temporary="$(mktemp -d)"
cleanup() {
  exit_code=$?
  trap - EXIT HUP INT TERM
  chmod -R u+w "$temporary" 2>/dev/null || true
  find "$temporary" -depth -delete
  exit "$exit_code"
}
trap cleanup EXIT HUP INT TERM

repository_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# The temporary v2 replacement tests unpublished source; published consumers
# must still be verified without replacements after the first v2 tag.
export GOCACHE="$temporary/gocache"
export GOMODCACHE="$temporary/gomodcache"
export GOPATH="$temporary/gopath"
export GOWORK=off

mkdir "$temporary/integration"
cat >"$temporary/integration/go.mod" <<EOF
module routerintegration

go 1.27.0

require (
  github.com/faustbrian/go-http-middleware v1.0.0
  github.com/faustbrian/go-jsonrpc v1.0.0
  github.com/faustbrian/go-router v1.0.0
  github.com/faustbrian/go-router/v2 v2.0.0
  github.com/faustbrian/go-service v1.0.0
)

replace github.com/faustbrian/go-router/v2 => $repository_root
EOF
cat >"$temporary/integration/integration_test.go" <<'EOF'
package integration_test

import (
  "context"
  "encoding/json"
  "net"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"

  middleware "github.com/faustbrian/go-http-middleware"
  jsonrpc "github.com/faustbrian/go-jsonrpc"
  router "github.com/faustbrian/go-router"
  routerv2 "github.com/faustbrian/go-router/v2"
  "github.com/faustbrian/go-service/serverhttp"
)

func TestOwnedHTTPBoundaries(t *testing.T) {
  registry := jsonrpc.NewRegistry()
  if err := registry.Register("ping", func(context.Context, json.RawMessage) (any, error) {
    return "pong", nil
  }); err != nil { t.Fatal(err) }
  rpc := jsonrpc.NewHTTPHandler(jsonrpc.NewDispatcher(registry))
  chain, err := middleware.New(func(next http.Handler) http.Handler { return next })
  if err != nil { t.Fatal(err) }
  wrappedRPC, err := chain.Handler(rpc)
  if err != nil { t.Fatal(err) }

  builder := router.New()
  if err := builder.Mount("/rpc", wrappedRPC, router.MountOptions{StripPrefix: true}); err != nil { t.Fatal(err) }
  if err := builder.Register(router.Route{
    Name: "track.webhook", Methods: []string{http.MethodPost}, Path: "/webhooks/track",
    Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }),
  }); err != nil { t.Fatal(err) }
  compiled, err := builder.Compile()
  if err != nil { t.Fatal(err) }

  var serverConstructor func(net.Listener, http.Handler, ...serverhttp.Option) (*serverhttp.Server, error) = serverhttp.New
  _ = serverConstructor
  request := httptest.NewRequest(http.MethodPost, "/rpc/", strings.NewReader(`{"jsonrpc":"2.0","method":"ping","id":1}`))
  request.Header.Set("Content-Type", "application/json")
  response := httptest.NewRecorder()
  compiled.ServeHTTP(response, request)
  if response.Code != http.StatusOK { t.Fatalf("RPC status: %d", response.Code) }
}

func TestCandidateV2HTTPBoundaries(t *testing.T) {
  registry := jsonrpc.NewRegistry()
  if err := registry.Register("ping", func(context.Context, json.RawMessage) (any, error) {
    return "pong", nil
  }); err != nil { t.Fatal(err) }
  rpc := jsonrpc.NewHTTPHandler(jsonrpc.NewDispatcher(registry))
  chain, err := middleware.New(func(next http.Handler) http.Handler { return next })
  if err != nil { t.Fatal(err) }
  wrappedRPC, err := chain.Handler(rpc)
  if err != nil { t.Fatal(err) }

  builder := routerv2.New()
  if err := builder.Mount("/rpc", wrappedRPC, routerv2.MountOptions{StripPrefix: true}); err != nil { t.Fatal(err) }
  compiled, err := builder.Compile()
  if err != nil { t.Fatal(err) }

  request := httptest.NewRequest(http.MethodPost, "/rpc/", strings.NewReader(`{"jsonrpc":"2.0","method":"ping","id":1}`))
  request.Header.Set("Content-Type", "application/json")
  response := httptest.NewRecorder()
  compiled.ServeHTTP(response, request)
  if response.Code != http.StatusOK { t.Fatalf("v2 RPC status: %d", response.Code) }

  trace := httptest.NewRecorder()
  compiled.ServeHTTP(trace, httptest.NewRequest(http.MethodTrace, "/rpc/", nil))
  if trace.Code != http.StatusMethodNotAllowed || trace.Header().Get("Allow") != "DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT" {
    t.Fatalf("v2 TRACE status=%d allow=%q", trace.Code, trace.Header().Get("Allow"))
  }
}
EOF
(cd "$temporary/integration" && go mod tidy && go test -race ./...)
