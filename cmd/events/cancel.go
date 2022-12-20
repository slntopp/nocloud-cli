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

	pb "github.com/slntopp/nocloud-proto/events"
	"github.com/spf13/cobra"
)

// GetCmd represents the list command
var CancelCmd = &cobra.Command{
	Use:     "cancel [[flags]]",
	Aliases: []string{"canc", "abort", "remove", "purge"},
	Short:   "Cance Event from Event Bus",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeEventsServiceClientOrFail()

		req := &pb.CancelRequest{}

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

		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		} else {
			req.Id = id
		}

		if _, err := client.Cancel(ctx, req); err != nil {
			return err
		}

		fmt.Println("Successfully removed event from the event bus.")

		return nil
	},
}

func init() {
	CancelCmd.Flags().StringP("type", "t", "", "Type of event")
	CancelCmd.Flags().StringP("uuid", "u", "", "Uuid of recipient")
	CancelCmd.Flags().StringP("id", "i", "", "Id of event")
}
