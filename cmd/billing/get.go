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
	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/billing"
	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get [UUID]",
	Short: "Get NoCloud BillingPlan",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeBillingServiceClientOrFail()
		res, err := client.GetPlan(ctx, &pb.Plan{
			Uuid: args[0],
		})
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			PrintPlan(res)
		}

		return nil
	},
}
