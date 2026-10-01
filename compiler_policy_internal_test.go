package router

import (
	"errors"
	"net/http"
	"testing"
)

func TestCompileRejectsInvalidGlobalMiddlewarePolicyBeforeConstruction(t *testing.T) {
	t.Parallel()

	constructions := 0
	builder := New(WithMiddleware(NamedMiddleware{
		Middleware: func(next http.Handler) http.Handler {
			constructions++
			return next
		},
	}))
	if err := builder.Register(Route{Methods: []string{"GET"}, Path: "/", Handler: valueHandler{}}); err != nil {
		t.Fatal(err)
	}
	// Exercise the compile-time defense independently of New's option checks.
	builder.globalMiddleware[0].ExclusionPolicy = MiddlewareExclusionPolicy(255)

	compiled, err := builder.Compile()
	var detail *Error
	if !errors.Is(err, ErrInvalidRoute) || !errors.As(err, &detail) || detail.Field != "middleware" {
		t.Fatalf("invalid policy compile = (%v, %v); want middleware ErrInvalidRoute", compiled, err)
	}
	if compiled != nil || constructions != 0 || builder.compiled {
		t.Fatalf("rejected compile published a router or constructed middleware: router=%v calls=%d compiled=%v", compiled, constructions, builder.compiled)
	}

	builder.globalMiddleware[0].ExclusionPolicy = MiddlewareExclusionDenied
	compiled, err = builder.Compile()
	if err != nil || compiled == nil || constructions != 1 || !builder.compiled {
		t.Fatalf("corrected policy compile = (%v, %v), calls=%d compiled=%v", compiled, err, constructions, builder.compiled)
	}
}
