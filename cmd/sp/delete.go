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

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/services_providers"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete [uuid] [[flags]]",
	Short: "Delete NoCloud Services Providers",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeServicesProviderServiceClientOrFail()
		request := pb.DeleteRequest{Uuid: args[0]}
		res, err := client.Delete(ctx, &request)

		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			if !res.GetResult() {
				fmt.Println("Can't delete ServicesProvider, some services are still UP and provisioned")
				fmt.Println("Services:")
				for _, uuid := range res.GetServices() {
					fmt.Println("- ", uuid)
				}
				return nil
			}
			fmt.Println("Done.")
		}
		return nil
	},
}
