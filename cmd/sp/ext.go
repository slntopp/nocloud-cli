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

// ExtCmd represents the ext command
var ExtCmd = &cobra.Command{
	Use:   "ext",
	Short: "Manage ServicesProviders extentions",
}

var listExtCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered ServicesProviders extentions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeServicesProviderServiceClientOrFail()
		req := &pb.ListRequest{}
		res, err := client.ListExtentions(ctx, req)
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			for key, ext := range res.GetTypes() {
				fmt.Println(key, " - ", ext)
			}
		}

		return nil
	},
}

func init() {
	ExtCmd.AddCommand(listExtCmd)
}
