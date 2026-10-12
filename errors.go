// Package router compiles explicit route descriptors into immutable net/http
// handlers. Start with New, register routes or groups, then Compile; the builder
// is single-owner, while the compiled indexes support concurrent dispatch,
// introspection and named URL generation. Handler and middleware synchronization,
// cancellation, panics and lifecycle remain the application's responsibility.
//
// Matching uses ServeMux's literal and wildcard grammar. Redirect serialization
// deliberately retains Go 1.26.6 behavior; HTTP method and URI-generation choices
// are recorded in docs/specification-decisions.md. This standard-library-only
// package starts no servers or goroutines and supplies no authentication, proxy
// trust, service discovery or cgo dependency.
package router

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	// ErrInvalidRoute identifies a malformed route descriptor.
	ErrInvalidRoute = errors.New("invalid route")
	// ErrConflict identifies ambiguous or duplicate semantic routes.
	ErrConflict = errors.New("route conflict")
	// ErrDuplicateName identifies a repeated stable route name.
	ErrDuplicateName = errors.New("duplicate route name")
	// ErrInvalidParameter identifies malformed URL-generation parameters.
	ErrInvalidParameter = errors.New("invalid route parameter")
	// ErrGeneration identifies a named-route URL generation failure.
	ErrGeneration = errors.New("route generation failed")
	// ErrUnsupported identifies deliberately unsupported routing, including CONNECT.
	ErrUnsupported = errors.New("unsupported routing behavior")
	// ErrCompileState identifies nil builders/routers or a successfully frozen builder.
	ErrCompileState = errors.New("invalid router compile state")
	// ErrLimitExceeded identifies a configured resource-budget violation.
	ErrLimitExceeded = errors.New("router limit exceeded")
)

// Error is a startup or generation diagnostic. Kind supports errors.Is; callers
// may use errors.As to inspect its unsanitized fields. Error sanitizes controls
// and bounds rendered text, but does not redact caller-provided secrets. Treat
// source labels and custom Kind values as disclosure-bearing inputs. A zero
// Error and a nil *Error both render as "router error".
type Error struct {
	// Kind is the underlying category returned by Unwrap; nil has no category.
	Kind error
	// Field identifies the rejected input dimension; empty omits it from text.
	Field string
	// Source is an optional caller-owned diagnostic label, not a filesystem read.
	// Keep it non-sensitive; rendered text truncates it independently.
	Source string
	// Detail describes the refusal. Error bounds and sanitizes its rendered form.
	Detail string
}

// Error returns single-line valid UTF-8 bounded to 160 bytes, including any
// caller-supplied source label. It replaces controls and truncates, not redacts.
func (e *Error) Error() string {
	if e == nil {
		return "router error"
	}
	kind := "router error"
	if e.Kind != nil {
		kind = sanitizeDiagnostic(e.Kind.Error(), 80)
	}
	message := kind
	if e.Field != "" {
		message += ": " + sanitizeDiagnostic(e.Field, 32)
	}
	if e.Detail != "" {
		message += ": " + sanitizeDiagnostic(e.Detail, 80)
	}
	if e.Source != "" {
		message += fmt.Sprintf(" (source %s)", sanitizeDiagnostic(e.Source, 32))
	}
	return bounded(message, 160)
}

// Unwrap exposes the stable error category.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Kind
}

// Bound input before UTF-8 repair so hostile diagnostics cannot cause work
// proportional to their unbounded tail. Repair can expand bytes, requiring a
// second boundary check; mutation_boundaries_internal_test.go protects this.
func bounded(value string, limit int) string {
	if limit < 1 {
		return ""
	}
	inputTruncated := len(value) > limit
	if inputTruncated {
		value = value[:limit]
	}
	value = strings.ToValidUTF8(value, "�")
	if !inputTruncated && len(value) <= limit {
		return value
	}
	if limit <= 3 {
		return value[:validUTF8Boundary(value, limit)]
	}
	return value[:validUTF8Boundary(value, limit-3)] + "..."
}

func validUTF8Boundary(value string, boundary int) int {
	if boundary >= len(value) {
		return len(value)
	}
	for boundary != 0 && !utf8.RuneStart(value[boundary]) {
		boundary--
	}
	return boundary
}

func sanitizeDiagnostic(value string, limit int) string {
	value = bounded(value, limit)
	return strings.Map(func(character rune) rune {
		if unicode.IsControl(character) {
			return '?'
		}
		return character
	}, value)
}
