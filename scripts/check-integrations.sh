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

export GOCACHE="$temporary/gocache"
export GOMODCACHE="$temporary/gomodcache"
export GOPATH="$temporary/gopath"
export GOWORK=off

mkdir "$temporary/integration"
cat >"$temporary/integration/go.mod" <<EOF
module routerintegration

go 1.26.6

require (
  github.com/faustbrian/go-http-middleware v1.0.0
  github.com/faustbrian/go-jsonrpc v1.0.0
  github.com/faustbrian/go-router v1.0.0
  github.com/faustbrian/go-service v1.0.0
)
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
EOF
(cd "$temporary/integration" && go mod tidy && go test -race ./...)
