package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperauteur/go-base-service/foundation/web"
)

// Logger ...
func Logger(log *log.Logger) web.Middleware {
	// readiness from check.go
	m := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// readiness from check.go
			err := before(ctx, w, r)

			// BOILERPLATE - LOGGING
			log.Println(r)

			return err
		}
		return h
	}
	return m
}
