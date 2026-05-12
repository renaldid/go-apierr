package apierr

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

const problemContentType = "application/problem+json"

// Problem is an RFC 9457 Problem Details object.
// Extension members in Extra are merged at the top level when marshaled to JSON.
type Problem struct {
	Type     string
	Title    string
	Status   int
	Detail   string
	Instance string
	Extra    map[string]any
}

// MarshalJSON serializes p as an RFC 9457 JSON object.
// Detail and Instance are omitted when empty. Extra fields are merged at the top level.
func (p Problem) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, 5+len(p.Extra))
	m["type"] = p.Type
	m["title"] = p.Title
	m["status"] = p.Status
	if p.Detail != "" {
		m["detail"] = p.Detail
	}
	if p.Instance != "" {
		m["instance"] = p.Instance
	}
	for k, v := range p.Extra {
		m[k] = v
	}
	return json.Marshal(m)
}

// ToProblem converts err into an RFC 9457 Problem.
// If err is (or wraps) an *Error, Problem fields are populated from it.
// For any other error a generic 500 Internal Server Error Problem is returned.
func ToProblem(err error) Problem {
	p := Problem{
		Type:   "about:blank",
		Status: http.StatusInternalServerError,
		Title:  http.StatusText(http.StatusInternalServerError),
	}

	var e *Error
	if errors.As(err, &e) {
		p.Status = e.code.HTTPStatus()
		p.Title = http.StatusText(p.Status)
		p.Detail = e.message
		if len(e.attrs) > 0 {
			p.Extra = make(map[string]any, len(e.attrs))
			for _, a := range e.attrs {
				p.Extra[a.Key] = attrValue(a.Value)
			}
		}
	}
	return p
}

// WriteHTTP writes err as an RFC 9457 Problem Details JSON response.
// Content-Type is set to "application/problem+json".
// If r is non-nil, Problem.Instance is set to r.URL.RequestURI().
func WriteHTTP(w http.ResponseWriter, r *http.Request, err error) {
	p := ToProblem(err)
	if r != nil {
		p.Instance = r.URL.RequestURI()
	}
	w.Header().Set("Content-Type", problemContentType)
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(p)
}

// attrValue converts a slog.Value to a JSON-compatible Go primitive.
// Numeric and boolean kinds are preserved; all others fall back to their string representation.
func attrValue(v slog.Value) any {
	rv := v.Resolve()
	switch rv.Kind() {
	case slog.KindInt64:
		return rv.Int64()
	case slog.KindUint64:
		return rv.Uint64()
	case slog.KindFloat64:
		return rv.Float64()
	case slog.KindBool:
		return rv.Bool()
	default:
		return rv.String()
	}
}
