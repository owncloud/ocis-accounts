package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-accounts/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"ACCOUNTS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"ACCOUNTS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"ACCOUNTS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			DefaultText: "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "namespace",
			Aliases:     []string{"ns"},
			Value:       "com.owncloud.api",
			DefaultText: "com.owncloud.api",
			Usage:       "namespace",
			EnvVars:     []string{"ACCOUNTS_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
		&cli.StringFlag{
			Name:        "address",
			Aliases:     []string{"addr"},
			Value:       "localhost:9180",
			DefaultText: "localhost:9180",
			Usage:       "service endpoint",
			EnvVars:     []string{"ACCOUNTS_ADDRESS"},
			Destination: &cfg.Server.Address,
		},
		// LDAP

		&cli.StringFlag{
			Name:        "ldap-hostname",
			Value:       "localhost",
			Usage:       "LDAP hostname",
			EnvVars:     []string{"ACCOUNTS_LDAP_HOSTNAME"},
			Destination: &cfg.LDAP.Hostname,
		},
		&cli.IntFlag{
			Name:        "ldap-port",
			Value:       9126,
			Usage:       "LDAP port",
			EnvVars:     []string{"ACCOUNTS_LDAP_PORT"},
			Destination: &cfg.LDAP.Port,
		},
		&cli.StringFlag{
			Name:        "ldap-base-dn",
			Value:       "dc=example,dc=org",
			Usage:       "LDAP basedn",
			EnvVars:     []string{"ACCOUNTS_LDAP_BASE_DN"},
			Destination: &cfg.LDAP.BaseDN,
		},
		&cli.StringFlag{
			Name:        "ldap-userfilter",
			Value:       "(&(objectclass=posixAccount)(cn=%s))",
			Usage:       "LDAP userfilter",
			EnvVars:     []string{"ACCOUNTS_LDAP_USERFILTER"},
			Destination: &cfg.LDAP.UserFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-groupfilter",
			Value:       "(&(objectclass=posixGroup)(cn=%s))",
			Usage:       "LDAP groupfilter",
			EnvVars:     []string{"ACCOUNTS_LDAP_GROUPFILTER"},
			Destination: &cfg.LDAP.GroupFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-dn",
			Value:       "cn=reva,ou=sysusers,dc=example,dc=org",
			Usage:       "LDAP bind dn",
			EnvVars:     []string{"ACCOUNTS_LDAP_BIND_DN"},
			Destination: &cfg.LDAP.BindDN,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-password",
			Value:       "reva",
			Usage:       "LDAP bind password",
			EnvVars:     []string{"ACCOUNTS_LDAP_BIND_PASSWORD"},
			Destination: &cfg.LDAP.BindPassword,
		},
		&cli.StringFlag{
			Name:        "ldap-idp",
			Value:       "https://localhost:9200",
			Usage:       "Identity provider to use for users",
			EnvVars:     []string{"ACCOUNTS_LDAP_IDP"},
			Destination: &cfg.LDAP.IDP,
		},
		// ldap dn is always the dn
		&cli.StringFlag{
			Name: "ldap-schema-account-id",
			// TODO write down LDAP schema & register OID ownclouduuid
			//... use 'sourceAnchor','immutableid' see https://docs.microsoft.com/en-us/azure/active-directory/hybrid/plan-connect-design-concepts#sourceanchor
			// or 'ms-DS-ConsistencyGuid' see https://docs.microsoft.com/en-us/azure/active-directory/hybrid/plan-connect-design-concepts
			// or build a scim schema for ldap? https://ldapwiki.com/wiki/SCIM%20Common%20Attribute
			// glauth -> support id and externalid from scim
			Value:       "uidNumber",
			Usage:       "LDAP account id attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_ACCOUNTID"},
			Destination: &cfg.LDAP.Schema.AccountID,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-username",
			Value:       "uid",
			Usage:       "LDAP username attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_USERNAME"},
			Destination: &cfg.LDAP.Schema.Username,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-password",
			Value:       "authPassword",
			Usage:       "LDAP password attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_PASSWORD"},
			Destination: &cfg.LDAP.Schema.Password,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-uid",
			Value:       "uidnumber",
			Usage:       "LDAP uid attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_UID"},
			Destination: &cfg.LDAP.Schema.UID,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-gid",
			Value:       "gidnumber",
			Usage:       "LDAP gid attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_GID"},
			Destination: &cfg.LDAP.Schema.GID,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-displayname",
			Value:       "displayname",
			Usage:       "LDAP displayname attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.LDAP.Schema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-surname",
			Value:       "sn",
			Usage:       "LDAP surname attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_SURNAME"},
			Destination: &cfg.LDAP.Schema.Surname,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-givenname",
			Value:       "sn",
			Usage:       "LDAP givenname attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_GIVENNAME"},
			Destination: &cfg.LDAP.Schema.GivenName,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-mail",
			Value:       "mail",
			Usage:       "LDAP mail attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_MAIL"},
			Destination: &cfg.LDAP.Schema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-memberof",
			Value:       "memberof",
			Usage:       "LDAP memberof attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_MEMBEROF"},
			Destination: &cfg.LDAP.Schema.Groups,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-description",
			Value:       "description",
			Usage:       "LDAP description attribute",
			EnvVars:     []string{"ACCOUNTS_LDAP_SCHEMA_Description"},
			Destination: &cfg.LDAP.Schema.Description,
		},
	}
}
