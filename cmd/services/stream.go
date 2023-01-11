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
	"log"

	pb "github.com/slntopp/nocloud-proto/services"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var StreamCmd = &cobra.Command{
	Use:     "stream [service_id] [[flags]]",
	Aliases: []string{},
	Short:   "NoCloud Service Stream",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx, client := MakeServicesServiceClientOrFail()

		req := pb.StreamRequest{Uuid: args[0]}

		resp, err := client.Stream(ctx, &req)
		if err != nil {
			log.Fatal(err)
			return err
		}

		for {
			respObject, err := resp.Recv()
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println(respObject)
		}
	},
}
