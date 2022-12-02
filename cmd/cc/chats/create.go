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

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create [flags]",
	Short: "Create chat",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx, client := helpers.MakeChatsServiceClientOrFail()
		var title string
		if title, err = cmd.Flags().GetString("title"); err != nil || title == "" {
			return errors.New("empty uuid")
		}

		chat, err := client.CreateChat(ctx, &proto.CreateChatRequest{
			Chat: &proto.Chat{
				Title: title,
			},
		})

		if err != nil {
			fmt.Printf("Error while creating chat %s. Reason: %v.\n", title, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, chat)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly created chat %s.\n", chat.GetUuid())
		}

		return err
	},
}

func init() {
	CreateCmd.Flags().StringP("title", "t", "", "Chat title")
}
