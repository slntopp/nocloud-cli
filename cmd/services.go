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
	"github.com/slntopp/nocloud-cli/cmd/services"
	"github.com/spf13/cobra"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "service",
	Aliases: []string{"srv", "services"},
	Short: "Manage NoCloud Services | Aliases: srv, services",
}

func init() {
	servicesCmd.AddCommand(services.TestCmd)
	servicesCmd.AddCommand(services.CreateCmd)
	
	servicesCmd.AddCommand(services.UpCmd)
	servicesCmd.AddCommand(services.DownCmd)

	servicesCmd.AddCommand(services.GetCmd)
	servicesCmd.AddCommand(services.ListCmd)
	servicesCmd.AddCommand(services.DeleteCmd)
	
	rootCmd.AddCommand(servicesCmd)
}
