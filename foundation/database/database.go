// Package database provides support for access to the database.
package database

import (
	"time"

	_ "github.com/lib/pq" // Calls init function.

	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sqlx.DB, error) {
	// Disabling SSL to simplify the dev process.
	// Shouldn't do this for prod env.
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	// Postgres specific. It's url based.
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	// Construct URL needed for connectivity string.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}