package handlers_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/services/ldap"

	"net/http"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestSignInByLdapHandler_Correct(t *testing.T) {
	setupBus(t)

	server := mock.NewServer()
	code, _ := server.
		OnTenant(mock.AvengersTenant).
		ExecutePostAsJSON(handlers.SignInByLdap(), `{"provider": "ldap_test", "username": "developer", "password": "developer_password"}`)

	Expect(code).Equals(http.StatusOK)
}

func TestSignInByLdapHandler_MissingParams(t *testing.T) {
	setupBus(t)

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.AvengersTenant).
		ExecutePostAsJSON(handlers.SignInByLdap(), `{}`)

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("errors[0].message")).Equals("Username is required.")
}

func TestSignInByLdapHandler_MissingPassword(t *testing.T) {
	setupBus(t)

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.AvengersTenant).
		ExecutePostAsJSON(handlers.SignInByLdap(), `{"provider": "ldap_test", "username": "developer"}`)

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("errors[0].message")).Equals("Password is required.")
}

func TestSignInByLdapHandler_IncorrectProvider(t *testing.T) {
	setupBus(t)

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.AvengersTenant).
		ExecutePostAsJSON(handlers.SignInByLdap(), `{"provider": "incorrect", "username": "developer", "password": "developer_password"}`)

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("errors[0].message")).Equals("Username or password is invalid.")
}

func TestSignInByLdapHandler_IncorrectUser(t *testing.T) {
	setupBus(t)

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.AvengersTenant).
		ExecutePostAsJSON(handlers.SignInByLdap(), `{"provider": "ldap_test", "username": "incorrect", "password": "developer_password"}`)

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("errors[0].message")).Equals("Username or password is invalid.")
}

func TestSignInByLdapHandler_IncorrectPassword(t *testing.T) {
	setupBus(t)

	server := mock.NewServer()
	code, response := server.
		OnTenant(mock.AvengersTenant).
		ExecutePostAsJSON(handlers.SignInByLdap(), `{"provider": "ldap_test", "username": "developer", "password": "incorrect"}`)

	Expect(code).Equals(http.StatusBadRequest)
	Expect(response.String("errors[0].message")).Equals("Username or password is invalid.")
}

func setupBus(t *testing.T) {
	RegisterT(t)
	bus.Init(&ldap.Service{})

	ldapProvider := &entity.LdapConfig{
		ID:                    1,
		Provider:              "ldap_test",
		DisplayName:           "Testing Ldap Server",
		Status:                2,
		Protocol:              1,
		CertCheck:             false,
		LdapHostname:          "localhost",
		LdapPort:              "389",
		BindUsername:          "cn=readonly,dc=example,dc=org",
		BindPassword:          "readonly_password",
		RootDN:                "dc=example,dc=org",
		Scope:                 3,
		UserSearchFilter:      "(objectClass=inetOrgPerson)",
		UsernameLdapAttribute: "uid",
		NameLdapAttribute:     "displayName",
		MailLdapAttribute:     "mail",
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCustomLdapConfigByProvider) error {
		if q.Provider == "ldap_test" {
			q.Result = ldapProvider
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *cmd.RegisterUser) error {
		return nil
	})
}
