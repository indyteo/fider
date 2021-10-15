package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// GeneralSettingsPage is the general settings page
func GeneralSettingsPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(web.Props{
			Title:     "General · Site Settings",
			ChunkName: "GeneralSettings.page",
		})
	}
}

// AdvancedSettingsPage is the advanced settings page
func AdvancedSettingsPage() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.Page(web.Props{
			Title:     "Advanced · Site Settings",
			ChunkName: "AdvancedSettings.page",
			Data: web.Map{
				"customCSS": c.Tenant().CustomCSS,
			},
		})
	}
}

// UpdateSettings update current tenant' settings
func UpdateSettings() web.HandlerFunc {
	return func(c *web.Context) error {
		action := actions.NewUpdateTenantSettings()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c,
			&cmd.UploadImage{
				Image:  action.Logo,
				Folder: "logos",
			},
			&cmd.UpdateTenantSettings{
				Logo:           action.Logo,
				Title:          action.Title,
				Invitation:     action.Invitation,
				WelcomeMessage: action.WelcomeMessage,
				CNAME:          action.CNAME,
				Locale:         action.Locale,
			},
		); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UpdateAdvancedSettings update current tenant' advanced settings
func UpdateAdvancedSettings() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdateTenantAdvancedSettings)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, &cmd.UpdateTenantAdvancedSettings{
			CustomCSS: action.CustomCSS,
		}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UpdatePrivacy update current tenant's privacy settings
func UpdatePrivacy() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdateTenantPrivacy)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		updateSettings := &cmd.UpdateTenantPrivacySettings{
			IsPrivate: action.IsPrivate,
		}
		if err := bus.Dispatch(c, updateSettings); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UpdateEmailAuthAllowed update current tenant's allow email auth settings
func UpdateEmailAuthAllowed() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdateTenantEmailAuthAllowed)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		updateSettings := &cmd.UpdateTenantEmailAuthAllowedSettings{
			IsEmailAuthAllowed: action.IsEmailAuthAllowed,
		}
		if err := bus.Dispatch(c, updateSettings); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// ManageMembers is the page used by administrators to change member's role
func ManageMembers() web.HandlerFunc {
	return func(c *web.Context) error {
		allUsers := &query.GetAllUsers{}
		if err := bus.Dispatch(c, allUsers); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Manage Members · Site Settings",
			ChunkName: "ManageMembers.page",
			Data: web.Map{
				"users": allUsers.Result,
			},
		})
	}
}

// ManageAuthentication is the page used by administrators to change site authentication settings
func ManageAuthentication() web.HandlerFunc {
	return func(c *web.Context) error {
		listOauthProviders := &query.ListAllOAuthProviders{}
		listLdapProviders := &query.ListAllLdapProviders{}
		if err := bus.Dispatch(c, listOauthProviders, listLdapProviders); err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title:     "Authentication · Site Settings",
			ChunkName: "ManageAuthentication.page",
			Data: web.Map{
				"oauthProviders": listOauthProviders.Result,
				"ldapProviders":  listLdapProviders.Result,
			},
		})
	}
}

// GetOAuthConfig returns OAuth config based on given provider
func GetOAuthConfig() web.HandlerFunc {
	return func(c *web.Context) error {
		getConfig := &query.GetCustomOAuthConfigByProvider{
			Provider: c.Param("provider"),
		}
		if err := bus.Dispatch(c, getConfig); err != nil {
			return c.Failure(err)
		}

		return c.Ok(getConfig.Result)
	}
}

// SaveOAuthConfig is used to create/edit OAuth configurations
func SaveOAuthConfig() web.HandlerFunc {
	return func(c *web.Context) error {
		action := actions.NewCreateEditOAuthConfig()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c,
			&cmd.UploadImage{
				Image:  action.Logo,
				Folder: "logos",
			},
			&cmd.SaveCustomOAuthConfig{
				ID:                action.ID,
				Logo:              action.Logo,
				Provider:          action.Provider,
				Status:            action.Status,
				DisplayName:       action.DisplayName,
				ClientID:          action.ClientID,
				ClientSecret:      action.ClientSecret,
				AuthorizeURL:      action.AuthorizeURL,
				TokenURL:          action.TokenURL,
				Scope:             action.Scope,
				ProfileURL:        action.ProfileURL,
				JSONUserIDPath:    action.JSONUserIDPath,
				JSONUserNamePath:  action.JSONUserNamePath,
				JSONUserEmailPath: action.JSONUserEmailPath,
			},
		); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// GetOAuthConfig returns OAuth config based on given provider
func GetLdapConfig() web.HandlerFunc {
	return func(c *web.Context) error {
		getConfig := &query.GetCustomLdapConfigByProvider{
			Provider: c.Param("provider"),
		}
		if err := bus.Dispatch(c, getConfig); err != nil {
			return c.Failure(err)
		}

		return c.Ok(getConfig.Result)
	}
}

// SaveLdapConfig is used to create/edit Ldap configurations
func SaveLdapConfig() web.HandlerFunc {
	return func(c *web.Context) error {
		action := actions.NewCreateEditLdapConfig()
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}
		if err := bus.Dispatch(c,
			&cmd.SaveCustomLdapConfig{
				ID:                    action.ID,
				Provider:              action.Provider,
				DisplayName:           action.DisplayName,
				Status:                action.Status,
				Protocol:              action.Protocol,
				CertCheck:             action.CertCheck,
				LdapHostname:          action.LdapHostname,
				LdapPort:              action.LdapPort,
				BindUsername:          action.BindUsername,
				BindPassword:          action.BindPassword,
				RootDN:                action.RootDN,
				Scope:                 action.Scope,
				UserSearchFilter:      action.UserSearchFilter,
				UsernameLdapAttribute: action.UsernameLdapAttribute,
				NameLdapAttribute:     action.NameLdapAttribute,
				MailLdapAttribute:     action.MailLdapAttribute,
			},
		); err != nil {
			return c.Failure(err)
		}
		return c.Ok(web.Map{})
	}
}

// TestLdapServer is used to test an LDAP provider
func TestLdapServer() web.HandlerFunc {
	return func(c *web.Context) error {
		testLdapServer := &cmd.TestLdapServer{Provider: c.Param("provider")}
		if err := bus.Dispatch(c, testLdapServer); err != nil {
			return c.GatewayTimeout()
		}

		return c.Ok(web.Map{})
	}
}
