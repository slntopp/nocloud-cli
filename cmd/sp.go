/*
Copyright © 2021 Nikita Ivanovski info@slnt-opp.xyz

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
package cmd

import (
	"github.com/slntopp/nocloud-cli/cmd/sp"
	"github.com/spf13/cobra"
)

// spCmd represents the sp command
var spCmd = &cobra.Command{
	Use:   "sp",
	Short: "Manage Services Providers | Doesn't do anything by default",
}

func init() {
	spCmd.AddCommand(sp.CreateCmd)
	spCmd.AddCommand(sp.TestCmd)
	
	spCmd.AddCommand(sp.GetCmd)
	spCmd.AddCommand(sp.ListCmd)

	spCmd.AddCommand(sp.ExtCmd)

	rootCmd.AddCommand(spCmd)
}
