// Package apierr provides a structured API error type that carries a semantic Code,
// a human-readable message, slog attributes, and an optional wrapped cause.
//
// A single *Error value satisfies the error interface, implements slog.LogValuer,
// and renders as RFC 9457 Problem Details JSON via [ToProblem] and [WriteHTTP].
// The Code aligns semantically with gRPC status codes; the companion grpcerr module
// provides direct gRPC status conversion.
package apierr

import (
	"errors"
	"fmt"
	"log/slog"
)

// Error is a structured API error.
//
// Two *Error values are considered equal by [errors.Is] when their Codes match,
// enabling sentinel-style checking:
//
//	var ErrNotFound = apierr.New(CodeNotFound, "")
//	errors.Is(err, ErrNotFound) // true for any CodeNotFound error
type Error struct {
	code    Code
	message string
	attrs   []slog.Attr
	cause   error
}

// New returns a new *Error with the given code, message, and optional slog attributes.
func New(code Code, message string, attrs ...slog.Attr) *Error {
	e := &Error{code: code, message: message}
	if len(attrs) > 0 {
		e.attrs = attrs
	}
	return e
}

// Newf returns a new *Error with a [fmt.Sprintf]-formatted message.
func Newf(code Code, format string, args ...any) *Error {
	return &Error{code: code, message: fmt.Sprintf(format, args...)}
}

// Wrap returns a new *Error that wraps cause with the given code, message, and optional
// slog attributes. If cause is nil, Wrap returns nil.
func Wrap(cause error, code Code, message string, attrs ...slog.Attr) *Error {
	if cause == nil {
		return nil
	}
	e := &Error{code: code, message: message, cause: cause}
	if len(attrs) > 0 {
		e.attrs = attrs
	}
	return e
}

// Wrapf returns a new *Error wrapping cause with a [fmt.Sprintf]-formatted message.
// If cause is nil, Wrapf returns nil.
func Wrapf(cause error, code Code, format string, args ...any) *Error {
	if cause == nil {
		return nil
	}
	return &Error{code: code, message: fmt.Sprintf(format, args...), cause: cause}
}

// Error implements the error interface. If a cause is present, it is appended
// after the message, separated by ": ".
func (e *Error) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

// Unwrap returns the wrapped cause, enabling [errors.Is] and [errors.As] chain traversal.
func (e *Error) Unwrap() error { return e.cause }

// Is reports whether e matches target.
// Two *Error values match when their Codes are equal, regardless of message or attrs.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && e.code == t.code
}

// Code returns the semantic code of e.
func (e *Error) Code() Code { return e.code }

// Message returns e's message without the cause chain.
func (e *Error) Message() string { return e.message }

// Attrs returns the slog attributes attached directly to e, not the cause chain.
func (e *Error) Attrs() []slog.Attr { return e.attrs }

// CodeOf returns the Code of the first *Error found in err's chain.
// Returns CodeUnknown if err does not contain an *Error.
func CodeOf(err error) Code {
	var e *Error
	if errors.As(err, &e) {
		return e.code
	}
	return CodeUnknown
}

// Canceled returns a new *Error with CodeCanceled.
func Canceled(message string, attrs ...slog.Attr) *Error {
	return New(CodeCanceled, message, attrs...)
}

// InvalidArgument returns a new *Error with CodeInvalidArgument.
func InvalidArgument(message string, attrs ...slog.Attr) *Error {
	return New(CodeInvalidArgument, message, attrs...)
}

// DeadlineExceeded returns a new *Error with CodeDeadlineExceeded.
func DeadlineExceeded(message string, attrs ...slog.Attr) *Error {
	return New(CodeDeadlineExceeded, message, attrs...)
}

// NotFound returns a new *Error with CodeNotFound.
func NotFound(message string, attrs ...slog.Attr) *Error {
	return New(CodeNotFound, message, attrs...)
}

// AlreadyExists returns a new *Error with CodeAlreadyExists.
func AlreadyExists(message string, attrs ...slog.Attr) *Error {
	return New(CodeAlreadyExists, message, attrs...)
}

// PermissionDenied returns a new *Error with CodePermissionDenied.
func PermissionDenied(message string, attrs ...slog.Attr) *Error {
	return New(CodePermissionDenied, message, attrs...)
}

// ResourceExhausted returns a new *Error with CodeResourceExhausted.
func ResourceExhausted(message string, attrs ...slog.Attr) *Error {
	return New(CodeResourceExhausted, message, attrs...)
}

// FailedPrecondition returns a new *Error with CodeFailedPrecondition.
func FailedPrecondition(message string, attrs ...slog.Attr) *Error {
	return New(CodeFailedPrecondition, message, attrs...)
}

// Aborted returns a new *Error with CodeAborted.
func Aborted(message string, attrs ...slog.Attr) *Error {
	return New(CodeAborted, message, attrs...)
}

// Unimplemented returns a new *Error with CodeUnimplemented.
func Unimplemented(message string, attrs ...slog.Attr) *Error {
	return New(CodeUnimplemented, message, attrs...)
}

// Internal returns a new *Error with CodeInternal.
func Internal(message string, attrs ...slog.Attr) *Error {
	return New(CodeInternal, message, attrs...)
}

// Unavailable returns a new *Error with CodeUnavailable.
func Unavailable(message string, attrs ...slog.Attr) *Error {
	return New(CodeUnavailable, message, attrs...)
}

// Unauthenticated returns a new *Error with CodeUnauthenticated.
func Unauthenticated(message string, attrs ...slog.Attr) *Error {
	return New(CodeUnauthenticated, message, attrs...)
}
