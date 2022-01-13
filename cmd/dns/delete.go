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
package dns

import (
	"encoding/json"
	"fmt"

	pb "github.com/slntopp/nocloud/pkg/dns/proto"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the dump command
var DeleteCmd = &cobra.Command{
	Use:   "delete [zone]",
	Short: "Delete Zone config",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeDNSClientOrFail()
		
		zone := args[0]
		request := pb.Zone{Name: zone}
		res, err := client.Delete(ctx, &request)
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(res)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			fmt.Print("Keys deleted", res.GetResult())
		}

		return nil
	},
}