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
package services

import (
	"encoding/json"
	"fmt"

	pb "github.com/slntopp/nocloud/pkg/services/proto"
	"github.com/spf13/cobra"
)

// GetCmd represents the list command
var InvokeCmd = &cobra.Command{
	Use:   "invoke [uuid] [uuid] [uuid] [action] [[flags]]",
	Aliases: []string{"call", "perform"},
	Short: "Invokes NoCloud Service Group Instance Action",
	Long: `Invokes Instance method, requires fuul UUIDs path, so args are: <service uuid> <group uuid> <instance uuid> <action key> --meta <json data>`,
	Args: cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeServicesServiceClientOrFail()
		request := pb.PerformActionRequest{
			Service: args[0],
			Group: args[1],
			Instance: args[2],
			Action: args[3],
		}
		data, err := cmd.Flags().GetString("meta")
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(data), &request.Data)
		if err != nil {
			return err
		}

		res, err := client.PerformServiceAction(ctx, &request)

		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); !printJson {
			return PrintServiceActionResponse(res)
		}
		
		meta, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(string(meta))
		return nil
	},
}

func init() {
	InvokeCmd.Flags().StringP("data", "d", "", "Data to be sent along with invoke request. Must be Valid JSON")
}