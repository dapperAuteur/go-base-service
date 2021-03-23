// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// Handler ...
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App ...
type App struct {
	*httptreemux.ContextMux
}

// NewApp ...
func NewApp() *App {
	app := App{
		ContextMux: httptreemux.NewContextMux(),
	}
	return &app
}

// Handle
func (a *App) Handle(method string, path string, readiness Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {
		// r is the readiness function in check.go (check.readiness)
		// BOILERPLATE
		if err := readiness(r.Context(), w, r); err != nil {
			// Handle error
		}

		// BOILERPLATE
	}

	a.ContextMux.Handle(method, path, h)
}
