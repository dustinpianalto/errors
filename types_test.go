package errors_test

import (
	"testing"

	"github.com/dustinpianalto/errors"
)

func TestString(t *testing.T) {
	tt := []struct {
		Name string
		Kind errors.Kind
		Out  string
	}{
		{"other", errors.Other, "other error"},
		{"internal", errors.Internal, "internal error"},
		{"invalid", errors.Invalid, "invalid operation"},
		{"incorrect", errors.Incorrect, "incorrect configuration"},
		{"permission", errors.Permission, "permission denied"},
		{"io", errors.IO, "I/O error"},
		{"conflict", errors.Conflict, "item already exists"},
		{"not found", errors.NotFound, "item does not exist"},
		{"malformed", errors.Malformed, "malformed request"},
		{"unknown kind", errors.Kind(65535), "unknown type"},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			out := tc.Kind.String()
			if out != tc.Out {
				t.Fatalf("Expected: %#v Got: %#v", tc.Out, out)
			}
		})
	}
}
