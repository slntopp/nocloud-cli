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
package playbooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/ansible"
)

var CreateCmd = &cobra.Command{
	Use:     "create [path to template]",
	Aliases: []string{"crt", "c"},
	Short:   "Create Playbook",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

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

		playbook := &pb.Playbook{}
		err = json.Unmarshal(template, playbook)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		ctx, client := MakePlaybooksServiceClientOrFail()
		res, err := client.Create(ctx, &pb.CreatePlaybookRequest{
			Playbook: playbook,
		})
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, map[string]string{"uuid": res.Playbook.Uuid})
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Result:", res.Playbook.Uuid)
		}

		return nil
	},
}
