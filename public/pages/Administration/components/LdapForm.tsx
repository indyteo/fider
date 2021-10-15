import React, { useState } from "react";
import { LdapConfig, LdapConfigStatus, LdapProtocols, LdapScopeStatus } from "@fider/models";
import { useFider } from "@fider/hooks";
import { Failure, actions } from "@fider/services";
import { Form, Button, Input, Toggle, Field, Select, SelectOption } from "@fider/components";
import { HStack } from "@fider/components/layout";

interface LdapFormProps {
    config?: LdapConfig
    onCancel: () => void
    cantDisable: boolean
  }

export const LdapForm: React.FC<LdapFormProps> = props => {

    const fider = useFider();
    const [provider] = useState((props.config && props.config.provider) || "");
    const [displayName, setDisplayName] = useState((props.config && props.config.displayName) || "");
    const [enabled, setEnabled] = useState((props.config && props.config.status === LdapConfigStatus.Enabled) || false);
    const [protocol,_setProtocol] = useState((props.config && props.config.protocol) || 1);
    const [certCheck, setCertCheck] = useState(!props.config || props.config.certCheck);
    const [scope, _setScope] = useState((props.config && props.config.scope) || 3);
    const [ldapHostname, setLdapHostname] = useState((props.config && props.config.ldapHostname) || "");
    const [ldapPort,setLdapPort] = useState((props.config && props.config.ldapPort) || "389");
    const [bindUsername, setBindUsername] = useState((props.config && props.config.bindUsername) || "");
    const [bindPassword, setBindPassword] = useState((props.config && props.config.bindPassword) || "");
    const [bindPasswordEnabled,setBindPasswordEnabled] = useState(!props.config);
    const [rootDN, setRootDN] = useState((props.config && props.config.rootDN) || "");
    const [userSearchFilter, setUserSearchFilter] = useState((props.config && props.config.userSearchFilter) || "");
    const [usernameLdapAttribute, setUsernameLdapAttribute] = useState((props.config && props.config.usernameLdapAttribute) || "");
    const [nameLdapAttribute, setNameLdapAttribute] = useState((props.config && props.config.nameLdapAttribute) || "");
    const [mailLdapAttribute, setMailLdapAttribute] = useState((props.config && props.config.mailLdapAttribute) || "");
    const [error, setError] = useState<Failure | undefined>();

    const handleSave = async () => {
        const result = await actions.saveLdapConfig({
          provider,
          status: enabled ? LdapConfigStatus.Enabled : LdapConfigStatus.Disabled,
          protocol,
          certCheck,
          displayName,
          ldapHostname,
          ldapPort,
          bindUsername,
          bindPassword: bindPasswordEnabled ? bindPassword : "",
          rootDN,
          scope,
          userSearchFilter,
          usernameLdapAttribute,
          nameLdapAttribute,
          mailLdapAttribute,
        });
        if (result.ok) {
          location.reload();
        } else {
          setError(result.error);
        }
      };

    const handleCancel = async () => {
        props.onCancel();
    };

    const enableBindPassword = () => {
        setBindPassword("");
        setBindPasswordEnabled(true);
    };

    const setProtocol = (opt?: SelectOption) => {
        if(opt) {
            let proto = parseInt(opt.value);
            _setProtocol(proto);
            if(proto === 3 /* LDAPS */)
            {
                setLdapPort("636")
            }
            else
            {
                setLdapPort("389")
            }
        }
    }

    const setScope = (opt?: SelectOption) => {
        if(opt) {
            _setScope(parseInt(opt.value));
        }
    }

    const title = props.config ? `LDAP Provider: ${props.config.displayName}` : "New LDAP Provider";
    return (
        <>
            <h2 className="text-title mb-2">{title}</h2>
            <Form error={error}>

                <Input
                    field="displayName"
                    label="Display Name"
                    maxLength={50}
                    value={displayName}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setDisplayName}
                    placeholder="My LDAP server"
                >
                <p className="text-muted">The name that will be displayed in the login form</p>
                </Input>

                <Select 
                    field="protocol"
                    label="LDAP protocol"
                    defaultValue={protocol.toString()}
                    options={Object.entries(LdapProtocols).map(([k, v]) => ({
                        value: v.toString(),
                        label: k,
                      }))}
                    onChange={setProtocol}
                />

                {(protocol === 2 /* LDAP + TLS */ || protocol === 3 /* LDAPS */) && <Field label="Check certificate">
                <Toggle active={certCheck} onToggle={setCertCheck} />
                    {certCheck ? <>
                        <span>Enabled</span>
                        <p className="info">Certificate <strong>will</strong> be checked.</p>
                        </> : <>
                        <span>Disabled</span>
                        <p className="info">Certificate <strong>won't</strong> be checked.</p>
                        </>}
                </Field>}

                <Input
                    field="ldapHostname"
                    label="LDAP Hostname"
                    maxLength={300}
                    value={ldapHostname}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setLdapHostname}
                    placeholder="localhost"
                />

                <Input
                    field="ldapPort"
                    label="LDAP Port"
                    maxLength={10}
                    value={ldapPort}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setLdapPort}
                />

                <Input
                    field="bindUsername"
                    label="Bind Username"
                    maxLength={100}
                    value={bindUsername}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setBindUsername}
                    placeholder="cn=readonly,dc=example,dc=org"
                />

                <Input
                    field="bindPassword"
                    label="Bind Password"
                    password={true}
                    maxLength={100}
                    value={bindPassword}
                    disabled={!bindPasswordEnabled}
                    onChange={setBindPassword}
                    afterLabel={
                        !bindPasswordEnabled ? (
                        <>
                            <span className="info">omitted for security reasons</span>
                            <span className="info clickable" onClick={enableBindPassword}>
                            change
                            </span>
                        </>
                        ) : undefined
                    }
                    placeholder="readonly_password"
                />

                <Input
                    field="rootDN"
                    label="Root DN"
                    maxLength={300}
                    value={rootDN}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setRootDN}
                    placeholder="dc=example,dc=org"
                />

                <Select 
                    field="scope"
                    label="Search scope"
                    defaultValue={scope.toString()}
                    options={Object.entries(LdapScopeStatus).map(([k, v]) => ({
                        value: v.toString(),
                        label: k,
                      }))}
                    onChange={setScope}
                />  
                
                <Input
                    field="userSearchFilter"
                    label="User Search Filter"
                    maxLength={500}
                    value={userSearchFilter}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setUserSearchFilter}
                    placeholder="(objectClass=inetOrgPerson)"
                />

                <Input
                    field="usernameLdapAttribute"
                    label="Username Ldap Attribute"
                    maxLength={100}
                    value={usernameLdapAttribute}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setUsernameLdapAttribute}
                    placeholder="uid"
                />

                <Input
                    field="nameLdapAttribute"
                    label="Full Name Ldap Attribute"
                    maxLength={100}
                    value={nameLdapAttribute}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setNameLdapAttribute}
                    placeholder="displayName"
                />

                <Input
                    field="mailLdapAttribute"
                    label="Mail Ldap Attribute"
                    maxLength={100}
                    value={mailLdapAttribute}
                    disabled={!fider.session.user.isAdministrator}
                    onChange={setMailLdapAttribute}
                    placeholder="mail"
                />

                <Field label="Status">
                    <Toggle field="status" active={enabled} onToggle={setEnabled} />
                    <div className="mt-1">
                    {enabled ? <>
                        <span>Enabled</span>
                        <p className="info">
                        This provider will be available for everyone to use during the sign in process. It is recommended that
                        you keep it disable and test it before enabling it. The Test button is available after saving this
                        configuration.
                        </p>
                    </> : <>
                        <span>Disabled</span>
                        <p className="info">Users won't be able to sign in with this Provider.</p>
                    </>}
                    </div>
                </Field>

                <HStack className="mt-2">
                <Button variant="primary" onClick={handleSave}>
                    Save {title}
                </Button>
                <Button variant="tertiary" onClick={handleCancel}>
                    Cancel
                </Button>
                </HStack>

            </Form>
        </>
        );
}