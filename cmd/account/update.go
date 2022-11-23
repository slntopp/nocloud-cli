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

	"github.com/slntopp/nocloud-cli/pkg/tools"
	accountspb "github.com/slntopp/nocloud-proto/registry/accounts"
	"github.com/spf13/cobra"
)

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:     "update [account UUID] [flags]",
	Aliases: []string{"upd"},
	Short:   "Update NoCloud Account",
	Long:    "In order to execute request you must change at least one field.",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeAccountsServiceClientOrFail()

		updated := false
		req := accountspb.Account{
			Uuid: args[0],
		}
		title, _ := cmd.Flags().GetString("title")
		if title != "" {
			req.Title = title
			updated = true
		}

		if !updated {
			return errors.New("no fields updated, exiting")
		}

		res, err := client.Update(ctx, &req)
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Result: %t\n", res.GetResult())
		}

		return nil
	},
}

func init() {
	UpdateCmd.Flags().String("title", "", "Change account title to given")
}
