create table if not exists ldap_providers (
  id                       serial not null, 
  tenant_id                int not null,
  provider                 varchar(30) not null,
  display_name             varchar(50) not null,
  status                   int not null,
  ldap_hostname            varchar(300) not null,
  ldap_port                int not null,
  cert_check               bool not null,
  bind_username            varchar(100) not null,
  bind_password            varchar(100) not null,
  root_dn                  varchar(250) not null,
  scope                    int not null,
  user_search_filter       varchar(500) not null,
  username_ldap_attribute  varchar(100) not null,
  name_ldap_attribute      varchar(100) not null,
  mail_ldap_attribute      varchar(100) not null,
  protocol                 int not null,
  created_on               timestamptz not null default now(),
  primary key (id),
  foreign key (tenant_id) references tenants(id)
);

CREATE UNIQUE INDEX tenant_id_ldap_provider_key ON ldap_providers (tenant_id, provider);
