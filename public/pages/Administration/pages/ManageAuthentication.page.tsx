import React from "react"

import { Button, OAuthProviderLogo, Icon, Field, Toggle, Form } from "@fider/components"
import { OAuthConfig, OAuthProviderOption, LdapConfig, LdapProviderOption } from "@fider/models"
import { OAuthForm } from "../components/OAuthForm"
import { LdapForm } from "../components/LdapForm"
import { actions, notify, Fider, Failure } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"

import IconPlay from "@fider/assets/images/heroicons-play.svg"
import IconPencilAlt from "@fider/assets/images/heroicons-pencil-alt.svg"

import { HStack, VStack } from "@fider/components/layout"

interface ManageAuthenticationPageProps {
  oauthProviders: OAuthProviderOption[];
  ldapProviders: LdapProviderOption[];
}

interface ManageAuthenticationPageState {
  isAddingOauth: boolean
  isAddingLdap: boolean
  editingOauth?: OAuthConfig
  editingLdap?: LdapConfig
  isEmailAuthAllowed: boolean
  canDisableEmailAuth: boolean
  error?: Failure
}

export default class ManageAuthenticationPage extends AdminBasePage<ManageAuthenticationPageProps, ManageAuthenticationPageState> {
  public id = "p-admin-authentication"
  public name = "authentication"
  public title = "Authentication"
  public subtitle = "Manage your site authentication"

  constructor(props: ManageAuthenticationPageProps) {
    super(props)
    this.state = {
      isAddingOauth: false,
      isAddingLdap: false,
      isEmailAuthAllowed: Fider.session.tenant.isEmailAuthAllowed,
      // TODO: Validate this
      canDisableEmailAuth: props.oauthProviders.map((o) => o.isEnabled).reduce((a, b) => a || b, false) || props.ldapProviders.map((o) => o.isEnabled).reduce((a, b) => a || b, false),
    }
  }

  /* OAUTH */

  private addNewOauth = async () => {
    this.setState({ isAddingOauth: true, editingOauth: undefined })
  }

  private editOauth = async (provider: string) => {
    const result = await actions.getOAuthConfig(provider)
    if (result.ok) {
      this.setState({ editingOauth: result.data, isAddingOauth: false })
    } else {
      notify.error("Failed to retrieve OAuth configuration. Try again later")
    }
  }

  private startTestOauth = async (provider: string) => {
    const redirect = `${Fider.settings.baseURL}/oauth/${provider}/echo`
    window.open(`/oauth/${provider}?redirect=${redirect}`, "oauth-test", "width=1100,height=600,status=no,menubar=no")
  }

  private cancelOauth = async () => {
    this.setState({ isAddingOauth: false, editingOauth: undefined })
  }

  /* LDAP */

  private addNewLdap = async () => {
    this.setState({ isAddingLdap: true, editingLdap: undefined })
  }

  private editLdap = async (provider: string) => {
    const result = await actions.getLdapConfig(provider)
    if (result.ok) {
      this.setState({ editingLdap: result.data, isAddingLdap: false })
    } else {
      notify.error("Failed to retrieve Ldap configuration. Try again later")
    }
  }

  private startTestLdap = async (provider: string, displayName: string) => {
    const result = await actions.testLdapServer(provider);
    if (result.ok)
    {
      notify.success("Success ! "+displayName+" is a valid LDAP server.")
    } else {
      notify.error(displayName+" is not available. Please review LDAP configuration")
    }
  };

  private cancelLdap = async () => {
    this.setState({ isAddingLdap: false, editingLdap: undefined })
  }

  private toggleEmailAuth = async (active: boolean) => {
    this.setState(
      () => ({
        isEmailAuthAllowed: active,
      }),
      async () => {
        const response = await actions.updateTenantEmailAuthAllowed(this.state.isEmailAuthAllowed)
        if (response.ok) {
          notify.success(`You successfully ${this.state.isEmailAuthAllowed ? "allowed" : "disallowed"} email authentication.`)
        } else {
          this.setState(
            () => ({
              isEmailAuthAllowed: !active,
              error: response.error,
            }),
            () => notify.error("Unable to save this setting.")
          )
        }
      }
    )
  }


