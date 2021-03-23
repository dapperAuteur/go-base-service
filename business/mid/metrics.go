package mid

import (
	"context"
	"expvar"
	"net/http"

	"github.com/dapperauteur/go-base-service/foundation/web"
)

// m contains the global program counters for the application.
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics updates program counters.
func Metrics() web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Call the next handler.
			err := handler(ctx, w, r)

			// Increment the request counter.
			m.req.Add(1)
		}
		return h
	}
	return m

}
