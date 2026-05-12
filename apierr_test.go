package apierr_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/renaldid/go-apierr"
)

func TestNew_NoAttrs(t *testing.T) {
	err := apierr.New(apierr.CodeNotFound, "user not found")

	if err.Code() != apierr.CodeNotFound {
		t.Fatalf("code: got %v, want %v", err.Code(), apierr.CodeNotFound)
	}
	if err.Message() != "user not found" {
		t.Fatalf("message: got %q, want %q", err.Message(), "user not found")
	}
	if len(err.Attrs()) != 0 {
		t.Fatalf("attrs: got %d, want 0", len(err.Attrs()))
	}
}

func TestNew_WithAttrs(t *testing.T) {
	err := apierr.New(apierr.CodeNotFound, "user not found", slog.String("id", "123"))

	if len(err.Attrs()) != 1 {
		t.Fatalf("attrs: got %d, want 1", len(err.Attrs()))
	}
	if err.Attrs()[0].Key != "id" {
		t.Fatalf("attr key: got %q, want %q", err.Attrs()[0].Key, "id")
	}
}

func TestNewf(t *testing.T) {
	err := apierr.Newf(apierr.CodeNotFound, "user %s not found", "abc")

	if err.Message() != "user abc not found" {
		t.Fatalf("message: got %q, want %q", err.Message(), "user abc not found")
	}
	if err.Code() != apierr.CodeNotFound {
		t.Fatalf("code: got %v, want %v", err.Code(), apierr.CodeNotFound)
	}
}

func TestWrap_NilCause(t *testing.T) {
	if got := apierr.Wrap(nil, apierr.CodeInternal, "msg"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestWrap_NonNilCause(t *testing.T) {
	cause := errors.New("db error")
	err := apierr.Wrap(cause, apierr.CodeInternal, "operation failed")

	if err.Code() != apierr.CodeInternal {
		t.Fatalf("code: got %v, want %v", err.Code(), apierr.CodeInternal)
	}
	if err.Message() != "operation failed" {
		t.Fatalf("message: got %q, want %q", err.Message(), "operation failed")
	}
	if !errors.Is(err, cause) {
		t.Fatal("Wrap must preserve cause in chain")
	}
}

func TestWrap_WithAttrs(t *testing.T) {
	cause := errors.New("cause")
	err := apierr.Wrap(cause, apierr.CodeInternal, "msg", slog.String("k", "v"))

	if len(err.Attrs()) != 1 {
		t.Fatalf("attrs: got %d, want 1", len(err.Attrs()))
	}
}

func TestWrapf_NilCause(t *testing.T) {
	if got := apierr.Wrapf(nil, apierr.CodeInternal, "msg %d", 1); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestWrapf_NonNilCause(t *testing.T) {
	cause := errors.New("cause")
	err := apierr.Wrapf(cause, apierr.CodeUnavailable, "retry after %d", 30)

	if err.Message() != "retry after 30" {
		t.Fatalf("message: got %q, want %q", err.Message(), "retry after 30")
	}
	if !errors.Is(err, cause) {
		t.Fatal("Wrapf must preserve cause in chain")
	}
}

func TestError_WithoutCause(t *testing.T) {
	err := apierr.New(apierr.CodeNotFound, "not found")

	if err.Error() != "not found" {
		t.Fatalf("got %q, want %q", err.Error(), "not found")
	}
}

func TestError_WithCause(t *testing.T) {
	cause := errors.New("connection refused")
	err := apierr.Wrap(cause, apierr.CodeUnavailable, "service down")

	want := "service down: connection refused"
	if err.Error() != want {
		t.Fatalf("got %q, want %q", err.Error(), want)
	}
}

func TestUnwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := apierr.Wrap(cause, apierr.CodeInternal, "msg")

	if errors.Unwrap(err) != cause {
		t.Fatal("Unwrap must return the original cause")
	}
}

func TestIs_SameCode(t *testing.T) {
	err := apierr.NotFound("resource missing")
	sentinel := apierr.New(apierr.CodeNotFound, "")

	if !errors.Is(err, sentinel) {
		t.Fatal("errors.Is must return true for matching Code")
	}
}

func TestIs_DifferentCode(t *testing.T) {
	err := apierr.NotFound("resource missing")
	sentinel := apierr.New(apierr.CodeInternal, "")

	if errors.Is(err, sentinel) {
		t.Fatal("errors.Is must return false for different Code")
	}
}

func TestIs_NonErrorTarget(t *testing.T) {
	err := apierr.NotFound("resource missing")

	if errors.Is(err, errors.New("plain")) {
		t.Fatal("errors.Is must return false for non-*Error target")
	}
}

func TestIs_WrappedChain(t *testing.T) {
	inner := apierr.NotFound("resource missing")
	outer := apierr.Wrap(inner, apierr.CodeInternal, "operation failed")
	sentinel := apierr.New(apierr.CodeNotFound, "")

	if !errors.Is(outer, sentinel) {
		t.Fatal("errors.Is must traverse the error chain")
	}
}

func TestCodeOf_DirectError(t *testing.T) {
	if got := apierr.CodeOf(apierr.NotFound("x")); got != apierr.CodeNotFound {
		t.Fatalf("got %v, want %v", got, apierr.CodeNotFound)
	}
}

func TestCodeOf_WrappedError(t *testing.T) {
	inner := apierr.NotFound("x")
	outer := apierr.Wrap(inner, apierr.CodeInternal, "wrapped")

	// errors.As finds the outermost *Error first
	if got := apierr.CodeOf(outer); got != apierr.CodeInternal {
		t.Fatalf("got %v, want %v", got, apierr.CodeInternal)
	}
}

func TestCodeOf_PlainError(t *testing.T) {
	if got := apierr.CodeOf(errors.New("plain")); got != apierr.CodeUnknown {
		t.Fatalf("got %v, want %v", got, apierr.CodeUnknown)
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *apierr.Error
		want apierr.Code
	}{
		{"Canceled", apierr.Canceled("m"), apierr.CodeCanceled},
		{"InvalidArgument", apierr.InvalidArgument("m"), apierr.CodeInvalidArgument},
		{"DeadlineExceeded", apierr.DeadlineExceeded("m"), apierr.CodeDeadlineExceeded},
		{"NotFound", apierr.NotFound("m"), apierr.CodeNotFound},
		{"AlreadyExists", apierr.AlreadyExists("m"), apierr.CodeAlreadyExists},
		{"PermissionDenied", apierr.PermissionDenied("m"), apierr.CodePermissionDenied},
		{"ResourceExhausted", apierr.ResourceExhausted("m"), apierr.CodeResourceExhausted},
		{"FailedPrecondition", apierr.FailedPrecondition("m"), apierr.CodeFailedPrecondition},
		{"Aborted", apierr.Aborted("m"), apierr.CodeAborted},
		{"Unimplemented", apierr.Unimplemented("m"), apierr.CodeUnimplemented},
		{"Internal", apierr.Internal("m"), apierr.CodeInternal},
		{"Unavailable", apierr.Unavailable("m"), apierr.CodeUnavailable},
		{"Unauthenticated", apierr.Unauthenticated("m"), apierr.CodeUnauthenticated},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Code() != tc.want {
				t.Fatalf("code: got %v, want %v", tc.err.Code(), tc.want)
			}
			if tc.err.Message() != "m" {
				t.Fatalf("message: got %q, want %q", tc.err.Message(), "m")
			}
		})
	}
}
