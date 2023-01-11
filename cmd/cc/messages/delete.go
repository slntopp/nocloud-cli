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

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete [[flags]]",
	Short: "Delete message",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx, client := helpers.MakeChatsServiceClientOrFail()

		var uuid string
		if uuid, err = cmd.Flags().GetString("uuid"); err != nil || uuid == "" {
			return errors.New("empty uuid")
		}
		resp, err := client.DeleteChatMessage(ctx, &proto.DeleteChatMessageRequest{
			Uuid: uuid,
		})

		if err != nil {
			fmt.Printf("Error while deleting message %s. Reason: %v.\n", uuid, err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, resp)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly deleted message %s.\n", uuid)
		}

		return err
	},
}

func init() {
	DeleteCmd.Flags().StringP("uuid", "u", "", "Message uuid")
}
