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
	"fmt"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud/pkg/services/proto"
	"github.com/spf13/cobra"
)

// SuspendCmd represents the suspend command
var SuspendCmd = &cobra.Command{
	Use:     "suspend [service_id] [flags]",
	Aliases: []string{"sus", "susp"},
	Short:   "NoCloud Service Suspend",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		uuid := args[0]
		ctx, client := MakeServicesServiceClientOrFail()

		req := pb.SuspendRequest{Uuid: uuid}

		res, err := client.Suspend(ctx, &req)
		if err != nil {
			fmt.Printf("Error while suspending Service %s. Reason: %v.\n", uuid, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly suspended Service %s.\n", uuid)
		}

		return nil
	},
}
