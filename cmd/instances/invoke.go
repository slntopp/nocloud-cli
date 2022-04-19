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
package instances

import (
	"encoding/json"
	"fmt"

	pb "github.com/slntopp/nocloud/pkg/instances/proto"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
)

// GetCmd represents the list command
var InvokeCmd = &cobra.Command{
	Use:   "invoke [uuid] [action] [[flags]]",
	Aliases: []string{"call", "perform"},
	Short: "Invokes NoCloud Instance Action",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeInstancesServiceClientOrFail()
		request := pb.InvokeRequest{
			Uuid: args[0],
			Method: args[1],
		}
		data, err := cmd.Flags().GetString("data")
		if err != nil {
			return err
		}
		
		if data != "" {
			var dataMap map[string]interface{}
			err = json.Unmarshal([]byte(data), &dataMap)
			if err != nil {
				return err
			}
			dataStruct, err := structpb.NewStruct(dataMap)
			if err != nil {
				return err
			}
			request.Params = dataStruct.GetFields()
		}

		res, err := client.Invoke(ctx, &request)

		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); !printJson {
			return PrintInstanceInvokeResponse(res)
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