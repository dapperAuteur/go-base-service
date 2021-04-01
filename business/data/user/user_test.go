package user_test

import (
	"testing"
	"time"

	"github.com/dapperauteur/go-base-service/business/auth"
	"github.com/dapperauteur/go-base-service/business/data/user"
	"github.com/dapperauteur/go-base-service/business/tests"
)

func TestUser(t *testing.T) {
	log, db, teardown := tests.NewUnit(t, dbc)
	t.Cleanup(teardown)

	u := user.New(log, db)

	t.Log("Given the need to work with User records.")
	{
		ctx := tests.Context()
		now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
		traceID := "00000000-0000-0000-0000-000000000000"

		nu := user.NewUser{
			Name:            "awe ful",
			Email:           "aweful@awews.com",
			Roles:           []string{auth.RoleAdmin},
			Password:        "spread love",
			PasswordConfirm: "spread love",
		}
	}
}
