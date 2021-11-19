package errors_test

import (
	stdioerrors "errors"
	"testing"

	"github.com/dustinpianalto/errors"
)

func TestE(t *testing.T) {
	var emptyError = &errors.Error{Message: errors.Message("foo")}
	var stdioError = stdioerrors.New("foo")
	tt := []struct {
		Name     string
		Kind     errors.Kind
		Method   errors.Method
		Username errors.Username
		Message  errors.Message
		String   string
		Err      error
		Out      *errors.Error
	}{
		{Name: "just kind", Kind: errors.Permission, Out: &errors.Error{Kind: errors.Permission}},
		{Name: "just method", Method: errors.Method("foo"), Out: &errors.Error{Kind: errors.Other, Method: errors.Method("foo")}},
		{Name: "just username", Username: errors.Username("foo"), Out: &errors.Error{Kind: errors.Other, Username: errors.Username("foo")}},
		{Name: "just message", Message: errors.Message("foo"), Out: &errors.Error{Kind: errors.Other, Message: errors.Message("foo")}},
		{Name: "just string", String: "foo", Out: &errors.Error{Kind: errors.Other, Message: errors.Message("foo")}},
		{Name: "just stdio error", Err: stdioError, Out: &errors.Error{Kind: errors.Other, Err: stdioError}},
		{Name: "just error", Err: emptyError, Out: &errors.Error{Kind: errors.Other, Err: emptyError}},
		{
			Name:     "fully populated stdio error",
			Kind:     errors.Permission,
			Method:   errors.Method("foo"),
			Username: errors.Username("foo"),
			Message:  errors.Message("foo"),
			Err:      stdioError,
			Out: &errors.Error{
				Kind:     errors.Permission,
				Method:   errors.Method("foo"),
				Username: errors.Username("foo"),
				Message:  errors.Message("foo"),
				Err:      stdioError,
			}},
		{
			Name:     "fully populated stdio error",
			Kind:     errors.Permission,
			Method:   errors.Method("foo"),
			Username: errors.Username("foo"),
			Message:  errors.Message("foo"),
			Err:      emptyError,
			Out: &errors.Error{
				Kind:     errors.Permission,
				Method:   errors.Method("foo"),
				Username: errors.Username("foo"),
				Message:  errors.Message("foo"),
				Err:      emptyError,
			}},
		{
			Name: "nil error",
			Err:  nil,
			Out:  &errors.Error{Message: errors.Message("unknown type <nil> with value <nil> in error call")},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var e error
			if tc.Err != nil || tc.Name == "nil error" {
				e = errors.E(tc.Kind, tc.Method, tc.Username, tc.Message, tc.String, tc.Err)
			} else {
				e = errors.E(tc.Kind, tc.Method, tc.Username, tc.Message, tc.String)
			}
			if !errors.Match(tc.Out, e) {
				t.Fatalf("Expected: %#v Got: %#v", tc.Out, e)
			}
		})
		t.Run("empty args", func(t *testing.T) {
			e := errors.E()
			if e != nil {
				t.Fatalf("Expected: %#v Got: %#v", nil, e)
			}
		})
	}
}

