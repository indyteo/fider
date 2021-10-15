import "./SignInControl.scss"

import React, { useState } from "react"
import { SocialSignInButton, Form, Button, Input, Message, Select, SelectOption} from "@fider/components"
import { Divider, HStack } from "@fider/components/layout"
import { device, actions, Failure, isCookieEnabled } from "@fider/services"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/macro"

interface SignInControlProps {
  useEmail: boolean
  redirectTo?: string
  onEmailSent?: (email: string) => void
}

export const SignInControl: React.FunctionComponent<SignInControlProps> = (props) => {
  const fider = useFider()
  const oauthProvidersLen = fider.settings.oauth.length;
  const ldapProvidersLen = fider.settings.ldap.length;
  const [showEmailForm, setShowEmailForm] = useState(fider.session.tenant ? fider.session.tenant.isEmailAuthAllowed : true)
  const [email, setEmail] = useState("")
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const [error, setError] = useState<Failure | undefined>(undefined)
  const [ldapError, setLdapError] = useState<Failure | undefined>(undefined);
  const [ldapProvider, _setLdapProvider] = useState((ldapProvidersLen > 0 && fider.settings.ldap[0].provider) || "");

  const forceShowEmailForm = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault()
    setShowEmailForm(true)
  }

  const signIn = async () => {
    const result = await actions.signIn(email)
    if (result.ok) {
      setEmail("")
      setError(undefined)
      if (props.onEmailSent) {
        props.onEmailSent(email)
      }
    } else if (result.error) {
      setError(result.error)
    }
  }

  const ldapSignIn = async () => {
    const result = await actions.ldapSignIn(username, password, ldapProvider);
    if (result.ok) {
      location.reload();
    } else if (result.error) {
      setLdapError(result.error);
    }
  };

  const setLdapProvider = (opt?: SelectOption) => {
    if(opt) {
      _setLdapProvider(opt.value);
    }
}

  if (!isCookieEnabled()) {
    return (
      <Message type="error">
        <h3 className="text-display">Cookies Required</h3>
        <p>Cookies are not enabled on your browser. Please enable cookies in your browser preferences to continue.</p>
      </Message>
    )
  }

  return (
    <div className="c-signin-control">

      {ldapProvidersLen > 0 && (
        <>
          <div className="l-signin-ldap">
          <p className="text-muted">
            <Trans id="signin.message.ldap">Connect using your LDAP account</Trans>
          </p>
          <Form error={ldapError}>
            <Input
              field="ldapUsername"
              value={username}
              placeholder="username"
              onChange={setUsername}
            />
            <Input
              field="ldapPassword"
              value={password}
              placeholder="password"
              onChange={setPassword}
              password={true}
            />
            <HStack justify="full" center={true} >
              <span className="text-category">
                <Trans id="signin.select.ldap">LDAP server</Trans>
              </span>
              <Select 
                field="ldapprovider"
                defaultValue={ldapProvider}
                options={fider.settings.ldap.map(x => ({
                  value: x.provider,
                  label: x.displayName,
                }))}
                onChange={setLdapProvider}
              />
            </HStack>
            <Button variant="primary" disabled={username === "" || password === ""} onClick={ldapSignIn}>
              Sign In
            </Button>
          </Form>
          </div>
          {(props.useEmail || oauthProvidersLen > 0) && <Divider />}
        </>
      )}

      {oauthProvidersLen > 0 && (
        <>
          <div className="c-signin-control__oauth mb-2">
            {fider.settings.oauth.map((o) => (
              <React.Fragment key={o.provider}>
                <SocialSignInButton option={o} redirectTo={props.redirectTo} />
              </React.Fragment>
            ))}
          </div>
          {props.useEmail && <Divider />}
        </>
      )}

      {props.useEmail &&
        (showEmailForm ? (
          <div>
            <p>
              <Trans id="signin.message.email">Enter your email address to sign in</Trans>
            </p>
            <Form error={error}>
              <Input
                field="email"
                value={email}
                autoFocus={!device.isTouch()}
                onChange={setEmail}
                placeholder="yourname@example.com"
                suffix={
                  <Button type="submit" variant="primary" disabled={email === ""} onClick={signIn}>
                    <Trans id="action.signin">Sign in</Trans>
                  </Button>
                }
              />
            </Form>
            {!fider.session.tenant.isEmailAuthAllowed && (
              <p className="text-red-700 mt-1">
                <Trans id="signin.message.onlyadmins">Currently only allowed to sign in to an administrator account</Trans>
              </p>
            )}
          </div>
        ) : (
          <div>
            <p className="text-muted">
              <Trans id="signin.message.emaildisabled">
                Email authentication has been disabled by an administrator. If you have an administrator account and need to bypass this restriction, please{" "}
                <a href="#" className="text-bold" onClick={forceShowEmailForm}>
                  click here
                </a>
                .
              </Trans>
            </p>
          </div>
        ))}
    </div>
  )
}
