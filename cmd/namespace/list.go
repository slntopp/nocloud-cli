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
package namespace

import (
	"encoding/json"
	"fmt"

	"github.com/slntopp/nocloud/pkg/accounting/namespacespb"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List NoCloud Namespaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeNamespacesServiceClientOrFail()
		request := namespacespb.ListRequest{}

		d, _ := cmd.Flags().GetInt32("depth")
		request.Depth = &d

		res, err := client.List(ctx, &request)

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
			PrintNamespacesPool(res.Pool)
		}

		return nil
	},
}

func init() {
	ListCmd.Flags().Int32P("depth", "d", 4, "Accounts Search(Traversal) depth")
}