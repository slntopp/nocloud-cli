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
package messages

import (
	"errors"
	"fmt"

	"github.com/slntopp/nocloud-cli/cmd/cc/helpers"
	"github.com/slntopp/nocloud-cli/pkg/tools"
	proto "github.com/slntopp/nocloud-proto/cc"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List messages from chat",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		var uuid string
		if uuid, err = cmd.Flags().GetString("uuid"); err != nil || uuid == "" {
			return errors.New("empty chat uuid")
		}

		ctx, client := helpers.MakeChatsServiceClientOrFail()
		resp, err := client.ListChatMessages(ctx, &proto.ListChatMessagesRequest{
			ChatUuid: uuid,
		})

		if err != nil {
			fmt.Printf("Error while getting messages from chat %s. Reason: %v.\n", uuid, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, resp)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly fetched messages %v.\n", resp.Messages)
		}

		return err
	},
}

func init() {
	ListCmd.Flags().StringP("uuid", "u", "", "Chat uuid from where to fetch messages")
}
