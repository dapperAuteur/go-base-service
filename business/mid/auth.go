package mid

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dapperauteur/go-base-service/business/auth"
	"github.com/dapperauteur/go-base-service/foundation/web"
	"go.opentelemetry.io/otel/trace"
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

			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.mid.authenticate")
			defer span.End()

			// Expecting: bearer <token>
			authStr := r.Header.Get("authorization")

			// Parse the authorization header.
			parts := strings.Split(authStr, " ")

			// Expected header is of the format `Bearer <token>`.
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				// log.Fatalln(err)
				return web.NewRequestError(err, http.StatusUnauthorized)

			}

			// Validate the token is signed by us.
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			// Add claims to the context so they can be retrieved later.
			ctx = context.WithValue(ctx, auth.Key, claims)

			// Call the next handler.
			return handler(ctx, w, r)
		}
		return h
	}
	return m
}

// Authorize validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
func Authorize(log *log.Logger, roles ...string) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.mid.authorize")
			defer span.End()

			// If the context is missing this value return failure.
			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from context")
			}

			if !claims.Authorize(roles...) {
				log.Printf("mid: authorize: claims: %v exp: %v", claims.Roles, roles)
				// return validate.NewRequestError(
				// 	fmt.Errorf("you are not authorized for that action: claims: %v exp: %v", claims.Roles, roles),
				// 	http.StatusForbidden,
				// )
				return ErrForbidden
			}

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
