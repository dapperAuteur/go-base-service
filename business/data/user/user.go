// Package user contains user related CRUD functionality.
package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/dapperauteur/go-base-service/business/auth"
	"github.com/dapperauteur/go-base-service/foundation/database"
	"go.opentelemetry.io/otel/trace"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempted action is not allowed")
)

// User manages the set of API's for user access.
type User struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a user for api access.
func New(log *log.Logger, db *sqlx.DB) User {
	return User{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (u User) Create(ctx context.Context, traceID string, nu NewUser, now time.Time) (Info, error) {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.create")
	defer span.End()

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return Info{}, errors.Wrap(err, "generating password hash")
	}

	usr := Info{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)`

	u.log.Printf("%s : %s : QUERY : %s", traceID, "user.Create",
		database.Log(q, usr.ID, usr.Name, usr.Email, usr.PasswordHash, usr.Roles, usr.DateCreated, usr.DateUpdated),
	)

	if _, err := u.db.ExecContext(ctx, q, usr.ID, usr.Name, usr.Email, usr.PasswordHash, usr.Roles, usr.DateCreated, usr.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting user")
	}

	return usr, nil
}

// Update replaces a user document in the database.
func (u User) Update(ctx context.Context, traceID string, claims auth.Claims, userID string, uu UpdateUser, now time.Time) error {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.update")
	defer span.End()

	usr, err := u.QueryByID(ctx, traceID, claims, userID)
	if err != nil {
		// return errors.Wrap(err, "updating user")
		return err
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "generating password hash")
		}
		usr.PasswordHash = pw
	}
	usr.DateUpdated = now

	const q = `
	UPDATE
		users
	SET 
		"name" = $2,
		"email" = $3,
		"roles" = $4,
		"password_hash" = $5,
		"date_updated" = $6
	WHERE
		user_id = $1`

	u.log.Printf("%s : %s : QUERY : %s", traceID, "user.Update",
		database.Log(q, usr.ID, usr.Name, usr.Email, usr.PasswordHash, usr.Roles, usr.DateCreated, usr.DateUpdated),
	)

	if _, err := u.db.ExecContext(ctx, q, usr.ID, usr.Name, usr.Email, usr.PasswordHash, usr.Roles, usr.DateUpdated); err != nil {
		return errors.Wrapf(err, "updating user %s", usr.ID)
	}

	return nil
}

// Delete removes a user from the database.
func (u User) Delete(ctx context.Context, traceID string, userID string) error {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.delete")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = $1`

	u.log.Printf("%s: %s : QUERY : %s", traceID, "user.Delete",
		database.Log(q, userID),
	)

	if _, err := u.db.ExecContext(ctx, q, userID); err != nil {
		return errors.Wrapf(err, "deleting user %s", userID)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (u User) Query(ctx context.Context, traceID string, pageNumber int, rowsPerPage int) ([]Info, error) {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.query")
	defer span.End()

	const q = `
	SELECT
		*
	FROM
		users
	ORDER BY
		user_id
	OFFSET $1
	ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	u.log.Printf("%s : %s : QUERY : %s", traceID, "user.Query",
		database.Log(q),
	)

	users := []Info{}
	if err := u.db.SelectContext(ctx, &users, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}

	return users, nil
}

// QueryByID gets the specified user from the database.
func (u User) QueryByID(ctx context.Context, traceID string, claims auth.Claims, userID string) (Info, error) {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.queryByID")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return Info{}, ErrInvalidID
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorize(auth.RoleAdmin) && claims.Subject != userID {
		return Info{}, ErrForbidden
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = $1`

	u.log.Printf("%s : %s : QUERY : %s", traceID, "user.QueryByID",
		database.Log(q, userID),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting user %q", userID)
	}

	return usr, nil

}

// QueryByEmail gets the specified user from the database by email.
func (u User) QueryByEmail(ctx context.Context, traceID string, claims auth.Claims, email string) (Info, error) {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.queryByEmail")
	defer span.End()

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = $1`

	u.log.Printf("%s : %s : QUERY : %s", traceID, "user.QueryByEmail",
		database.Log(q, email),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, email); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting user %q", email)
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorize(auth.RoleAdmin) && claims.Subject != usr.ID {
		return Info{}, ErrForbidden
	}

	return usr, nil
}

// Authenticate finds a user by their email and verifies their password.
// On success it returns a Claims User representing this user.
// The claims can be used to generate a token for future authentication.
func (u User) Authenticate(ctx context.Context, traceID string, now time.Time, email, password string) (auth.Claims, error) {

	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "business.data.user.authenticate")
	defer span.End()

	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`

	u.log.Printf("%s : %s : QUERY : %s", traceID, "user.Authenticate",
		database.Log(q, data),
	)

	var usr Info
	if err := database.NamedQueryStruct(ctx, u.db, q, data, &usr); err != nil {
		if err == database.ErrNotFound {
			return auth.Claims{}, database.ErrNotFound
		}
		return auth.Claims{}, errors.Wrapf(err, "selecting user %q", email)
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, database.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   usr.ID,
			ExpiresAt: now.Add(time.Hour).Unix(),
			IssuedAt:  now.Unix(),
		},
		Roles: usr.Roles,
	}

	return claims, nil
}
