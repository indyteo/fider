package ldap

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"time"

	// We use v3 version of the go-ldap package
	// Allows defining timeout and ssl behaviour when connecting to ldap
	ldap "github.com/go-ldap/ldap/v3"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "LDAP"
}

func (s Service) Category() string {
	return "LDAP"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	bus.AddHandler(getLdapProfile)
	bus.AddHandler(verifyLdapUser)
	bus.AddHandler(listActiveLdapProviders)
	bus.AddHandler(listAllLdapProviders)
	bus.AddHandler(testLdapServer)
}

// Package level default timeout
var defaultTimeout = 5 * time.Second

// getConfig returns the properties of a given LDAP provider
func getConfig(ctx context.Context, provider string) (*entity.LdapConfig, error) {

	getCustomLdap := &query.GetCustomLdapConfigByProvider{Provider: provider}
	err := bus.Dispatch(ctx, getCustomLdap)
	if err != nil {
		log.Warnf(ctx, " Could not get LDAP information for @{Provider}", dto.Props{"Provider": provider})
		return nil, err
	}

	return getCustomLdap.Result, nil
}

// getLdapConn returns a valid connection to the LDAP server
func getLdapConn(ctx context.Context, provider string) (*ldap.Conn, *entity.LdapConfig, string, error) {

	// Get LDAP provider configuration from database
	ldapConfig, err := getConfig(ctx, provider)
	if err != nil {
		return nil, ldapConfig, "", err
	}

	// Get protocol from LDAP provider configuration
	protocol := "ldap://"
	if ldapConfig.Protocol == enum.LDAPS {
		protocol = "ldaps://"
	}

	ldapURL := protocol + ldapConfig.LdapHostname + ":" + ldapConfig.LdapPort

	// Connect to LDAP with short timeout
	// https://github.com/go-ldap/ldap/issues/310
	tlsConfig := &tls.Config{InsecureSkipVerify: !ldapConfig.CertCheck}
	l, err := ldap.DialURL(ldapURL, ldap.DialWithDialer(&net.Dialer{Timeout: defaultTimeout}), ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		log.Errorf(ctx, "Could not dial LDAP url : @{LdapURL}", dto.Props{"LdapURL": ldapURL})
		return nil, ldapConfig, ldapURL, err
	}

	// Reconnect with ldap+TLS if necessary
	if ldapConfig.Protocol == enum.LDAPTLS {
		err = l.StartTLS(tlsConfig)
		if err != nil {
			log.Errorf(ctx, "Could not activate TLS for @{LdapURL}", dto.Props{"LdapURL": ldapURL})
			return nil, ldapConfig, ldapURL, err
		}
	}

	return l, ldapConfig, ldapURL, nil
}

// testLdapServer test if LDAP server can be accessed by the read only user
func testLdapServer(ctx context.Context, c *cmd.TestLdapServer) error {

	l, ldapConfig, ldapURL, err := getLdapConn(ctx, c.Provider)
	if err != nil {
		return err
	}

	// Bind with read only user
	err = l.Bind(ldapConfig.BindUsername, ldapConfig.BindPassword)
	if err != nil {
		log.Errorf(ctx, "Could not bind with @{Username} for @{LdapURL}", dto.Props{"Username": ldapConfig.BindUsername, "LdapURL": ldapURL})
		return err
	}

	// closing the connection at the end of function
	defer l.Close()

	return nil
}

