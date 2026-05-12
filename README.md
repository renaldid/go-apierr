# go-apierr

**Structured API errors for Go — HTTP status codes, gRPC codes, slog attributes, and RFC 9457 Problem Details in one ergonomic zero-dependency type.**

[![Go Reference](https://pkg.go.dev/badge/github.com/renaldid/go-apierr.svg)](https://pkg.go.dev/github.com/renaldid/go-apierr)
[![CI](https://github.com/renaldid/go-apierr/actions/workflows/ci.yml/badge.svg)](https://github.com/renaldid/go-apierr/actions/workflows/ci.yml)
[![Coverage](https://codecov.io/gh/renaldid/go-apierr/badge.svg)](https://codecov.io/gh/renaldid/go-apierr)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## The Problem

Every Go backend service eventually needs to map errors to HTTP status codes, gRPC codes, and structured logs, but these concerns are typically scattered across the codebase:

```go
// Without go-apierr: repeated, inconsistent, hard to maintain
if errors.Is(err, ErrUserNotFound) {
    http.Error(w, "not found", 404)
    slog.Error("user missing", "id", id)
    return status.Error(codes.NotFound, "not found") // gRPC
}
```

`go-apierr` solves this by encoding all three mappings into a single `*Error` value, defined once.

## Features

- **Single `*Error` type** — carries a semantic code, message, slog attributes, and an optional wrapped cause
- **HTTP status mapping** — converts any `*Error` to the correct HTTP status code automatically
- **RFC 9457 Problem Details** — writes spec-compliant `application/problem+json` responses out of the box
- **`slog` integration** — implements `slog.LogValuer`; pass `*Error` directly to any slog call
- **`errors.Is` / `errors.As` compatible** — sentinel-style matching by `Code`
- **gRPC support** — optional [`grpcerr`](./grpcerr) sub-module maps every `Code` to a `google.golang.org/grpc/codes.Code`
- **Zero dependencies** — core module requires only the Go standard library (Go 1.21+)
- **100% test coverage** — race-detector clean

## Installation

```bash
go get github.com/renaldid/go-apierr
```

For gRPC support (separate module, optional):

```bash
go get github.com/renaldid/go-apierr/grpcerr
```

## Quick Start

```go
package main

import (
    "log/slog"
    "net/http"

    "github.com/renaldid/go-apierr"
)

func getUser(w http.ResponseWriter, r *http.Request) {
    userID := r.PathValue("id")
    user, err := db.FindUser(r.Context(), userID)
    if err != nil {
        apierr.WriteHTTP(w, r, apierr.NotFound("user not found",
            slog.String("user_id", userID),
        ))
        return
    }
    // ...
}
```

The response body:

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "user not found",
  "instance": "/users/42",
  "user_id": "42"
}
```

## Usage

### Creating Errors

```go
// Convenience constructors — one per semantic code
err := apierr.NotFound("user not found", slog.String("id", id))
err := apierr.InvalidArgument("email is required", slog.String("field", "email"))
err := apierr.Unauthenticated("token expired")
err := apierr.Internal("unexpected db error")

// Generic constructor
err := apierr.New(apierr.CodeNotFound, "product not found",
    slog.String("sku", sku),
    slog.String("warehouse", wh),
)

// Formatted message
err := apierr.Newf(apierr.CodeNotFound, "user %s not found", userID)
```

### Wrapping Errors

```go
row, dbErr := pool.QueryRow(ctx, query, id)
if dbErr != nil {
    return apierr.Wrap(dbErr, apierr.CodeInternal, "failed to fetch user",
        slog.String("user_id", id),
    )
}
```

### Sentinel Matching

```go
var ErrNotFound = apierr.New(apierr.CodeNotFound, "")

// Matches any *Error with CodeNotFound, regardless of message
if errors.Is(err, ErrNotFound) {
    // handle not found
}

// Or check the code directly
if apierr.CodeOf(err) == apierr.CodeNotFound {
    // handle not found
}
```

### HTTP Handlers (REST API)

```go
// WriteHTTP sets Content-Type: application/problem+json and the correct status code
apierr.WriteHTTP(w, r, err)

// Convert to Problem struct for custom rendering
p := apierr.ToProblem(err)
p.Instance = "/custom/path"
json.NewEncoder(w).Encode(p)
```

### Structured Logging with slog

```go
// *Error implements slog.LogValuer — use it directly
slog.Error("request failed", "err", err)
// -> {"err":{"code":"not_found","message":"user not found","id":"42"}}

// Or spread flat attributes
slog.Error("request failed", apierr.LogAttrs(err)...)
// -> {"error.code":"not_found","error.message":"user not found","id":"42"}
```

### gRPC Handlers

Import the optional [`grpcerr`](./grpcerr) sub-module:

```go
import "github.com/renaldid/go-apierr/grpcerr"

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := s.db.FindUser(ctx, req.Id)
    if err != nil {
        return nil, grpcerr.ToGRPC(err) // -> gRPC NotFound status
    }
    return toProto(user), nil
}
```

Converting gRPC errors back to `*Error`:

```go
err := grpcerr.FromGRPC(grpcErr) // gRPC status -> *apierr.Error
```

## Error Codes and Status Mappings

| Code | String | HTTP Status | gRPC Code |
|------|--------|-------------|-----------|
| `CodeUnknown` | `unknown` | 500 | `Unknown` |
| `CodeCanceled` | `canceled` | 408 | `Canceled` |
| `CodeInvalidArgument` | `invalid_argument` | 400 | `InvalidArgument` |
| `CodeDeadlineExceeded` | `deadline_exceeded` | 504 | `DeadlineExceeded` |
| `CodeNotFound` | `not_found` | 404 | `NotFound` |
| `CodeAlreadyExists` | `already_exists` | 409 | `AlreadyExists` |
| `CodePermissionDenied` | `permission_denied` | 403 | `PermissionDenied` |
| `CodeResourceExhausted` | `resource_exhausted` | 429 | `ResourceExhausted` |
| `CodeFailedPrecondition` | `failed_precondition` | 400 | `FailedPrecondition` |
| `CodeAborted` | `aborted` | 409 | `Aborted` |
| `CodeUnimplemented` | `unimplemented` | 501 | `Unimplemented` |
| `CodeInternal` | `internal` | 500 | `Internal` |
| `CodeUnavailable` | `unavailable` | 503 | `Unavailable` |
| `CodeUnauthenticated` | `unauthenticated` | 401 | `Unauthenticated` |

## RFC 9457 Problem Details

`WriteHTTP` and `ToProblem` produce [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) compliant responses. slog attributes attached to the error are serialized as extension members at the top level:

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "user not found",
  "instance": "/api/v1/users/42",
  "user_id": "42",
  "trace_id": "01HXK..."
}
```

## Contributing

Contributions, issues, and feature requests are welcome. Please open an issue before submitting a pull request.

```bash
git clone https://github.com/renaldid/go-apierr
cd go-apierr
go test -race ./...
```

## License

MIT — see [LICENSE](LICENSE).
