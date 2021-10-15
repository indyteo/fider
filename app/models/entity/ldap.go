package entity

import "encoding/json"

// LdapConfig is the configuration of a custom LDAP provider
type LdapConfig struct {
	ID                    int
	Provider              string
	DisplayName           string
	Status                int
	Protocol              int
	CertCheck             bool
	LdapHostname          string
	LdapPort              string
	BindUsername          string
	BindPassword          string
	RootDN                string
	Scope                 int
	UserSearchFilter      string
	UsernameLdapAttribute string
	NameLdapAttribute     string
	MailLdapAttribute     string
}

// MarshalJSON returns the JSON encoding of LdapConfig
func (o LdapConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":                    o.ID,
		"provider":              o.Provider,
		"displayName":           o.DisplayName,
		"status":                o.Status,
		"protocol":              o.Protocol,
		"certCheck":             o.CertCheck,
		"ldapHostname":          o.LdapHostname,
		"ldapPort":              o.LdapPort,
		"bindUsername":          o.BindUsername,
		"bindPassword":          "password will remain secret",
		"rootDN":                o.RootDN,
		"scope":                 o.Scope,
		"userSearchFilter":      o.UserSearchFilter,
		"usernameLdapAttribute": o.UsernameLdapAttribute,
		"nameLdapAttribute":     o.NameLdapAttribute,
		"mailLdapAttribute":     o.MailLdapAttribute,
	})
}
