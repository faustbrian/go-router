// Package routertest provides small consumer-facing helpers for compiled route
// tests without introducing a parallel runtime API. Start with MustCompile and
// Serve using a testing.TB-compatible value. Helpers invoke handlers synchronously
// in memory, start no server, add no deadline or panic recovery and do not own
// application resources; blocking and cancellation remain handler/caller policy.
package routertest

import (
	"context"
	"net/http"
	"net/http/httptest"

	router "github.com/faustbrian/go-router/v2"
)

// TestingT is the subset of testing.TB used by the helpers. It must be non-nil;
// real testing.TB values end the current test goroutine on Fatalf. Helpers return
// defensively after reporting failure even if a test double's Fatalf returns.
type TestingT interface {
	// Helper marks the caller as test support for failure locations.
	Helper()
	// Fatalf reports the formatted failure and normally ends the test goroutine.
	Fatalf(format string, arguments ...any)
}

// MustCompile compiles builder once or fails the test, returning nil if Fatalf
// returns. A nil builder is reported through its typed compile-state error;
// constructors' panics still propagate and successful compilation freezes it.
func MustCompile(testingT TestingT, builder *router.Builder) *router.Router {
	testingT.Helper()
	compiled, err := builder.Compile()
	if err != nil {
		testingT.Fatalf("compile router: %v", err)
		return nil
	}
	return compiled
}

// Serve synchronously sends a bodyless request with context.Background to
// handler and returns its recorder. Invalid method/target fails the test and
// returns nil if Fatalf returns. handler must be non-nil and owns blocking and
// panics; construct your own request when testing cancellation or deadlines.
func Serve(testingT TestingT, handler http.Handler, method, target string) *httptest.ResponseRecorder {
	testingT.Helper()
	request, err := http.NewRequestWithContext(context.Background(), method, target, nil)
	if err != nil {
		testingT.Fatalf("create request: %v", err)
		return nil
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

// AssertStatus fails the test for a nil response or unexpected recorded code;
// a nil response is reported as actual status zero.
func AssertStatus(testingT TestingT, response *httptest.ResponseRecorder, expected int) {
	testingT.Helper()
	if response == nil || response.Code != expected {
		actual := 0
		if response != nil {
			actual = response.Code
		}
		testingT.Fatalf("response status: got %d, want %d", actual, expected)
	}
}

// RouteTable returns a copied compiled route table or fails for a nil router.
func RouteTable(testingT TestingT, compiled *router.Router) []router.RouteInfo {
	testingT.Helper()
	if compiled == nil {
		testingT.Fatalf("route table: nil compiled router")
		return nil
	}
	return compiled.Routes()
}
