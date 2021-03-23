package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperauteur/go-base-service/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			// Run the next handler and catch any propagated error.
			if err := handler(ctx, w, r); err != nil {

				// Log the error.
				log.Printf("TraceID %s: ERROR : %v", v.TraceID, err)

				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er); err != nil {
					return err
				}

			return nil
		}
		return h
	}

	return m
}
