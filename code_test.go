package apierr_test

import (
	"net/http"
	"testing"

	"github.com/renaldid/go-apierr"
)

func TestCode_HTTPStatus(t *testing.T) {
	tests := []struct {
		code apierr.Code
		want int
	}{
		{apierr.CodeUnknown, http.StatusInternalServerError},
		{apierr.CodeCanceled, http.StatusRequestTimeout},
		{apierr.CodeInvalidArgument, http.StatusBadRequest},
		{apierr.CodeDeadlineExceeded, http.StatusGatewayTimeout},
		{apierr.CodeNotFound, http.StatusNotFound},
		{apierr.CodeAlreadyExists, http.StatusConflict},
		{apierr.CodePermissionDenied, http.StatusForbidden},
		{apierr.CodeResourceExhausted, http.StatusTooManyRequests},
		{apierr.CodeFailedPrecondition, http.StatusBadRequest},
		{apierr.CodeAborted, http.StatusConflict},
		{apierr.CodeUnimplemented, http.StatusNotImplemented},
		{apierr.CodeInternal, http.StatusInternalServerError},
		{apierr.CodeUnavailable, http.StatusServiceUnavailable},
		{apierr.CodeUnauthenticated, http.StatusUnauthorized},
		{apierr.Code(99), http.StatusInternalServerError}, // unrecognized → 500
	}
	for _, tc := range tests {
		t.Run(tc.code.String(), func(t *testing.T) {
			if got := tc.code.HTTPStatus(); got != tc.want {
				t.Fatalf("HTTPStatus: got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestCode_String(t *testing.T) {
	tests := []struct {
		code apierr.Code
		want string
	}{
		{apierr.CodeUnknown, "unknown"},
		{apierr.CodeCanceled, "canceled"},
		{apierr.CodeInvalidArgument, "invalid_argument"},
		{apierr.CodeDeadlineExceeded, "deadline_exceeded"},
		{apierr.CodeNotFound, "not_found"},
		{apierr.CodeAlreadyExists, "already_exists"},
		{apierr.CodePermissionDenied, "permission_denied"},
		{apierr.CodeResourceExhausted, "resource_exhausted"},
		{apierr.CodeFailedPrecondition, "failed_precondition"},
		{apierr.CodeAborted, "aborted"},
		{apierr.CodeUnimplemented, "unimplemented"},
		{apierr.CodeInternal, "internal"},
		{apierr.CodeUnavailable, "unavailable"},
		{apierr.CodeUnauthenticated, "unauthenticated"},
		{apierr.Code(99), "unknown"}, // unrecognized → "unknown"
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.code.String(); got != tc.want {
				t.Fatalf("String: got %q, want %q", got, tc.want)
			}
		})
	}
}
