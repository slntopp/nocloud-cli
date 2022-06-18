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
	"io"

	pb "github.com/slntopp/nocloud/pkg/services/proto"
	"github.com/spf13/cobra"
)

var instUUIDs []string

var StateCmd = &cobra.Command{
	Use:     "state [UUID] [flags]",
	Aliases: []string{"st", "s"},
	Short:   "Instances state streaming",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		streamReq := &pb.StreamDataRequest{
			Uuid: args[0],
		}

		ctx, client := MakeServicesServiceClientOrFail()

		stream, err := client.StreamData(ctx, streamReq)
		if err != nil {
			return err
		}

		fmt.Printf("Stream started")

		for {
			state, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("State stream finished")
				break
			}
			if err != nil {
				return err
			}

			if len(instUUIDs) > 0 {
				for _, uuid := range instUUIDs {
					if uuid == state.Uuid {
						fmt.Printf("%s: %s", uuid, state.State.State.String())
					}
				}
			} else {
				fmt.Printf("%s: %s", state.Uuid, state.State.State.String())
			}
		}

		return nil
	},
}

func init() {
	StateCmd.Flags().StringArrayVarP(&instUUIDs, "instances", "i", []string{}, "Instance UUIDs filter")
}
