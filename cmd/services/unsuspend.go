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
	pb "github.com/slntopp/nocloud-proto/services"
	"github.com/spf13/cobra"
)

// Unsuspend represents the unsuspend command
var UnsuspendCmd = &cobra.Command{
	Use:     "unsuspend [service_id] [[flags]]",
	Aliases: []string{"unsus, uns, unsusp"},
	Short:   "NoCloud Service Unsuspend",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		uuid := args[0]
		ctx, client := MakeServicesServiceClientOrFail()

		req := pb.UnsuspendRequest{Uuid: uuid}

		res, err := client.Unsuspend(ctx, &req)
		if err != nil {
			fmt.Printf("Error while unsuspending Service %s. Reason: %v.\n", uuid, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly unsuspended Service %s.\n", uuid)
		}

		return nil
	},
}
