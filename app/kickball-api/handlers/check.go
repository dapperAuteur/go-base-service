package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/dapperauteur/go-base-service/foundation/web"
)

type check struct {
	log *log.Logger
}

func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	// to simulate ERRORS and PANICS
	// trusted error
	// if n := rand.Intn(100); n%2 == 0 {
	// return web.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	// untrusted error
	// return errors.New("untrusted error")
	// force panic
	// panic("forcing panic")
	// force shutdown
	// return web.NewShutdownError("forcing shutdown")
	// }

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
