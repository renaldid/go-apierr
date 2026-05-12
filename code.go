package apierr

import "net/http"

// Code identifies the class of an API error.
// Codes are semantically aligned with gRPC status codes to simplify cross-protocol mapping.
type Code uint32

const (
	CodeUnknown            Code = 0
	CodeCanceled           Code = 1
	CodeInvalidArgument    Code = 2
	CodeDeadlineExceeded   Code = 3
	CodeNotFound           Code = 4
	CodeAlreadyExists      Code = 5
	CodePermissionDenied   Code = 6
	CodeResourceExhausted  Code = 7
	CodeFailedPrecondition Code = 8
	CodeAborted            Code = 9
	CodeUnimplemented      Code = 10
	CodeInternal           Code = 11
	CodeUnavailable        Code = 12
	CodeUnauthenticated    Code = 13
)

var codeStrings = map[Code]string{
	CodeUnknown:            "unknown",
	CodeCanceled:           "canceled",
	CodeInvalidArgument:    "invalid_argument",
	CodeDeadlineExceeded:   "deadline_exceeded",
	CodeNotFound:           "not_found",
	CodeAlreadyExists:      "already_exists",
	CodePermissionDenied:   "permission_denied",
	CodeResourceExhausted:  "resource_exhausted",
	CodeFailedPrecondition: "failed_precondition",
	CodeAborted:            "aborted",
	CodeUnimplemented:      "unimplemented",
	CodeInternal:           "internal",
	CodeUnavailable:        "unavailable",
	CodeUnauthenticated:    "unauthenticated",
}

var httpStatuses = map[Code]int{
	CodeUnknown:            http.StatusInternalServerError,
	CodeCanceled:           http.StatusRequestTimeout,
	CodeInvalidArgument:    http.StatusBadRequest,
	CodeDeadlineExceeded:   http.StatusGatewayTimeout,
	CodeNotFound:           http.StatusNotFound,
	CodeAlreadyExists:      http.StatusConflict,
	CodePermissionDenied:   http.StatusForbidden,
	CodeResourceExhausted:  http.StatusTooManyRequests,
	CodeFailedPrecondition: http.StatusBadRequest,
	CodeAborted:            http.StatusConflict,
	CodeUnimplemented:      http.StatusNotImplemented,
	CodeInternal:           http.StatusInternalServerError,
	CodeUnavailable:        http.StatusServiceUnavailable,
	CodeUnauthenticated:    http.StatusUnauthorized,
}

// HTTPStatus returns the HTTP status code corresponding to c.
// Unrecognized codes map to 500 Internal Server Error.
func (c Code) HTTPStatus() int {
	if s, ok := httpStatuses[c]; ok {
		return s
	}
	return http.StatusInternalServerError
}

// String returns c's snake_case name (e.g. "not_found").
// Unrecognized codes return "unknown".
func (c Code) String() string {
	if s, ok := codeStrings[c]; ok {
		return s
	}
	return "unknown"
}
