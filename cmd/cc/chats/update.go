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
package chats

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/slntopp/nocloud-cli/cmd/cc/helpers"
	"github.com/slntopp/nocloud-cli/pkg/tools"
	proto "github.com/slntopp/nocloud-proto/cc"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update [path to template] [flags]",
	Short: "Update chat",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return errors.New("template doesn't exist at path " + args[0])
		}

		ctx, client := helpers.MakeChatsServiceClientOrFail()

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
			return errors.New("unsupported template format " + format)
		}

		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		fmt.Println("Template", string(template))
		request := &proto.Chat{}
		err = json.Unmarshal(template, request)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		_, err = client.UpdateChat(ctx, request)

		if err != nil {
			fmt.Printf("Error while updating chat. Reason: %v.\n", err)
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, request)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("Successfuly updated chat %s.\n", request)
		}

		return err
	},
}

func init() {
	UpdateCmd.Flags().StringP("title", "t", "", "Chat title")
}
