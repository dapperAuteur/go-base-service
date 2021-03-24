package mid

import (
	"context"
	"errors"
	"net/http"

	"github.com/dapperauteur/go-base-service/business/auth"
	"github.com/dapperauteur/go-base-service/foundation/web"
)

// ErrForbidden is returned when an authenticated user does not have a sufficient role for an action
var ErrForbidden = web.NewRequestError(
	errors.New("you are not authorized for that action"),
	http.StatusForbidden,
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return nil
		}
		return h
	}
	return m
}

// Authorize validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
func Authorize(roles ...string) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(after web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return nil
		}
		return h
	}
	return m
}
