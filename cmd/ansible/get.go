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
package ansible

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	pb "github.com/slntopp/nocloud-proto/ansible"
)

var GetCmd = &cobra.Command{
	Use:     "get [uuid]",
	Aliases: []string{"get", "fetch"},
	Short:   "Get Ansible Run",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		uuid := args[0]
		if uuid == "" {
			return errors.New("empty uuid")
		}

		ctx, client := MakeAnsibleServiceCleintOrFail()

		res, err := client.Get(ctx, &pb.GetRunRequest{
			Uuid: uuid,
		})
		if err != nil {
			return err
		}

		output, err := json.MarshalIndent(res, "-", " ")
		if err != nil {
			fmt.Println(res)
			return err
		}

		fmt.Println("Result: ", string(output))

		return nil
	},
}
