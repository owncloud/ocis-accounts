package command

import (
	"fmt"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/flagset"
	accounts "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"google.golang.org/genproto/protobuf/field_mask"
	"os"
)

// UpdateAccount command for modifying accounts including password policies
func UpdateAccount(cfg *config.Config) *cli.Command {
	a := &accounts.Account{
		PasswordProfile: &accounts.PasswordProfile{},
	}
	return &cli.Command{
		Name:      "update",
		Usage:     "Make changes to an existing account",
		ArgsUsage: "id",
		Flags:     flagset.UpdateAccountWithConfig(cfg, a),
		Before: func(c *cli.Context) error {
			if len(c.StringSlice("password_policies")) > 0 {
				// StringSliceFlag doesn't support Destination
				a.PasswordProfile.PasswordPolicies = c.StringSlice("password_policies")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				fmt.Println("Please provide an account-id")
				os.Exit(1)
			}

			a.Id = c.Args().First()
			accSvcID := cfg.GRPC.Namespace + "." + cfg.Server.Name
			accSvc := accounts.NewAccountsService(accSvcID, grpc.NewClient())
			_, err := accSvc.UpdateAccount(c.Context, &accounts.UpdateAccountRequest{
				Account:    a,
				UpdateMask: buildAccUpdateMask(c.FlagNames()),
			})

			if err != nil {
				fmt.Println(fmt.Errorf("could not update account %w", err))
				return err
			}

			return nil
		}}
}

// buildAccUpdateMask by mapping passed update flags to account fieldNames.
//
// The UpdateMask is passed with the update-request to the server so that
// only the modified values are transferred.
func buildAccUpdateMask(setFlags []string) *field_mask.FieldMask {
	var flagToPath = map[string]string{
		"enabled":                   "AccountEnabled",
		"displayname":               "DisplayName",
		"preferred-name":            "PreferredName",
		"uidnumber":                 "UidNumber",
		"gidnumber":                 "GidNumber",
		"mail":                      "Mail",
		"description":               "Description",
		"password":                  "PasswordProfile.Password",
		"password-policies":         "PasswordProfile.PasswordPolicies",
		"force-password-change":     "PasswordProfile.ForceChangePasswordNextSignIn",
		"force-password-change-mfa": "PasswordProfile.ForceChangePasswordNextSignInWithMfa",
	}

	updatedPaths := make([]string, 0)

	for _, v := range setFlags {
		if _, ok := flagToPath[v]; ok {
			updatedPaths = append(updatedPaths, flagToPath[v])
		}
	}

	return &field_mask.FieldMask{Paths: updatedPaths}
}
