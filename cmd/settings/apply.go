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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/settings"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// ApplyCmd represents the list command
var ApplyCmd = &cobra.Command{
	Use:   "apply [path]",
	Short: "Add NoCloud Setting",
	Long: `Add NoCloud Setting from yaml/json config, like:
key: setting_key
value: some value
description: Just a setting Key
level: 0-4 setting visibility
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		var request pb.PutRequest
		err = json.Unmarshal(template, &request)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		ctx, client := MakeSettingsServiceClientOrFail()
		res, err := client.Put(ctx, &request)
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Added:", res.GetKey())
		}

		return nil
	},
}
