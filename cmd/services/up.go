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
	pb "github.com/slntopp/nocloud-proto/services"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var UpCmd = &cobra.Command{
	Use:     "up [service_id] [[flags]]",
	Aliases: []string{},
	Short:   "NoCloud Service Up",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx, client := MakeServicesServiceClientOrFail()

		req := pb.UpRequest{Uuid: args[0]}

		_, err = client.Up(ctx, &req)
		return err
	},
}
