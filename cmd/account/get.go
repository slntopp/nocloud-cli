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
	"encoding/json"
	"fmt"

	accountspb "github.com/slntopp/nocloud/pkg/registry/proto/accounts"
	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get [UUID]",
	Short: "Get NoCloud Account Data",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeAccountsServiceClientOrFail()
		res, err := client.Get(ctx, &accountspb.GetRequest{
			Uuid: args[0],
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
			PrintAccount(res)
		}

		return nil
	},
}
