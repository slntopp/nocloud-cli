/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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

	"sigs.k8s.io/yaml"

	pb "github.com/slntopp/nocloud-proto/services_providers"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update [path to template] [flags]",
	Short: "Update Services Provider Config",
	Args:  cobra.ExactArgs(1),
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

		fmt.Println("Template", string(template))
		var request pb.ServicesProvider
		err = json.Unmarshal(template, &request)
		if err != nil {
			fmt.Println("Error while parsing template")
			return err
		}

		ctx, client := MakeServicesProviderServiceClientOrFail()
		res, err := client.Get(ctx, &pb.GetRequest{Uuid: request.GetUuid()})
		if err != nil {
			return err
		}
		if res.GetUuid() == "" {
			errMsg := fmt.Sprintf("Service Provider with given Uuid %v is not found", request.GetUuid())
			return errors.New(errMsg)
		}

		res, err = client.Update(ctx, &request)
		if err != nil {
			return err
		}

		fmt.Println("Service Provider Updated, UUID:", res.GetUuid())
		return nil
	},
}
