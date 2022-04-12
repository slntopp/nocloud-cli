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
package cmd

import (
	"github.com/slntopp/nocloud-cli/cmd/billing"
	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var billingCmd = &cobra.Command{
	Use:   "billing",
	Aliases: []string{"bill"},
	Short: "Manage Billing Plans",
}

func init() {
	billingCmd.AddCommand(billing.CreateCmd)
	billingCmd.AddCommand(billing.GetCmd)
	billingCmd.AddCommand(billing.ListCmd)
	billingCmd.AddCommand(billing.DeleteCmd)

	rootCmd.AddCommand(billingCmd)
}

