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

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "get message",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx, client := helpers.MakeChatsServiceClientOrFail()

		var uuid string
		if uuid, err = cmd.Flags().GetString("uuid"); err != nil || uuid == "" {
			return errors.New("empty uuid")
		}
		message, err := client.GetChatMessage(ctx, &proto.GetChatMessageRequest{
			Uuid: uuid,
		})

		if err != nil {
			fmt.Printf("Error while fetching message %s. Reason: %v.\n", uuid, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, message)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly fetched message %v\n", message)
		}

		return err
	},
}

func init() {
	GetCmd.Flags().StringP("uuid", "u", "", "Message uuid")
}
