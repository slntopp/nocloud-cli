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
package settings

import (
	"fmt"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/settings"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// GetCmd represents the list command
var GetCmd = &cobra.Command{
	Use:   "get [...keys]",
	Short: "Get NoCloud Settings",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeSettingsServiceClientOrFail()
		request := pb.GetRequest{
			Keys: args,
		}

		res, err := client.Get(ctx, &request)

		if err != nil {
			return err
		}

		result := res.AsMap()

		ok, err := tools.PrintJsonDataQ(cmd, result)
		if err != nil {
			return err
		}
		if !ok {
			data, err := yaml.Marshal(result)
			if err != nil {
				return err
			}
			fmt.Print(string(data))
		}

		return nil
	},
}
