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
package events

import (
	"fmt"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/events"
	"github.com/spf13/cobra"
)

// GetCmd represents the list command
var SubCmd = &cobra.Command{
	Use:     "sub [[flags]]",
	Aliases: []string{"subscribe", "consume"},
	Short:   "Consume Events",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeEventsServiceClientOrFail()

		req := &pb.ConsumeRequest{}

		t, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		} else {
			req.Type = t
		}

		uuid, err := cmd.Flags().GetString("uuid")
		if err != nil {
			return err
		} else {
			req.Uuid = uuid
		}

		stream, err := client.Consume(ctx, req)
		if err != nil {
			return err
		}

		for {
			msg, err := stream.Recv()
			if err != nil {
				return err
			}

			ok, err := tools.PrintJsonDataQ(cmd, msg)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println(msg)
			}
		}
	},
}

func init() {
	SubCmd.Flags().StringP("type", "t", "", "Type of event")
	SubCmd.Flags().StringP("uuid", "u", "", "Uuid of recipient")
}
