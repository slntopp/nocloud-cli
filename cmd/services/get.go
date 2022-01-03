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
var GetCmd = &cobra.Command{
	Use:   "get [uuid] [[flags]]",
	Short: "Get NoCloud Service",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeServicesServiceClientOrFail()
		request := pb.GetRequest{Uuid: args[0]}
		res, err := client.Get(ctx, &request)

		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); !printJson {
			return PrintService(res)
		}
		
		data, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}
