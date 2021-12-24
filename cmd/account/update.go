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
package account

import (
	"encoding/json"
	"errors"
	"fmt"

	accountspb "github.com/slntopp/nocloud/pkg/registry/proto/accounts"
	"github.com/spf13/cobra"
)

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use: "update [account UUID] [flags]",
	Short: "Update NoCloud Account",
	Long: "In order to execute request you must change at least one field.",
	Args: cobra.MinimumNArgs(1),
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
			return errors.New("No fields updated, exiting")
		}

		res, err := client.Update(ctx, &req)
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
			fmt.Printf("Result: %t\n", res.GetResult())
		}

		return nil
	},
}

func init() {
	UpdateCmd.Flags().String("title", "", "Change account title to given")
}