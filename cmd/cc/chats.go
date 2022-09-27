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
package cc

import (
	"github.com/slntopp/nocloud-cli/cmd/cc/chats"
	"github.com/spf13/cobra"
)

// chatsCmd represents the chats command
var ChatsCmd = &cobra.Command{
	Use:     "chat",
	Aliases: []string{"cht", "ch", "chats"},
	Short:   "Manage nocloud-cc chats",
}

func init() {
	ChatsCmd.AddCommand(chats.CreateCmd)
	ChatsCmd.AddCommand(chats.GetCmd)
	ChatsCmd.AddCommand(chats.DeleteCmd)
	ChatsCmd.AddCommand(chats.UpdateCmd)
	ChatsCmd.AddCommand(chats.GetCmd)
	ChatsCmd.AddCommand(chats.InviteCmd)
	ChatsCmd.AddCommand(chats.StreamCmd)
}
