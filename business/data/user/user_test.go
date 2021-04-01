package user_test

import (
	"testing"
	"time"

	"github.com/dapperauteur/go-base-service/business/auth"
	"github.com/dapperauteur/go-base-service/business/data/user"
	"github.com/dapperauteur/go-base-service/business/tests"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp"
)

func TestUser(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
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

		usr, err := u.Create(ctx, traceID, nu, now)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

		claims := auth.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "service project",
				Subject:   usr.ID,
				ExpiresAt: now.Add(time.Hour).Unix(),
				IssuedAt:  now.Unix(),
			},
			Roles: []string{auth.RoleUser},
		}

		saved, err := u.QueryByID(ctx, traceID, claims, usr.ID)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)

		if diff := cmp.Diff(usr, saved); diff != "" {
			t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", tests.Failed, testID, diff)
		}
		t.Logf("\t%s\tTest %d:\tShould get back the same user.", tests.Success, testID)
	}
}
