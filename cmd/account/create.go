/*
Copyright © 2021-2022 Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package account

import (
	"errors"
	"fmt"
	"strings"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	accountspb "github.com/slntopp/nocloud-proto/registry/accounts"
	"github.com/spf13/cobra"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:     "create [title] [namespace-uuid] [[flags]]",
	Aliases: []string{"crt", "c"},
	Short:   "Create NoCloud Account",
	Long:    "Authorization data flags must be given('auth-type', 'auth-data')",
	Args:    cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeAccountsServiceClientOrFail()

		authType, _ := cmd.Flags().GetString("auth-type")
		authData, _ := cmd.Flags().GetStringSlice("auth-data")
		if strings.Join(authData, "") == "" {
			return errors.New("args: Authorization Data wasn't given")
		}
		credentials := accountspb.Credentials{
			Type: (authType), Data: authData,
		}

		access, _ := cmd.Flags().GetInt32("access-level")

		req := accountspb.CreateRequest{
			Title: args[0], Namespace: args[1],
			Auth: &credentials, Access: &access,
		}
		res, err := client.Create(ctx, &req)
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("UUID:", res.GetUuid())
		}

		return nil
	},
}

func init() {
	CreateCmd.Flags().String("auth-type", "standard", "Authorization Type")
	CreateCmd.Flags().StringSlice("auth-data", []string{"", ""}, "Authorization Data(Credentials) as comma separated list. For example, =username,password")
	CreateCmd.Flags().Int32("access-level", 1, "New Account Access level to Namespace")
}
