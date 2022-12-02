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
package chats

import (
	"errors"
	"fmt"

	"github.com/slntopp/nocloud-cli/cmd/cc/helpers"
	"github.com/slntopp/nocloud-cli/pkg/tools"
	proto "github.com/slntopp/nocloud-proto/cc"
	"github.com/spf13/cobra"
)

// InviteCmd represents the invite command
var InviteCmd = &cobra.Command{
	Use:   "invite [flags]",
	Short: "Invite to chat",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx, client := helpers.MakeChatsServiceClientOrFail()
		var account string
		if account, err = cmd.Flags().GetString("account"); err != nil || account == "" {
			return errors.New("empty account uui")
		}
		var chat string
		if chat, err = cmd.Flags().GetString("chat"); err != nil || chat == "" {
			return errors.New("empty chat uuid")
		}

		_, err = client.Invite(ctx, &proto.InviteChatRequest{
			ChatUuid: chat,
			UserUuid: account,
		})

		if err != nil {
			fmt.Printf("Error while inviting to chat %s. Reason: %v.\n", chat, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, chat)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly invited account %s to chat %s.\n", account, chat)
		}

		return err
	},
}

func init() {
	InviteCmd.Flags().StringP("account", "a", "", "Account uuid")
	InviteCmd.Flags().StringP("chat", "c", "", "Chat uuid")
}
