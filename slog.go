package apierr

import (
	"errors"
	"log/slog"
)

// LogValue implements [slog.LogValuer], enabling *Error to be passed directly
// to slog calls as a structured group with "code", "message", and any attached attributes.
func (e *Error) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, 2+len(e.attrs))
	attrs = append(attrs,
		slog.String("code", e.code.String()),
		slog.String("message", e.message),
	)
	attrs = append(attrs, e.attrs...)
	return slog.GroupValue(attrs...)
}

// LogAttrs returns a flat []slog.Attr slice describing err, suitable for spreading
// into a slog call: slog.Info("request failed", apierr.LogAttrs(err)...).
//
// If err is (or wraps) an *Error, the result contains "error.code", "error.message",
// and any slog attributes attached to the *Error.
// For plain errors a single "error" string attribute is returned.
func LogAttrs(err error) []slog.Attr {
	var e *Error
	if errors.As(err, &e) {
		attrs := make([]slog.Attr, 0, 2+len(e.attrs))
		attrs = append(attrs,
			slog.String("error.code", e.code.String()),
			slog.String("error.message", e.message),
		)
		return append(attrs, e.attrs...)
	}
	return []slog.Attr{slog.String("error", err.Error())}
}
