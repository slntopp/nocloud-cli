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
	"fmt"

	"github.com/slntopp/nocloud/pkg/accounting/accountspb"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the get command
var DeleteCmd = &cobra.Command{
	Use:   "delete [UUID]",
	Short: "Delete NoCloud Account",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeAccountsServiceClientOrFail()
		res, err := client.Delete(ctx, &accountspb.DeleteRequest{
			Id: args[0],
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
			fmt.Printf("Result: %t\n", res.Result)
		}

		return nil
	},
}
