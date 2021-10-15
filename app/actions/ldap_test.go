package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/rand"
)

func TestCreateEditLdapConfig_Validate_InvalidInput(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		expected []string
		action   *actions.CreateEditLdapConfig
	}{
		{
			expected: []string{"status", "protocol", "scope", "displayName", "ldapHostname", "ldapPort", "bindUsername", "bindPassword", "rootDN", "userSearchFilter", "usernameLdapAttribute", "nameLdapAttribute", "mailLdapAttribute"},
			action:   &actions.CreateEditLdapConfig{},
		},
		{
			expected: []string{"status", "protocol", "scope", "displayName", "ldapHostname", "ldapPort", "bindUsername", "bindPassword", "rootDN", "userSearchFilter", "usernameLdapAttribute", "nameLdapAttribute", "mailLdapAttribute"},
			action: &actions.CreateEditLdapConfig{
				ID:                    0,
				Provider:              "",
				DisplayName:           rand.String(51),
				Status:                0,
				Protocol:              0,
				LdapHostname:          rand.String(301),
				LdapPort:              "12345678910",
				BindUsername:          rand.String(101),
				BindPassword:          rand.String(101),
				RootDN:                rand.String(251),
				Scope:                 0,
				UserSearchFilter:      rand.String(501),
				UsernameLdapAttribute: rand.String(101),
				NameLdapAttribute:     rand.String(101),
				MailLdapAttribute:     rand.String(101),
			},
		},
		{
			expected: []string{"ldapPort"},
			action: &actions.CreateEditLdapConfig{
				ID:                    0,
				Provider:              "",
				DisplayName:           "Test",
				Status:                enum.LdapConfigEnabled,
				Protocol:              enum.LDAP,
				LdapHostname:          "Hostname",
				LdapPort:              "Invalid",
				BindUsername:          "Bind Username",
				BindPassword:          "Bind Password",
				RootDN:                "Root DN",
				Scope:                 enum.ScopeBaseObject,
				UserSearchFilter:      "User Search Filter",
				UsernameLdapAttribute: "Username LDAP Attribute",
				NameLdapAttribute:     "Name LDAP Attribute",
				MailLdapAttribute:     "Mail LDAP Attribute",
			},
		},
	}

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{
		IsEmailAuthAllowed: true,
	})

	for _, testCase := range testCases {
		result := testCase.action.Validate(ctx, nil)
		ExpectFailed(result, testCase.expected...)
	}
}

func TestCreateEditLdapConfig_Validate_ValidInput(t *testing.T) {
	RegisterT(t)

	action := &actions.CreateEditLdapConfig{
		ID:                    0,
		Provider:              "",
		DisplayName:           "Test",
		Status:                enum.LdapConfigEnabled,
		Protocol:              enum.LDAP,
		LdapHostname:          "Hostname",
		LdapPort:              "1234",
		BindUsername:          "Bind Username",
		BindPassword:          "Bind Password",
		RootDN:                "Root DN",
		Scope:                 enum.ScopeBaseObject,
		UserSearchFilter:      "User Search Filter",
		UsernameLdapAttribute: "Username LDAP Attribute",
		NameLdapAttribute:     "Name LDAP Attribute",
		MailLdapAttribute:     "Mail LDAP Attribute",
	}

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{
		IsEmailAuthAllowed: true,
	})

	result := action.Validate(ctx, nil)
	ExpectSuccess(result)
}

func TestCreateEditLdapConfig_DefaultValues(t *testing.T) {
	RegisterT(t)

	action := &actions.CreateEditLdapConfig{}
	Expect(action.ID).Equals(0)
}

func TestCreateEditLdapConfig_IsAuthorized(t *testing.T) {
	RegisterT(t)

	action := &actions.CreateEditLdapConfig{}
	Expect(action.IsAuthorized(context.Background(), nil)).IsFalse()
	Expect(action.IsAuthorized(context.Background(), &entity.User{})).IsFalse()
	Expect(action.IsAuthorized(context.Background(), &entity.User{Role: enum.RoleAdministrator})).IsTrue()
}
