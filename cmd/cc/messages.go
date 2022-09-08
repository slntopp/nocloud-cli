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
	"github.com/slntopp/nocloud-cli/cmd/cc/messages"
	"github.com/spf13/cobra"
)

// messageCmd represents the message command
var MessagesCmd = &cobra.Command{
	Use:     "messages",
	Aliases: []string{"msg", "message"},
	Short:   "Manage nocloud-cc messages",
}

func init() {
	MessagesCmd.AddCommand(messages.SendCmd)
	MessagesCmd.AddCommand(messages.DeleteCmd)
	MessagesCmd.AddCommand(messages.UpdateCmd)
	MessagesCmd.AddCommand(messages.ListCmd)
	MessagesCmd.AddCommand(messages.GetCmd)
}
