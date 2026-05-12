package grpcerr_test

import (
	"errors"
	"testing"

	"github.com/renaldid/go-apierr"
	"github.com/renaldid/go-apierr/grpcerr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPC_Nil(t *testing.T) {
	if grpcerr.ToGRPC(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestToGRPC_ApiError(t *testing.T) {
	err := apierr.NotFound("user not found")
	s, ok := status.FromError(grpcerr.ToGRPC(err))
	if !ok {
		t.Fatal("expected a gRPC status error")
	}
	if s.Code() != codes.NotFound {
		t.Fatalf("Code: got %v, want %v", s.Code(), codes.NotFound)
	}
	if s.Message() != "user not found" {
		t.Fatalf("Message: got %q, want %q", s.Message(), "user not found")
	}
}

func TestToGRPC_WrappedApiError(t *testing.T) {
	inner := apierr.NotFound("not found")
	outer := apierr.Wrap(inner, apierr.CodeUnavailable, "service down")
	s, _ := status.FromError(grpcerr.ToGRPC(outer))
	if s.Code() != codes.Unavailable {
		t.Fatalf("Code: got %v, want %v", s.Code(), codes.Unavailable)
	}
}

func TestToGRPC_PlainError(t *testing.T) {
	s, _ := status.FromError(grpcerr.ToGRPC(errors.New("plain")))
	if s.Code() != codes.Unknown {
		t.Fatalf("Code: got %v, want %v", s.Code(), codes.Unknown)
	}
}

func TestToGRPC_UnrecognizedCode(t *testing.T) {
	err := apierr.New(apierr.Code(99), "unrecognized")
	s, _ := status.FromError(grpcerr.ToGRPC(err))
	if s.Code() != codes.Unknown {
		t.Fatalf("Code: got %v, want %v", s.Code(), codes.Unknown)
	}
}

func TestCode_Nil(t *testing.T) {
	if grpcerr.Code(nil) != codes.OK {
		t.Fatalf("got %v, want OK", grpcerr.Code(nil))
	}
}

func TestCode_ApiError(t *testing.T) {
	if got := grpcerr.Code(apierr.InvalidArgument("bad")); got != codes.InvalidArgument {
		t.Fatalf("got %v, want %v", got, codes.InvalidArgument)
	}
}

func TestCode_PlainError(t *testing.T) {
	if got := grpcerr.Code(errors.New("plain")); got != codes.Unknown {
		t.Fatalf("got %v, want %v", got, codes.Unknown)
	}
}

func TestCode_UnrecognizedCode(t *testing.T) {
	if got := grpcerr.Code(apierr.New(apierr.Code(99), "x")); got != codes.Unknown {
		t.Fatalf("got %v, want %v", got, codes.Unknown)
	}
}

func TestFromGRPC_Nil(t *testing.T) {
	if grpcerr.FromGRPC(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestFromGRPC_StatusError(t *testing.T) {
	grpcErr := status.Error(codes.NotFound, "resource missing")
	err := grpcerr.FromGRPC(grpcErr)
	if apierr.CodeOf(err) != apierr.CodeNotFound {
		t.Fatalf("Code: got %v, want %v", apierr.CodeOf(err), apierr.CodeNotFound)
	}
}

func TestFromGRPC_UnmappedCode(t *testing.T) {
	grpcErr := status.Error(codes.DataLoss, "data lost")
	err := grpcerr.FromGRPC(grpcErr)
	if apierr.CodeOf(err) != apierr.CodeUnknown {
		t.Fatalf("Code: got %v, want %v", apierr.CodeOf(err), apierr.CodeUnknown)
	}
}

func TestFromGRPC_NonStatusError(t *testing.T) {
	err := grpcerr.FromGRPC(errors.New("not a status error"))
	if apierr.CodeOf(err) != apierr.CodeInternal {
		t.Fatalf("Code: got %v, want %v", apierr.CodeOf(err), apierr.CodeInternal)
	}
}

func TestRoundTrip_AllMappedCodes(t *testing.T) {
	pairs := []struct {
		api  apierr.Code
		grpc codes.Code
	}{
		{apierr.CodeUnknown, codes.Unknown},
		{apierr.CodeCanceled, codes.Canceled},
		{apierr.CodeInvalidArgument, codes.InvalidArgument},
		{apierr.CodeDeadlineExceeded, codes.DeadlineExceeded},
		{apierr.CodeNotFound, codes.NotFound},
		{apierr.CodeAlreadyExists, codes.AlreadyExists},
		{apierr.CodePermissionDenied, codes.PermissionDenied},
		{apierr.CodeResourceExhausted, codes.ResourceExhausted},
		{apierr.CodeFailedPrecondition, codes.FailedPrecondition},
		{apierr.CodeAborted, codes.Aborted},
		{apierr.CodeUnimplemented, codes.Unimplemented},
		{apierr.CodeInternal, codes.Internal},
		{apierr.CodeUnavailable, codes.Unavailable},
		{apierr.CodeUnauthenticated, codes.Unauthenticated},
	}
	for _, tc := range pairs {
		t.Run(tc.api.String(), func(t *testing.T) {
			s, _ := status.FromError(grpcerr.ToGRPC(apierr.New(tc.api, "msg")))
			if s.Code() != tc.grpc {
				t.Fatalf("ToGRPC: got %v, want %v", s.Code(), tc.grpc)
			}
			roundtripped := grpcerr.FromGRPC(status.Error(tc.grpc, "msg"))
			if apierr.CodeOf(roundtripped) != tc.api {
				t.Fatalf("FromGRPC: got %v, want %v", apierr.CodeOf(roundtripped), tc.api)
			}
		})
	}
}
