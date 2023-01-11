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
package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/events"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// GetCmd represents the list command
var PubCmd = &cobra.Command{
	Use:     "pub [template] [[flags]]",
	Aliases: []string{"publish", "send"},
	Short:   "Publishes Event",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeEventsServiceClientOrFail()

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
			if err != nil {
				return err
			}
		default:
			return errors.New("Unsupported template format " + format)
		}

		event := &pb.Event{}
		err = json.Unmarshal(template, &event)
		if err != nil {
			return err
		}

		t, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		} else {
			event.Type = t
		}

		uuid, err := cmd.Flags().GetString("uuid")
		if err != nil {
			return err
		} else {
			event.Uuid = uuid
		}

		response, err := client.Publish(ctx, event)
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, response)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Successfully sent event")
		}

		return nil
	},
}

func init() {
	PubCmd.Flags().StringP("type", "t", "", "Type of event")
	PubCmd.Flags().StringP("uuid", "u", "", "Uuid of recipient")
}
