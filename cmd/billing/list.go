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
package billing

import (
	"encoding/json"
	"fmt"

	pb "github.com/slntopp/nocloud/pkg/billing/proto"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"l", "ls"},
	Short: "List NoCloud Billing Plans",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeBillingServiceClientOrFail()
		request := pb.ListRequest{}

		res, err := client.ListPlans(ctx, &request)

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
			fmt.Println("--------------------------------")
			for _, plan := range res.GetPool() {
				PrintPlan(plan)
				fmt.Println("--------------------------------")
			}
		}

		return nil
	},
}