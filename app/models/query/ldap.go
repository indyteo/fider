package query

import (
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
)

type GetCustomLdapConfigByProvider struct {
	Provider string

	Result *entity.LdapConfig
}

type ListCustomLdapConfig struct {
	Result []*entity.LdapConfig
}

type GetLdapProfile struct {
	Provider string
	Username string

	Result *dto.LdapUserProfile
}

type ListActiveLdapProviders struct {
	Result []*dto.LdapProviderOption
}

type ListAllLdapProviders struct {
	Result []*dto.LdapProviderOption
}
