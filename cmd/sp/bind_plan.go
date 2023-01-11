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
package sp

import (
	"fmt"

	pb "github.com/slntopp/nocloud-proto/services_providers"
	"github.com/spf13/cobra"
)

// BindPlanCmd represents the bind-plan Command
var BindPlanCmd = &cobra.Command{
	Use:   "bind-plan [sp-uuid] [plan-uuid] [[flags]]",
	Short: "Bind Billing Plan to Services Providers",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeServicesProviderServiceClientOrFail()
		request := pb.BindPlanRequest{Uuid: args[0], PlanUuid: args[1]}
		_, err := client.BindPlan(ctx, &request)
		if err != nil {
			return err
		}
		fmt.Println("Binding Completed")
		return nil
	},
}

// UnbindPlanCmd represents the unbind-plan Command
var UnbindPlanCmd = &cobra.Command{
	Use:   "unbind-plan [uuid] [plan_uuid] [[flags]]",
	Short: "Unbind Billing Plan",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeServicesProviderServiceClientOrFail()
		request := pb.UnbindPlanRequest{Uuid: args[0], PlanUuid: args[1]}
		_, err := client.UnbindPlan(ctx, &request)
		if err != nil {
			return err
		}
		fmt.Println("Unbinding Completed")
		return nil
	},
}
