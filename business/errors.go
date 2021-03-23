package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperauteur/go-base-service/foundation/web"
)

// Errors ...
func Errors(log *log.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				// HANDLE ERROR
				log.Println(err)

				// HANDLE ERROR RESPONSE

				// SHUTDOWN SIGNAL?
				// return err

				// DECISION
				return nil
			}

			return nil
		}
		return h
	}

	return m
}