func TestError(t *testing.T) {
	var emptyError = &errors.Error{Message: errors.Message("foo")}
	var stdioError = stdioerrors.New("foo")
	tt := []struct {
		Name string
		Err  error
		Out  string
	}{
		{Name: "just kind", Err: &errors.Error{Kind: errors.Permission}, Out: "permission denied"},
		{Name: "just method", Err: &errors.Error{Method: errors.Method("foo")}, Out: "foo"},
		{Name: "just username", Err: &errors.Error{Username: errors.Username("foo")}, Out: "foo"},
		{Name: "just message", Err: &errors.Error{Message: errors.Message("foo")}, Out: "foo"},
		{Name: "just stdio error", Err: &errors.Error{Err: stdioError}, Out: "foo"},
		{Name: "just error", Err: &errors.Error{Kind: errors.Other, Err: emptyError}, Out: "foo"},
		{
			Name: "fully populated stdio error",
			Err: &errors.Error{
				Kind:     errors.Permission,
				Method:   errors.Method("method"),
				Username: errors.Username("username"),
				Message:  errors.Message("message"),
				Err:      stdioError,
			},
			Out: "method: username: permission denied: message: foo",
		},
		{
			Name: "fully populated stdio error",
			Err: &errors.Error{
				Kind:     errors.Permission,
				Method:   errors.Method("method"),
				Username: errors.Username("username"),
				Message:  errors.Message("message"),
				Err:      emptyError,
			},
			Out: "method: username: permission denied: message:\n\tfoo",
		},
		{
			Name: "empty error",
			Err:  &errors.Error{},
			Out:  "no error message",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			out := tc.Err.Error()

			if out != tc.Out {
				t.Fatalf("Expected: %#v Got: %#v", tc.Out, out)
			}
		})
	}
}

func TestIs(t *testing.T) {
	tt := []struct {
		Name string
		Err  error
		Kind errors.Kind
		Out  bool
	}{
		{"is permission", &errors.Error{Kind: errors.Permission}, errors.Permission, true},
		{"is not permission", &errors.Error{}, errors.Permission, false},
		{"nested is permission", &errors.Error{Err: &errors.Error{Kind: errors.Permission}}, errors.Permission, true},
		{"is not errors.Error", stdioerrors.New("foo"), errors.Permission, false},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if out := errors.Is(tc.Kind, tc.Err); out != tc.Out {
				t.Fatalf("Expected: %#v Got: %#v", tc.Out, out)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	tt := []struct {
		Name string
		Err1 error
		Err2 error
		Out  bool
	}{
		{
			Name: "err1 not Error",
			Err1: stdioerrors.New("foo"),
			Err2: nil,
			Out:  false,
		},
		{
			Name: "err2 not Error",
			Err1: &errors.Error{},
			Err2: stdioerrors.New("foo"),
			Out:  false,
		},
		{
			Name: "kind not same",
			Err1: &errors.Error{Kind: errors.Permission},
			Err2: &errors.Error{Kind: errors.Conflict},
			Out:  false,
		},
		{
			Name: "method not same",
			Err1: &errors.Error{Method: errors.Method("foo")},
			Err2: &errors.Error{Method: errors.Method("bar")},
			Out:  false,
		},
		{
			Name: "username not same",
			Err1: &errors.Error{Username: errors.Username("foo")},
			Err2: &errors.Error{Username: errors.Username("bar")},
			Out:  false,
		},
		{
			Name: "message not same",
			Err1: &errors.Error{Message: errors.Message("foo")},
			Err2: &errors.Error{Message: errors.Message("bar")},
			Out:  false,
		},
		{
			Name: "err2 is nil",
			Err1: &errors.Error{Err: stdioerrors.New("foo")},
			Err2: &errors.Error{},
			Out:  false,
		},
		{
			Name: "nested errors not same",
			Err1: &errors.Error{Err: &errors.Error{Kind: errors.Permission}},
			Err2: &errors.Error{Err: &errors.Error{Kind: errors.Conflict}},
			Out:  false,
		},
		{
			Name: "nested errors are same",
			Err1: &errors.Error{Err: &errors.Error{}},
			Err2: &errors.Error{Err: &errors.Error{}},
			Out:  true,
		},
		{
			Name: "all same",
			Err1: &errors.Error{
				Kind:     errors.Permission,
				Method:   errors.Method("method"),
				Username: errors.Username("username"),
				Message:  errors.Message("message"),
				Err:      nil,
			},
			Err2: &errors.Error{
				Kind:     errors.Permission,
				Method:   errors.Method("method"),
				Username: errors.Username("username"),
				Message:  errors.Message("message"),
				Err:      nil,
			},
			Out: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			out := errors.Match(tc.Err1, tc.Err2)
			if out != tc.Out {
				t.Fatalf("Expected: %#v Got: %#v", tc.Out, out)
			}
		})
	}
}
