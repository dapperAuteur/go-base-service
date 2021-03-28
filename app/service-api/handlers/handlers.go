// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"log"
	"net/http"
	"os"


	"github.com/jmoiron/sqlx"

	"github.com/dapperauteur/go-base-service/business/auth"
	"github.com/dapperauteur/go-base-service/business/mid"
	"github.com/dapperauteur/go-base-service/foundation/web"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) *web.App {

	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	cg := checkGroup{
		build: build,
		db:   db,
	}

	app.Handle(http.MethodGet, "/readiness", cg.readiness, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	app.Handle(http.MethodGet, "/liveness", cg.liveness, mid.Authenticate(a), mid.Authorize(log, auth.RoleAdmin))
	return app
}
