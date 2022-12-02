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
package services

import (
	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/services"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:     "list [[NAMESPACE]] [[flags]]",
	Aliases: []string{"l"},
	Short:   "List NoCloud Services",
	Long:    `Add namespace UUID after list command, to filter services by namespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeServicesServiceClientOrFail()
		request := pb.ListRequest{}
		if len(args) > 0 {
			request.Namespace = &args[0]
		}

		a, _ := cmd.Flags().GetBool("show-deleted")
		if a {
			showDeletedValue := "true"
			request.ShowDeleted = &showDeletedValue
		}

		d, _ := cmd.Flags().GetInt32("depth")
		request.Depth = &d

		res, err := client.List(ctx, &request)

		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			PrintServicesPool(res.GetPool())
		}

		return nil
	},
}

func init() {
	ListCmd.Flags().BoolP("show-deleted", "a", false, "Show Deleted Services")
	ListCmd.Flags().Int32P("depth", "d", 4, "Services Search(Traversal) depth")
}
