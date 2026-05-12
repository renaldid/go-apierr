package apierr_test

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/renaldid/go-apierr"
)

func TestToProblem_WithError(t *testing.T) {
	err := apierr.NotFound("user not found", slog.String("id", "123"))
	p := apierr.ToProblem(err)

	if p.Status != http.StatusNotFound {
		t.Fatalf("Status: got %d, want %d", p.Status, http.StatusNotFound)
	}
	if p.Title != "Not Found" {
		t.Fatalf("Title: got %q, want %q", p.Title, "Not Found")
	}
	if p.Detail != "user not found" {
		t.Fatalf("Detail: got %q, want %q", p.Detail, "user not found")
	}
	if p.Type != "about:blank" {
		t.Fatalf("Type: got %q, want %q", p.Type, "about:blank")
	}
}

func TestToProblem_ExtraFields(t *testing.T) {
	err := apierr.NotFound("user not found", slog.String("id", "123"))
	p := apierr.ToProblem(err)

	if len(p.Extra) != 1 {
		t.Fatalf("Extra len: got %d, want 1", len(p.Extra))
	}
	if p.Extra["id"] != "123" {
		t.Fatalf("Extra[id]: got %v, want %q", p.Extra["id"], "123")
	}
}

func TestToProblem_NoAttrs(t *testing.T) {
	p := apierr.ToProblem(apierr.NotFound("not found"))

	if p.Extra != nil {
		t.Fatal("Extra should be nil when error has no attrs")
	}
}

func TestToProblem_PlainError(t *testing.T) {
	p := apierr.ToProblem(errors.New("unexpected failure"))

	if p.Status != http.StatusInternalServerError {
		t.Fatalf("Status: got %d, want %d", p.Status, http.StatusInternalServerError)
	}
	if p.Detail != "" {
		t.Fatalf("Detail should be empty for plain error, got %q", p.Detail)
	}
}

func TestToProblem_WrappedError(t *testing.T) {
	inner := apierr.NotFound("resource missing")
	outer := apierr.Wrap(inner, apierr.CodeInternal, "operation failed")
	p := apierr.ToProblem(outer)

	// errors.As finds the outermost *Error (outer = CodeInternal)
	if p.Status != http.StatusInternalServerError {
		t.Fatalf("Status: got %d, want %d", p.Status, http.StatusInternalServerError)
	}
}

func TestAttrValue_Int64(t *testing.T) {
	p := apierr.ToProblem(apierr.New(apierr.CodeInternal, "m", slog.Int64("n", 42)))
	if v, ok := p.Extra["n"].(int64); !ok || v != 42 {
		t.Fatalf("Extra[n]: got %T(%v), want int64(42)", p.Extra["n"], p.Extra["n"])
	}
}

func TestAttrValue_Uint64(t *testing.T) {
	p := apierr.ToProblem(apierr.New(apierr.CodeInternal, "m", slog.Uint64("n", 100)))
	if v, ok := p.Extra["n"].(uint64); !ok || v != 100 {
		t.Fatalf("Extra[n]: got %T(%v), want uint64(100)", p.Extra["n"], p.Extra["n"])
	}
}

func TestAttrValue_Float64(t *testing.T) {
	p := apierr.ToProblem(apierr.New(apierr.CodeInternal, "m", slog.Float64("r", 0.5)))
	if v, ok := p.Extra["r"].(float64); !ok || v != 0.5 {
		t.Fatalf("Extra[r]: got %T(%v), want float64(0.5)", p.Extra["r"], p.Extra["r"])
	}
}

func TestAttrValue_Bool(t *testing.T) {
	p := apierr.ToProblem(apierr.New(apierr.CodeInternal, "m", slog.Bool("ok", true)))
	if v, ok := p.Extra["ok"].(bool); !ok || !v {
		t.Fatalf("Extra[ok]: got %T(%v), want bool(true)", p.Extra["ok"], p.Extra["ok"])
	}
}

func TestAttrValue_StringDefault(t *testing.T) {
	p := apierr.ToProblem(apierr.New(apierr.CodeInternal, "m", slog.String("k", "val")))
	if v, ok := p.Extra["k"].(string); !ok || v != "val" {
		t.Fatalf("Extra[k]: got %T(%v), want string(val)", p.Extra["k"], p.Extra["k"])
	}
}

func TestAttrValue_TimeDefault(t *testing.T) {
	p := apierr.ToProblem(apierr.New(apierr.CodeInternal, "m", slog.Time("ts", time.Now())))
	if _, ok := p.Extra["ts"].(string); !ok {
		t.Fatalf("Extra[ts]: got %T, want string", p.Extra["ts"])
	}
}

func TestAttrValue_DurationDefault(t *testing.T) {
	p := apierr.ToProblem(apierr.New(apierr.CodeInternal, "m", slog.Duration("dur", time.Second)))
	if _, ok := p.Extra["dur"].(string); !ok {
		t.Fatalf("Extra[dur]: got %T, want string", p.Extra["dur"])
	}
}

func TestProblem_MarshalJSON_AllFields(t *testing.T) {
	p := apierr.Problem{
		Type:     "about:blank",
		Title:    "Not Found",
		Status:   404,
		Detail:   "user not found",
		Instance: "/users/123",
		Extra:    map[string]any{"user_id": "abc"},
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	checks := map[string]any{
		"type":     "about:blank",
		"title":    "Not Found",
		"status":   float64(404),
		"detail":   "user not found",
		"instance": "/users/123",
		"user_id":  "abc",
	}
	for k, want := range checks {
		if m[k] != want {
			t.Errorf("JSON[%q]: got %v, want %v", k, m[k], want)
		}
	}
}

func TestProblem_MarshalJSON_OptionalFieldsOmitted(t *testing.T) {
	p := apierr.Problem{Type: "about:blank", Title: "OK", Status: 200}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if _, ok := m["detail"]; ok {
		t.Error("detail must be omitted when empty")
	}
	if _, ok := m["instance"]; ok {
		t.Error("instance must be omitted when empty")
	}
}

func TestWriteHTTP_WithRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/42?v=1", nil)

	apierr.WriteHTTP(rec, req, apierr.NotFound("user not found"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Fatalf("Content-Type: got %q, want %q", ct, "application/problem+json")
	}
	var m map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatalf("Decode body: %v", err)
	}
	if m["instance"] != "/users/42?v=1" {
		t.Fatalf("instance: got %v, want %q", m["instance"], "/users/42?v=1")
	}
}

func TestWriteHTTP_NilRequest(t *testing.T) {
	rec := httptest.NewRecorder()

	apierr.WriteHTTP(rec, nil, apierr.Internal("server error"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	var m map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatalf("Decode body: %v", err)
	}
	if _, exists := m["instance"]; exists {
		t.Error("instance must not appear when request is nil")
	}
}

func TestWriteHTTP_PlainError(t *testing.T) {
	rec := httptest.NewRecorder()

	apierr.WriteHTTP(rec, nil, errors.New("unexpected"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
