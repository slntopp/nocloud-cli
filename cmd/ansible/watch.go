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
package ansible

import (
	"fmt"
	"io"
	"log"

	pb "github.com/slntopp/nocloud-ansible/proto/ansible"
	"github.com/spf13/cobra"
)

var WatchCmd = &cobra.Command{
	Use:     "watch [service_id]",
	Aliases: []string{},
	Short:   "Watch Run",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx, client := MakeAnsibleServiceCleintOrFail()

		req := pb.WatchRunRequest{Uuid: args[0]}

		resp, err := client.Watch(ctx, &req)
		if err != nil {
			log.Fatal(err)
			return err
		}

		for {
			respObject, err := resp.Recv()
			if err == io.EOF {
				fmt.Println("All done")
				return nil
			}
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println(respObject)
		}
	},
}
