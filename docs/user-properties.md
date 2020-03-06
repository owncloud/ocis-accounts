---
title: "User Properties"
date: 2018-05-02T00:00:00+00:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis-accounts
geekdocEditPath: edit/master/docs
geekdocFilePath: user-properties.md
---

{{< toc >}}

The ocis accounts service can store user properties. It currently uses the [OpenID Connect Standard claims](https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims).
For the UI we want to use the graph API, which has it's own understanding of [user properties](https://docs.microsoft.com/de-de/graph/api/resources/user?view=graph-rest-1.0#properties).

## Comparison

| oidc standard claims  | graphapi                                    | ldap posixaccount | ldap inetorgperson   | oc10 |
|-----------------------|---------------------------------------------|-------------------|----------------------|------|
| sub                   | - / maybe id                                |                   |                      |      |
| name                  | displayName                                 |                   | displayName          |      |
| given_name            | givenName                                   |                   | givenName            |      |
| family_name           | surname                                     |                   |                      |      |
| middle_name           | -                                           |                   |                      |      |
| nickname              | -                                           |                   |                      |      |
| preferred_username    | preferredName                               |                   |                      |      |
| profile               | -                                           |                   | Profile              |      |
| picture               | photo                                       |                   | jpegPhoto,ldapPhoto? |      |
| website               | mySite                                      |                   |                      |      |
| email                 | mail                                        |                   | mail                 |      |
| email_verified        | -                                           |                   |                      |      |
| gender                | -                                           |                   |                      |      |
| birthdate             | birthday                                    |                   |                      |      |
| zoneinfo              | -                                           |                   | Timezone             |      |
| locale                | preferredLanguage                           |                   | preferredLanguage    |      |
| phone_number          | businessPhones,mobilePhone                  |                   | homePhone,mobile     |      |
| phone_number_verified | -                                           |                   |                      |      |
| address               | country,state,postalCode,city,streetAddress |                   |                      |      |
| updated_at            | -                                           |                   |                      |      |
| -                     | aboutMe                                     | description       |                      |      |
| -                     | accountEnabled                              |                   | loginDisabled        |      |
| -                     | ageGroup                                    |                   |                      |      |
| -                     | assignedLicenses                            |                   |                      |      |
| -                     | assignedPlans                               |                   |                      |      |
| -                     | companyName                                 |                   |                      |      |
| -                     | consentProvidedForMinor                     |                   |                      |      |
| -                     | createdDateTime                             |                   |                      |      |
| -                     | creationType                                |                   |                      |      |
| -                     | deletedDateTime                             |                   |                      |      |
| -                     | department                                  |                   |                      |      |
| -                     | employeeId                                  |                   | employeeNumber       |      |
| -                     | faxNumber                                   |                   |                      |      |
| -                     | hireDate                                    |                   |                      |      |
| -                     | id                                          |                   |                      |      |
| -                     | identities                                  |                   |                      |      |
| -                     | imAddresses                                 |                   |                      |      |
| -                     | interests                                   |                   |                      |      |
| -                     | isResourceAccount                           |                   |                      |      |
| -                     | jobTitle                                    |                   |                      |      |
| -                     | lastPasswordChangeDateTime                  |                   |                      |      |
| -                     | legalAgeGroupClassification                 |                   |                      |      |
| -                     | licenseAssignmentStates                     |                   |                      |      |
| -                     | mailboxSettings                             |                   |                      |      |
| -                     | mailNickname                                |                   |                      |      |
| -                     | officeLocation                              |                   |                      |      |
| -                     | onPremisesDistinguishedName                 |                   |                      |      |
| -                     | onPremisesDomainName                        |                   |                      |      |
| -                     | onPremisesExtensionAttributes               |                   |                      |      |
| -                     | onPremisesImmutableId                       |                   |                      |      |
| -                     | onPremisesLastSyncDateTime                  |                   |                      |      |
| -                     | onPremisesProvisioningErrors                |                   |                      |      |
| -                     | onPremisesSamAccountName                    |                   |                      |      |
| -                     | onPremisesSecurityIdentifier                |                   |                      |      |
| -                     | onPremisesSyncEnabled                       |                   |                      |      |
| -                     | onPremisesUserPrincipalName                 |                   |                      |      |
| -                     | otherMails                                  |                   |                      |      |
| -                     | passwordPolicies                            |                   |                      |      |
| -                     | passwordProfile                             |                   |                      |      |
| -                     | pastProjects                                |                   |                      |      |
| -                     | preferredDataLocation                       |                   |                      |      |
| -                     | provisionedPlans                            |                   |                      |      |
| -                     | proxyAddresses                              |                   |                      |      |
| -                     | refreshTokensValidFromDateTime              |                   |                      |      |
| -                     | responsibilities                            |                   |                      |      |
| -                     | schools                                     |                   |                      |      |
| -                     | showInAddressList                           |                   |                      |      |
| -                     | skills                                      |                   |                      |      |
| -                     | signInSessionsValidFromDateTime             |                   |                      |      |
| -                     | usageLocation                               |                   |                      |      |
| -                     | userPrincipalName                           |                   |                      |      |
| -                     | userType                                    |                   |                      |      |
| -                     |                                             | uid               | uid                  |      |
| -                     |                                             | cn                |                      |      |
| -                     |                                             | uidNumber         |                      |      |
| -                     |                                             | gidNumber         |                      |      |
| -                     |                                             | unixHomeDirectory |                      |      |
| -                     |                                             | homeDirectory     |                      |      |
| -                     |                                             | userPassword      |                      |      |
| -                     |                                             | unixUserPassword  |                      |      |
| -                     |                                             | loginShell        |                      |      |
| -                     |                                             | gecos             |                      |      |
| -                     |                                             |                   |                      |      |
| -                     |                                             |                   |                      |      |
| -                     |                                             |                   |                      |      |
| -                     |                                             |                   |                      |      |
| -                     |                                             |                   |                      |      |
| -                     |                                             |                   |                      |      |
| -                     |                                             |                   |                      |      |
| -                     |                                             |                   |                      |      |



## Extensions in the graph API
The graph api can be used to store arbitrary kv pairs for many resources, eg. a user.
See https://docs.microsoft.com/de-de/graph/api/resources/opentypeextension?view=graph-rest-1.0 for details.

## Notas
- Table formatted with http://markdowntable.com/
- The ldapwiki has some explcicit mappings for some ldap attributes to oidc claims, eg. [Profile](https://ldapwiki.com/wiki/Profile#section-Profile-OpenIDConnectScopes). See https://ldapwiki.com/wiki/OpenID%20Connect%20Scopes for more details

## TODO
- add more inetorgperson attributes https://ldapwiki.com/wiki/InetOrgPerson
- add OC10 account properties and settings
- distinguish from user settings, eg. default sort order and column ... per directory? or is that a setting for a folder?
- check what konnectd can read and map from ldap
- check the [ldapwiki oidc connect scopes page](https://ldapwiki.com/wiki/OpenID%20Connect%20Scopes) for more *well known* mappings