package command

import (
	"fmt"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/flagset"
	accounts "github.com/owncloud/ocis-accounts/pkg/proto/v0"
)

// AddAccount command creates a new account
func AddAccount(cfg *config.Config) *cli.Command {
	a := &accounts.Account{
		PasswordProfile: &accounts.PasswordProfile{},
	}
	return &cli.Command{
		Name:    "add",
		Usage:   "Create a new account",
		Aliases: []string{"create", "a"},
		Flags:   flagset.AddAccountWithConfig(cfg, a),
		Action: func(c *cli.Context) error {
			accSvcID := cfg.GRPC.Namespace + "." + cfg.Server.Name
			accSvc := accounts.NewAccountsService(accSvcID, grpc.NewClient())
			_, err := accSvc.CreateAccount(c.Context, &accounts.CreateAccountRequest{
				Account: a,
			})

			if err != nil {
				fmt.Println(fmt.Errorf("could not create account %w", err))
				return err
			}

			return nil
		}}
}
