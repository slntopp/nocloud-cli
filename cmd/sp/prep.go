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
package sp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
	"sigs.k8s.io/yaml"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud-proto/services_providers"
)

var PrepCmd = &cobra.Command{
	Use:   "prep [path to template] [flags]",
	Short: "Prepare SP template by gathering data",
	Args:  cobra.ExactArgs(1),
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

		var request pb.PrepSP
		err = json.Unmarshal(template, &request)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		ctx, client := MakeServicesProviderServiceClientOrFail()
		res, err := client.Prep(ctx, &request)
		if err != nil {
			return err
		}

		ok, _ := tools.PrintJsonDataQ(cmd, res)
		if ok {
			return nil
		}

		out, err := yaml.Marshal(res)
		if err != nil {
			return err
		}

		fmt.Println(string(out))

		return nil
	},
}

var PrepIONeCmd = &cobra.Command{
	Use:   "ione [endpoint] [username] [password|token] [flags]",
	Short: "Gather data required for importing OpenNebula as SP",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		request := &pb.PrepSP{
			Sp: &pb.ServicesProvider{
				Type: "ione",
				Secrets: map[string]*structpb.Value{
					"host": structpb.NewStringValue(args[0]),
					"user": structpb.NewStringValue(args[1]),
					"pass": structpb.NewStringValue(args[2]),
				},
			},
		}

		ctx, client := MakeServicesProviderServiceClientOrFail()
		res, err := client.Prep(ctx, request)
		if err != nil {
			return err
		}

		ok, _ := tools.PrintJsonDataQ(cmd, res)
		if ok {
			return nil
		}

		out, err := yaml.Marshal(res)
		if err != nil {
			return err
		}

		fmt.Println(string(out))

		return nil
	},
}
