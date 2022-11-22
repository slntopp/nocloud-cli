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
	"log"

	"github.com/jroimartin/gocui"
	"github.com/slntopp/nocloud-cli/cmd/account"
	"github.com/slntopp/nocloud-cli/cmd/cc/helpers"
	proto "github.com/slntopp/nocloud-proto/cc"
	"github.com/spf13/cobra"
)

type Buffer struct {
	limit    int
	messages []*proto.ChatMessage
}

func (b *Buffer) Push(msg *proto.ChatMessage) {
	if len(b.messages) == b.limit {
		b.messages = b.messages[1:]
	}

	b.messages = append(b.messages, msg)
}

// StreamCmd represents the create stream command
var StreamCmd = &cobra.Command{
	Use:   "stream [flags]",
	Short: "Create chat stream for fetching messages",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx, client := helpers.MakeChatsServiceClientOrFail()
		ctxAcc, clientAcc := account.MakeAccountsServiceClientOrFail()

		var uuid string
		if uuid, err = cmd.Flags().GetString("uuid"); err != nil || uuid == "" {
			return errors.New("empty uuid")
		}

		stream, err := client.Stream(ctx, &proto.ChatMessageStreamRequest{
			Uuid: uuid,
		})
		if err != nil {
			fmt.Printf("Error while creating chat stream %s. Reason: %v.\n", uuid, err)
			return err
		}

		ui, err := NewUI(ctx, ctxAcc, client, clientAcc, stream, uuid)
		if err != nil {
			return err
		}
		defer ui.Close()
		ui.SetManagerFunc(ui.layout)
		if err := ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
			log.Fatalln(err)
		}
		if err := ui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, ui.sendMsg); err != nil {
			log.Fatalln(err)
		}
		go ui.receiveMsg()
		if err = ui.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Fatalln(err)
		}

		return nil
	},
}

func init() {
	StreamCmd.Flags().StringP("uuid", "u", "", "Chat uuid")
}
