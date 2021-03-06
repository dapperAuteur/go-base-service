// Package database provides support for access to the database.
package database

import (
	"time"

	_ "github.com/lib/pq" // Calls init function. DB driver in use.

	"context"
	"errors"
	"fmt"
	"net/url"

	// "reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempted action is not allowed")
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

// NamedQueryStruct is a helper function for executing queries that return a
// single value to be unmarshalled into a struct type.
func NamedQueryStruct(ctx context.Context, db *sqlx.DB, query string, data interface{}, dest interface{}) error {
	rows, err := db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return ErrNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	// First check we can ping the database.
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// Log provides a pretty print version of the query and parameters.
func Log(query string, args ...interface{}) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	// query = strings.Replace(query, "\t", "", -1)
	// query = strings.Replace(query, "\n", " ", -1)

	// return fmt.Sprintf("[%s]\n", strings.Trim(query, " "))
	return query
}
