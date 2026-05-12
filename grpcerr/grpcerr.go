// Package grpcerr bridges [apierr.Error] with gRPC status codes.
//
// This package is a separate Go module so that the core apierr package remains
// zero-dependency. Import it only when your service also imports google.golang.org/grpc.
//
//	grpcErr := grpcerr.ToGRPC(err)   // *Error → gRPC status error
//	apiErr  := grpcerr.FromGRPC(err) // gRPC status error → *Error
package grpcerr

import (
	"errors"

	"github.com/renaldid/go-apierr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var toGRPCCode = map[apierr.Code]codes.Code{
	apierr.CodeUnknown:            codes.Unknown,
	apierr.CodeCanceled:           codes.Canceled,
	apierr.CodeInvalidArgument:    codes.InvalidArgument,
	apierr.CodeDeadlineExceeded:   codes.DeadlineExceeded,
	apierr.CodeNotFound:           codes.NotFound,
	apierr.CodeAlreadyExists:      codes.AlreadyExists,
	apierr.CodePermissionDenied:   codes.PermissionDenied,
	apierr.CodeResourceExhausted:  codes.ResourceExhausted,
	apierr.CodeFailedPrecondition: codes.FailedPrecondition,
	apierr.CodeAborted:            codes.Aborted,
	apierr.CodeUnimplemented:      codes.Unimplemented,
	apierr.CodeInternal:           codes.Internal,
	apierr.CodeUnavailable:        codes.Unavailable,
	apierr.CodeUnauthenticated:    codes.Unauthenticated,
}

// fromGRPCCode is the reverse of toGRPCCode, built once at init time.
var fromGRPCCode = func() map[codes.Code]apierr.Code {
	m := make(map[codes.Code]apierr.Code, len(toGRPCCode))
	for k, v := range toGRPCCode {
		m[v] = k
	}
	return m
}()

// ToGRPC converts err to a gRPC status error.
// If err is (or wraps) an *apierr.Error, the gRPC code is derived from its Code.
// Plain errors and unrecognized codes map to codes.Unknown.
// A nil input returns nil.
func ToGRPC(err error) error {
	if err == nil {
		return nil
	}
	var e *apierr.Error
	if errors.As(err, &e) {
		c, ok := toGRPCCode[e.Code()]
		if !ok {
			c = codes.Unknown
		}
		return status.Error(c, e.Error())
	}
	return status.Error(codes.Unknown, err.Error())
}

// Code returns the gRPC codes.Code corresponding to the apierr.Code found in err's chain.
// Returns codes.OK for nil and codes.Unknown for plain errors or unrecognized codes.
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if c, ok := toGRPCCode[apierr.CodeOf(err)]; ok {
		return c
	}
	return codes.Unknown
}

// FromGRPC converts a gRPC status error to an *apierr.Error.
// Non-gRPC errors map to CodeInternal. Unmapped gRPC codes map to CodeUnknown.
// A nil input returns nil.
func FromGRPC(err error) error {
	if err == nil {
		return nil
	}
	s, ok := status.FromError(err)
	if !ok {
		return apierr.New(apierr.CodeInternal, err.Error())
	}
	if apiCode, ok := fromGRPCCode[s.Code()]; ok {
		return apierr.New(apiCode, s.Message())
	}
	return apierr.New(apierr.CodeUnknown, s.Message())
}
