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

// sendCmd represents the send message command
var SendCmd = &cobra.Command{
	Use:   "send [flags]",
	Short: "Send message",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx, client := helpers.MakeChatsServiceClientOrFail()

		var messageText string
		if messageText, err = cmd.Flags().GetString("message"); err != nil || messageText == "" {
			return errors.New("message text is empty")
		}
		var reciever string
		if reciever, err = cmd.Flags().GetString("to"); err != nil || reciever == "" {
			return errors.New("reciever is empty")
		}
		var entities []string
		if entities, err = cmd.Flags().GetStringSlice("entities"); err != nil {
			return err
		}

		msg, err := client.SendChatMessage(ctx, &proto.SendChatMessageRequest{
			Message: &proto.ChatMessage{
				To:      reciever,
				Message: messageText,
			},
			Entities: entities,
		})

		if err != nil {
			fmt.Printf("Error while sending message %s. Reason: %v.\n", messageText, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, msg)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly sent message %v.\n", msg)
		}

		return err
	},
}

func init() {
	SendCmd.Flags().StringP("to", "t", "", "Reciever uuid")
	SendCmd.Flags().StringP("message", "m", "", "Message text")
	SendCmd.Flags().StringSliceP("entities", "e", []string{}, "Message text")
}
