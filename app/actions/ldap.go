package actions

import (
	"context"
	"strconv"
	"strings"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/validate"
)

// IsInteger verifies if string is an integer
func IsInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// CreateEditLdapConfig is used to create/edit LDAP configuration
type CreateEditLdapConfig struct {
	ID                    int
	Provider              string `json:"provider"`
	DisplayName           string `json:"displayName"`
	Status                int    `json:"status"`
	Protocol              int    `json:"protocol"`
	CertCheck             bool   `json:"certCheck"`
	LdapHostname          string `json:"ldapHostname"`
	LdapPort              string `json:"ldapPort"`
	BindUsername          string `json:"bindUsername"`
	BindPassword          string `json:"bindPassword"`
	RootDN                string `json:"rootDN"`
	Scope                 int    `json:"scope"`
	UserSearchFilter      string `json:"userSearchFilter"`
	UsernameLdapAttribute string `json:"usernameLdapAttribute"`
	NameLdapAttribute     string `json:"nameLdapAttribute"`
	MailLdapAttribute     string `json:"mailLdapAttribute"`
}

func NewCreateEditLdapConfig() *CreateEditLdapConfig {
	return &CreateEditLdapConfig{}
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateEditLdapConfig) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *CreateEditLdapConfig) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Provider != "" {
		getConfig := &query.GetCustomLdapConfigByProvider{Provider: action.Provider}
		err := bus.Dispatch(ctx, getConfig)
		if err != nil {
			return validate.Error(err)
		}
		action.ID = getConfig.Result.ID
		if action.BindPassword == "" {
			action.BindPassword = getConfig.Result.BindPassword
		}

	} else {
		action.Provider = "_" + strings.ToLower(rand.String(10))
	}

	if action.Status != enum.LdapConfigEnabled &&
		action.Status != enum.LdapConfigDisabled {
		result.AddFieldFailure("status", "Invalid status.")
	}

	if action.Protocol != enum.LDAP &&
		action.Protocol != enum.LDAPTLS &&
		action.Protocol != enum.LDAPS {
		result.AddFieldFailure("protocol", "Invalid Protocol status.")
	}

	if action.Scope != enum.ScopeBaseObject &&
		action.Scope != enum.ScopeSingleLevel &&
		action.Scope != enum.ScopeWholeSubtree {
		result.AddFieldFailure("scope", "Invalid scope status.")
	}

	if action.DisplayName == "" {
		result.AddFieldFailure("displayName", "Display Name is required.")
	} else if len(action.DisplayName) > 50 {
		result.AddFieldFailure("displayName", "Display Name must have less than 50 characters.")
	}

	if action.LdapHostname == "" {
		result.AddFieldFailure("ldapHostname", "LDAP Domain is required.")
	} else if len(action.LdapHostname) > 300 {
		result.AddFieldFailure("ldapHostname", "LDAP Domain must have less than 300 characters.")
	}

	if action.LdapPort == "" {
		result.AddFieldFailure("ldapPort", "LDAP port is required.")
	} else if len(action.LdapPort) > 10 {
		result.AddFieldFailure("ldapPort", "LDAP port must be less than 10 digits.")
	} else if !IsInteger(action.LdapPort) {
		result.AddFieldFailure("ldapPort", "LDAP must be an integer")
	}

	if action.BindUsername == "" {
		result.AddFieldFailure("bindUsername", "Bind username is required.")
	} else if len(action.BindUsername) > 100 {
		result.AddFieldFailure("bindUsername", "Bind username must have less than 100 characters.")
	}

	if action.BindPassword == "" {
		result.AddFieldFailure("bindPassword", "Bind password is required.")
	} else if len(action.BindPassword) > 100 {
		result.AddFieldFailure("bindPassword", "Bind password must have less than 100 characters.")
	}

	if action.RootDN == "" {
		result.AddFieldFailure("rootDN", "Root DN is required.")
	} else if len(action.RootDN) > 250 {
		result.AddFieldFailure("rootDN", "Root DN must have less than 250 characters.")
	}

	if action.UserSearchFilter == "" {
		result.AddFieldFailure("userSearchFilter", "User Search Filter is required.")
	} else if len(action.UserSearchFilter) > 500 {
		result.AddFieldFailure("userSearchFilter", "User Search Filter must have less than 500 characters.")
	}

	if action.UsernameLdapAttribute == "" {
		result.AddFieldFailure("usernameLdapAttribute", "Username LDAP attribute is required.")
	} else if len(action.UsernameLdapAttribute) > 100 {
		result.AddFieldFailure("usernameLdapAttribute", "Username LDAP attribute must have less than 100 characters.")
	}

	if action.NameLdapAttribute == "" {
		result.AddFieldFailure("nameLdapAttribute", "Full Name LDAP attribute is required.")
	} else if len(action.NameLdapAttribute) > 100 {
		result.AddFieldFailure("nameLdapAttribute", "Full Name LDAP attribute must have less than 100 characters.")
	}

	if action.MailLdapAttribute == "" {
		result.AddFieldFailure("mailLdapAttribute", "Mail LDAP attribute is required.")
	} else if len(action.MailLdapAttribute) > 100 {
		result.AddFieldFailure("mailLdapAttribute", "Mail LDAP attribute must have less than 100 characters.")
	}

	return result
}
