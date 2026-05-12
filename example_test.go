package apierr_test

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/renaldid/go-apierr"
)

func ExampleNotFound() {
	err := apierr.NotFound("user not found")
	fmt.Println(err)
	// Output:
	// user not found
}

func ExampleNew_withAttrs() {
	err := apierr.New(apierr.CodeNotFound, "user not found",
		slog.String("user_id", "abc123"),
	)
	fmt.Println(err.Code())
	fmt.Println(err.Message())
	// Output:
	// not_found
	// user not found
}

func ExampleWrap() {
	cause := errors.New("connection refused")
	err := apierr.Wrap(cause, apierr.CodeUnavailable, "service unavailable")
	fmt.Println(err)
	// Output:
	// service unavailable: connection refused
}

func ExampleCodeOf() {
	err := apierr.NotFound("resource missing")
	fmt.Println(apierr.CodeOf(err))
	// Output:
	// not_found
}

func ExampleCode_HTTPStatus() {
	fmt.Println(apierr.CodeNotFound.HTTPStatus())
	fmt.Println(apierr.CodeUnauthenticated.HTTPStatus())
	fmt.Println(apierr.CodeResourceExhausted.HTTPStatus())
	// Output:
	// 404
	// 401
	// 429
}
