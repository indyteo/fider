export interface OAuthProviderOption {
  provider: string
  displayName: string
  clientID: string
  url: string
  callbackURL: string
  logoBlobKey: string
  isCustomProvider: boolean
  isEnabled: boolean
}

export interface SystemSettings {
  mode: string
  locale: string
  version: string
  environment: string
  domain: string
  hasLegal: boolean
  isBillingEnabled: boolean
  baseURL: string
  assetsURL: string
  oauth: OAuthProviderOption[]
  // Adding LDAP
  ldap: LdapProviderOption[]
}

export interface UserSettings {
  [key: string]: string
}

export const OAuthConfigStatus = {
  Disabled: 1,
  Enabled: 2,
}

export interface OAuthConfig {
  provider: string
  displayName: string
  status: number
  clientID: string
  clientSecret: string
  authorizeURL: string
  tokenURL: string
  profileURL: string
  logoBlobKey: string
  scope: string
  jsonUserIDPath: string
  jsonUserNamePath: string
  jsonUserEmailPath: string
}

export interface ImageUpload {
  bkey?: string
  upload?: {
    fileName?: string
    content?: string
    contentType?: string
  }
  remove: boolean
}

// LDAP provider options
export interface LdapProviderOption {
  provider: string
  displayName: string
  isEnabled: boolean
}

// Similarly to OauthConfigStatus
export const LdapConfigStatus = {
  Disabled: 1,
  Enabled: 2
}

// These values derive from go-ldap v3
// Values are +1 compared to go-ldap
export const LdapScopeStatus = {
  ScopeBaseObject : 1,
  ScopeSingleLevel: 2,
  ScopeWholeSubtree: 3
}

// These values derive from go-ldap v3
// Values are +1 compared to go-ldap
export const LdapProtocols = {
  "ldap://" : 1,
  "ldap:// + StartTLS": 2,
  "ldaps://": 3
}

// Full LDAP config
export interface LdapConfig {
  provider: string
  displayName: string
  status: number
  protocol: number
  certCheck: boolean
  ldapHostname: string
  ldapPort: string
  bindUsername: string
  bindPassword: string
  rootDN: string
  scope: number
  userSearchFilter: string
  usernameLdapAttribute: string
  nameLdapAttribute: string
  mailLdapAttribute: string
}
