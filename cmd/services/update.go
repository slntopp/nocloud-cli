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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	pb "github.com/slntopp/nocloud-proto/services"
)

// createCmd represents the create command
var UpdateCmd = &cobra.Command{
	Use:     "update [path to template] [[flags]]",
	Aliases: []string{"upd", "u"},
	Short:   "Update Service Config",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		if namespace == "" {
			return errors.New(" Namespace UUID isn't given")
		}

		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return errors.New("Template doesn't exist at path " + args[0])
		}

		var format string
		{
			pathSlice := strings.Split(args[0], ".")
			format = pathSlice[len(pathSlice)-1]
		}

		template, err := os.ReadFile(args[0])

		switch format {
		case "json":
		case "yml", "yaml":
			template, err = yaml.YAMLToJSON(template)
		default:
			return errors.New("Unsupported template format " + format)
		}

		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}
		var service pb.Service
		err = json.Unmarshal(template, &service)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		ctx, client := MakeServicesServiceClientOrFail()
		res, err := client.Update(ctx, &service)
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

func init() {
	UpdateCmd.Flags().StringP("namespace", "n", "", "Namespace UUID (required)")
}