  public content() {
    let enabledProvidersCount = 0
    for (const o of this.props.ldapProviders) {
      if (o.isEnabled) {
        enabledProvidersCount++
      }
    }
    for (const o of this.props.oauthProviders) {
      if (o.isEnabled) {
        enabledProvidersCount++
      }
    }
    const cantDisable = !this.state.isEmailAuthAllowed && enabledProvidersCount == 1

    if (this.state.isAddingOauth) {
      return <OAuthForm cantDisable={cantDisable} onCancel={this.cancelOauth} />
    }

    if (this.state.editingOauth) {
      return <OAuthForm cantDisable={cantDisable} config={this.state.editingOauth} onCancel={this.cancelOauth} />
    }

    if (this.state.isAddingLdap) {
      return <LdapForm cantDisable={cantDisable} onCancel={this.cancelLdap} />
    }

    if (this.state.editingLdap) {
      return <LdapForm cantDisable={cantDisable} config={this.state.editingLdap} onCancel={this.cancelLdap} />
    }

    const enabled = <span className="text-green-700">Enabled</span>
    const disabled = <span className="text-red-700">Disabled</span>

    return (
      <VStack spacing={8}>
        <div>
          <h2 className="text-display">General Authentication</h2>
          <Form error={this.state.error}>
            <Field label="Allow Email Authentication">
              <Toggle
                field="isEmailAuthAllowed"
                label={this.state.isEmailAuthAllowed ? "Allowed" : "Disallowed"}
                disabled={!Fider.session.user.isAdministrator || !this.state.canDisableEmailAuth}
                active={this.state.isEmailAuthAllowed}
                onToggle={this.toggleEmailAuth}
              />
              {!this.state.canDisableEmailAuth && (
                <p className="text-muted my-1">You need to configure another authentication provider before disabling email authentication.</p>
              )}
              <p className="text-muted my-1">
                When email-based authentication is disabled, users will not be allowed to sign in using their email. Thus, they will be forced to use another
                authentication provider, such as your preferred OAuth provider.{" "}
                <strong>Be sure to enable and test one before you turn this setting off!</strong>
              </p>
              <p className="text-muted mt-1">Note: Administrator accounts will still be allowed to sign in using their email.</p>
            </Field>
          </Form>
        </div>

        <div>
        <h2 className="text-display">LDAP Providers</h2>
        <p>
          You can use these section to add any authentication provider thats supports the LDAP protocol.
        </p>
        <VStack spacing={6}>
        {this.props.ldapProviders.map((o) => (
          <div key={o.provider}>
          <HStack justify="between">
            <HStack className="h-6">
              <strong>{o.displayName}</strong>
            </HStack>
            {(
              <HStack>
                {Fider.session.user.isAdministrator && (
                  <Button onClick={this.editLdap.bind(this, o.provider)} size="small">
                    <Icon sprite={IconPencilAlt} />
                    <span>Edit</span>
                  </Button>
                )}
                <Button onClick={this.startTestLdap.bind(this, o.provider, o.displayName)} size="small">
                  <Icon sprite={IconPlay} />
                  <span>Test</span>
                </Button>
              </HStack>
            )}
          </HStack>
          <div className="text-xs block my-1">{o.isEnabled ? enabled : disabled}</div>
        </div>
        ))}
        <div>
          {Fider.session.user.isAdministrator && (
            <Button variant="secondary" onClick={this.addNewLdap}>
                Add new LDAP server
            </Button>
          )}
        </div>
        </VStack>
        </div>

        <div>
          <h2 className="text-display">OAuth Providers</h2>
          <p>
            You can use these section to add any authentication provider thats supports the OAuth2 protocol. Additional information is available in our{" "}
            <a rel="noopener" className="text-link" target="_blank" href="https://getfider.com/docs/configuring-oauth/">
              OAuth Documentation
            </a>
            .
          </p>
          <VStack spacing={6}>
            {this.props.oauthProviders.map((o) => (
              <div key={o.provider}>
                <HStack justify="between">
                  <HStack className="h-6">
                    <OAuthProviderLogo option={o} />
                    <strong>{o.displayName}</strong>
                  </HStack>
                  {o.isCustomProvider && (
                    <HStack>
                      {Fider.session.user.isAdministrator && (
                        <Button onClick={this.editOauth.bind(this, o.provider)} size="small">
                          <Icon sprite={IconPencilAlt} />
                          <span>Edit</span>
                        </Button>
                      )}
                      <Button onClick={this.startTestOauth.bind(this, o.provider)} size="small">
                        <Icon sprite={IconPlay} />
                        <span>Test</span>
                      </Button>
                    </HStack>
                  )}
                </HStack>
                <div className="text-xs block my-1">{o.isEnabled ? enabled : disabled}</div>
                {o.isCustomProvider && (
                  <span className="text-muted">
                    <strong>Client ID:</strong> {o.clientID} <br />
                    <strong>Callback URL:</strong> {o.callbackURL}
                  </span>
                )}
              </div>
            ))}
            <div>
              {Fider.session.user.isAdministrator && (
                <Button variant="secondary" onClick={this.addNewOauth}>
                  Add new Oauth server
                </Button>
              )}
            </div>
          </VStack>
        </div>
      </VStack>
    )
  }
}
