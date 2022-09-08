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

	"github.com/slntopp/nocloud-cc/pkg/chats/proto"
	"github.com/slntopp/nocloud-cli/cmd/cc/helpers"
	"github.com/slntopp/nocloud-cli/pkg/tools"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update [flags]",
	Short: "Update message",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		var uuid string
		if uuid, err = cmd.Flags().GetString("uuid"); err != nil || uuid == "" {
			return errors.New("empty uuid")
		}
		var messageText string
		if messageText, err = cmd.Flags().GetString("message"); err != nil || messageText == "" {
			return errors.New("empty uuid")
		}

		ctx, client := helpers.MakeChatsServiceClientOrFail()
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		resp, err := client.UpdateChatMessage(ctx, &proto.ChatMessage{
			Message: messageText,
		})

		if err != nil {
			fmt.Printf("Error while updating message. Reason: %v.\n", err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, resp)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly updated message %s.\n", uuid)
		}

		return err
	},
}

func init() {
	UpdateCmd.Flags().StringP("uuid", "u", "", "Message uuid")
	UpdateCmd.Flags().StringP("message", "t", "", "New message text")
}
