package cmd

import (
	"github.com/getfider/fider/app/models/dto"
)

type TestLdapServer struct {
	Provider string
}

type VerifyLdapUser struct {
	Provider string
	Username string
	Password string
}

type SaveCustomLdapConfig struct {
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

type ParseLdapRawProfile struct {
	Provider string
	Body     string
	Result   *dto.LdapUserProfile
}
