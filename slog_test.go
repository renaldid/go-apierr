package apierr_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/renaldid/go-apierr"
)

func TestLogValue_WithAttrs(t *testing.T) {
	err := apierr.NotFound("user not found", slog.String("id", "123"))
	lv := err.LogValue()

	if lv.Kind() != slog.KindGroup {
		t.Fatalf("Kind: got %v, want KindGroup", lv.Kind())
	}
	attrs := lv.Group()
	if len(attrs) != 3 { // code, message, id
		t.Fatalf("group len: got %d, want 3", len(attrs))
	}
	if attrs[0].Key != "code" {
		t.Fatalf("attrs[0].Key: got %q, want %q", attrs[0].Key, "code")
	}
	if attrs[0].Value.String() != "not_found" {
		t.Fatalf("attrs[0].Value: got %q, want %q", attrs[0].Value.String(), "not_found")
	}
	if attrs[1].Key != "message" {
		t.Fatalf("attrs[1].Key: got %q, want %q", attrs[1].Key, "message")
	}
	if attrs[2].Key != "id" {
		t.Fatalf("attrs[2].Key: got %q, want %q", attrs[2].Key, "id")
	}
}

func TestLogValue_NoAttrs(t *testing.T) {
	lv := apierr.Internal("server error").LogValue()

	if len(lv.Group()) != 2 { // code + message only
		t.Fatalf("group len: got %d, want 2", len(lv.Group()))
	}
}

func TestLogAttrs_DirectError(t *testing.T) {
	err := apierr.NotFound("user not found", slog.String("id", "123"))
	attrs := apierr.LogAttrs(err)

	if len(attrs) != 3 {
		t.Fatalf("len: got %d, want 3", len(attrs))
	}
	if attrs[0].Key != "error.code" {
		t.Fatalf("attrs[0].Key: got %q, want %q", attrs[0].Key, "error.code")
	}
	if attrs[0].Value.String() != "not_found" {
		t.Fatalf("attrs[0].Value: got %q, want %q", attrs[0].Value.String(), "not_found")
	}
	if attrs[1].Key != "error.message" {
		t.Fatalf("attrs[1].Key: got %q, want %q", attrs[1].Key, "error.message")
	}
	if attrs[2].Key != "id" {
		t.Fatalf("attrs[2].Key: got %q, want %q", attrs[2].Key, "id")
	}
}

func TestLogAttrs_WrappedError(t *testing.T) {
	inner := apierr.NotFound("not found")
	outer := apierr.Wrap(inner, apierr.CodeInternal, "operation failed")
	attrs := apierr.LogAttrs(outer)

	// errors.As finds outermost *Error first
	if attrs[0].Value.String() != "internal" {
		t.Fatalf("code: got %q, want %q", attrs[0].Value.String(), "internal")
	}
}

func TestLogAttrs_NoAttrs(t *testing.T) {
	attrs := apierr.LogAttrs(apierr.Internal("err"))

	if len(attrs) != 2 { // error.code + error.message only
		t.Fatalf("len: got %d, want 2", len(attrs))
	}
}

func TestLogAttrs_PlainError(t *testing.T) {
	attrs := apierr.LogAttrs(errors.New("plain error"))

	if len(attrs) != 1 {
		t.Fatalf("len: got %d, want 1", len(attrs))
	}
	if attrs[0].Key != "error" {
		t.Fatalf("Key: got %q, want %q", attrs[0].Key, "error")
	}
	if attrs[0].Value.String() != "plain error" {
		t.Fatalf("Value: got %q, want %q", attrs[0].Value.String(), "plain error")
	}
}
