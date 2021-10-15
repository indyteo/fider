package postgres

import (
	"context"
	"strconv"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbLdapConfig struct {
	ID                    int    `db:"id"`
	Provider              string `db:"provider"`
	DisplayName           string `db:"display_name"`
	Status                int    `db:"status"`
	Protocol              int    `db:"protocol"`
	CertCheck             bool   `db:"cert_check"`
	LdapHostname          string `db:"ldap_hostname"`
	Port                  int    `db:"ldap_port"`
	BindUsername          string `db:"bind_username"`
	BindPassword          string `db:"bind_password"`
	RootDN                string `db:"root_dn"`
	Scope                 int    `db:"scope"`
	UserSearchFilter      string `db:"user_search_filter"`
	UsernameLdapAttribute string `db:"username_ldap_attribute"`
	NameLdapAttribute     string `db:"name_ldap_attribute"`
	MailLdapAttribute     string `db:"mail_ldap_attribute"`
}

func (m *dbLdapConfig) toModel() *entity.LdapConfig {
	return &entity.LdapConfig{
		ID:                    m.ID,
		Provider:              m.Provider,
		DisplayName:           m.DisplayName,
		Status:                m.Status,
		Protocol:              m.Protocol,
		CertCheck:             m.CertCheck,
		LdapHostname:          m.LdapHostname,
		LdapPort:              strconv.Itoa(m.Port),
		BindUsername:          m.BindUsername,
		BindPassword:          m.BindPassword,
		RootDN:                m.RootDN,
		Scope:                 m.Scope,
		UserSearchFilter:      m.UserSearchFilter,
		UsernameLdapAttribute: m.UsernameLdapAttribute,
		NameLdapAttribute:     m.NameLdapAttribute,
		MailLdapAttribute:     m.MailLdapAttribute,
	}
}

func getCustomLdapConfigByProvider(ctx context.Context, q *query.GetCustomLdapConfigByProvider) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if tenant == nil {
			return app.ErrNotFound
		}

		config := &dbLdapConfig{}
		err := trx.Get(config, `
		SELECT id, provider, display_name, status,
					ldap_hostname, ldap_port, 
					bind_username, bind_password, root_dn,
					scope, user_search_filter, username_ldap_attribute, name_ldap_attribute,
		       		mail_ldap_attribute, protocol, cert_check
		FROM ldap_providers
		WHERE tenant_id = $1 AND provider = $2
		`, tenant.ID, q.Provider)
		if err != nil {
			return err
		}

		q.Result = config.toModel()
		return nil
	})
}

func listCustomLdapConfig(ctx context.Context, q *query.ListCustomLdapConfig) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		configs := []*dbLdapConfig{}

		if tenant != nil {
			err := trx.Select(&configs, `
			SELECT id, provider, display_name, status,
						 ldap_hostname, ldap_port, 
						 bind_username, bind_password, root_dn,
						 scope, user_search_filter, username_ldap_attribute, protocol, cert_check
			FROM ldap_providers
			WHERE tenant_id = $1
			ORDER BY id`, tenant.ID)
			if err != nil {
				return err
			}
		}

		q.Result = make([]*entity.LdapConfig, len(configs))
		for i, config := range configs {
			q.Result[i] = config.toModel()
		}
		return nil
	})
}

func saveCustomLdapConfig(ctx context.Context, c *cmd.SaveCustomLdapConfig) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var err error

		if c.ID == 0 {
			query := `INSERT INTO ldap_providers (
				tenant_id, provider, display_name, status,
				ldap_hostname, ldap_port, bind_username,
				bind_password, root_dn, scope, user_search_filter,
				username_ldap_attribute, name_ldap_attribute, mail_ldap_attribute, protocol, cert_check
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
			RETURNING id`

			err = trx.Get(&c.ID, query, tenant.ID, c.Provider,
				c.DisplayName, c.Status, c.LdapHostname, c.LdapPort,
				c.BindUsername, c.BindPassword, c.RootDN,
				c.Scope, c.UserSearchFilter, c.UsernameLdapAttribute,
				c.NameLdapAttribute, c.MailLdapAttribute, c.Protocol, c.CertCheck)

		} else {
			query := `
				UPDATE ldap_providers 
				SET display_name = $3, status = $4, ldap_hostname = $5, ldap_port = $6, 
				bind_username = $7, bind_password = $8, root_dn = $9, scope = $10, 
				user_search_filter = $11, username_ldap_attribute = $12, name_ldap_attribute = $13,
				mail_ldap_attribute = $14, protocol = $15, cert_check = $16
			WHERE tenant_id = $1 AND id = $2`

			_, err = trx.Execute(query, tenant.ID, c.ID,
				c.DisplayName, c.Status, c.LdapHostname, c.LdapPort,
				c.BindUsername, c.BindPassword, c.RootDN,
				c.Scope, c.UserSearchFilter, c.UsernameLdapAttribute,
				c.NameLdapAttribute, c.MailLdapAttribute, c.Protocol, c.CertCheck)
		}

		if err != nil {
			return errors.Wrap(err, "failed to save Ldap Provider")
		}

		return nil
	})
}
