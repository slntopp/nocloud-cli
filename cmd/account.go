/*
Copyright © 2021 Nikita Ivanovski info@slnt-opp.xyz

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
package cmd

import (
	"encoding/json"
	"fmt"

	accountspb "github.com/slntopp/nocloud/pkg/registry/proto/accounts"
	"github.com/spf13/cobra"

	"github.com/slntopp/nocloud-cli/cmd/account"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts, prints info about current by default",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := account.MakeAccountsServiceClientOrFail()

		res, err := client.Get(ctx, &accountspb.GetRequest{
			Uuid: "me",
		})
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(res)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			account.PrintAccount(res)
		}

		return nil
	},
}

func init() {
	accountCmd.AddCommand(account.GetCmd)
	accountCmd.AddCommand(account.ListCmd)
	accountCmd.AddCommand(account.CreateCmd)
	accountCmd.AddCommand(account.UpdateCmd)
	accountCmd.AddCommand(account.DeleteCmd)

	rootCmd.AddCommand(accountCmd)
}

