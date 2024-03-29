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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	pb "github.com/slntopp/nocloud-proto/ansible"
)

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:     "create [path to template]",
	Aliases: []string{"crt", "c"},
	Short:   "Create Ansible Run",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return errors.New("Template doesn't exist at path " + args[0])
		}

		var format string
		{
			pathSlice := strings.Split(args[0], ".")
			format = pathSlice[len(pathSlice)-1]
		}

		template, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("Error reading template file")
			return err
		}

		switch format {
		case "json":
		case "yml", "yaml":
			template, err = yaml.YAMLToJSON(template)
		default:
			return errors.New("Unsupported template format " + format)
		}
		if err != nil {
			fmt.Println("Error while parsing template1")
			return err
		}

		run := &pb.Run{}
		err = json.Unmarshal(template, &run)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		ctx, client := MakeAnsibleServiceCleintOrFail()

		res, err := client.Create(ctx, &pb.CreateRunRequest{Run: run})
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