// verifyLdapUser checks that the user exists and its username/password are valid
func verifyLdapUser(ctx context.Context, c *cmd.VerifyLdapUser) error {

	l, ldapConfig, ldapURL, err := getLdapConn(ctx, c.Provider)
	if err != nil {
		return err
	}

	// First Bind with read only user
	err = l.Bind(ldapConfig.BindUsername, ldapConfig.BindPassword)
	if err != nil {
		log.Errorf(ctx, "Could not bind with @{Username} for @{LdapURL}", dto.Props{"Username": ldapConfig.BindUsername, "LdapURL": ldapURL})
		return err
	}

	// Search for given username
	var filter = "(&" + ldapConfig.UserSearchFilter + "(" + ldapConfig.UsernameLdapAttribute + "=" + c.Username + "))"
	searchRequest := ldap.NewSearchRequest(
		ldapConfig.RootDN,
		ldapConfig.Scope-1, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Errorf(ctx, "Could not search LDAP with @{Filter}", dto.Props{"Filter": filter})
		return err
	}

	// Verify search results
	// TODO : Deal with cases where len()>=2
	if len(sr.Entries) != 1 {
		log.Debugf(ctx, "@{Length} user found with filter (&@{Filter}", dto.Props{"Length": len(sr.Entries), "Filter": filter})
		return errors.New("User not found")
	}

	// Get DN of the user to be tested
	userDN := sr.Entries[0].DN

	// Bind as user to verify their password
	err = l.Bind(userDN, c.Password)
	if err != nil {
		log.Debugf(ctx, "Could not bind with @{User}. Possible password error.", dto.Props{"User": userDN})
		return err
	}

	// closing the connection at the end of function
	defer l.Close()

	return nil
}

// getLdapProfile gets user information in LDAP
func getLdapProfile(ctx context.Context, q *query.GetLdapProfile) error {

	l, ldapConfig, ldapURL, err := getLdapConn(ctx, q.Provider)
	if err != nil {
		return err
	}

	// Bind with read only user
	err = l.Bind(ldapConfig.BindUsername, ldapConfig.BindPassword)
	if err != nil {
		log.Errorf(ctx, "Could not bind with @{Username} for @{LdapURL}", dto.Props{"Username": ldapConfig.BindUsername, "LdapURL": ldapURL})
		return err
	}

	// Search for user id, name and email
	var filter = "(&" + ldapConfig.UserSearchFilter + "(" + ldapConfig.UsernameLdapAttribute + "=" + q.Username + "))"
	searchRequest := ldap.NewSearchRequest(
		ldapConfig.RootDN,
		ldapConfig.Scope-1, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{ldapConfig.MailLdapAttribute, ldapConfig.NameLdapAttribute, ldapConfig.UsernameLdapAttribute},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Errorf(ctx, "Could not search ldap with @{Filter}", dto.Props{"Filter": filter})
		return err
	}

	// Verify search results
	if len(sr.Entries) != 1 {
		log.Errorf(ctx, "@{Length} user found with @{Filter}", dto.Props{"Length": sr.Entries, "Filter": filter})
		return errors.New("User not found")
	}

	// Create user profile
	profile := &dto.LdapUserProfile{
		ID:    strings.TrimSpace(sr.Entries[0].GetAttributeValue(ldapConfig.UsernameLdapAttribute)),
		Name:  strings.TrimSpace(sr.Entries[0].GetAttributeValue(ldapConfig.NameLdapAttribute)),
		Email: strings.ToLower(sr.Entries[0].GetAttributeValue(ldapConfig.MailLdapAttribute)),
	}

	defer l.Close()

	q.Result = profile

	return nil
}

// listActiveLdapProviders returns a list of enabled LDAP providers */
func listActiveLdapProviders(ctx context.Context, q *query.ListActiveLdapProviders) error {

	allLdapProviders := &query.ListAllLdapProviders{}
	err := bus.Dispatch(ctx, allLdapProviders)
	if err != nil {
		return err
	}

	list := make([]*dto.LdapProviderOption, 0)

	for _, p := range allLdapProviders.Result {
		if p.IsEnabled {
			list = append(list, p)
		}
	}
	q.Result = list
	return nil
}

// listAllLdapProviders returns a list of all LDAP providers (disabled or enabled) */
func listAllLdapProviders(ctx context.Context, q *query.ListAllLdapProviders) error {

	ldapProviders := &query.ListCustomLdapConfig{}
	err := bus.Dispatch(ctx, ldapProviders)
	if err != nil {
		return errors.Wrap(err, "failed to get list of custom Ldap providers")
	}

	list := make([]*dto.LdapProviderOption, 0)

	for _, p := range ldapProviders.Result {
		list = append(list, &dto.LdapProviderOption{
			Provider:    p.Provider,
			DisplayName: p.DisplayName,
			IsEnabled:   p.Status == enum.LdapConfigEnabled,
		})
	}

	q.Result = list
	return nil
}
