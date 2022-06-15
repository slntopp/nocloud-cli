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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/slntopp/nocloud-cli/cmd/sp"
	pb "github.com/slntopp/nocloud/pkg/services/proto"
	sppb "github.com/slntopp/nocloud/pkg/services_providers/proto"
)

func SelectDeployPoliciesInteractive(ctx context.Context, cmd *cobra.Command, client pb.ServicesServiceClient, service *pb.Service) (res map[int32]string, err error) {
	if err != nil {
		return nil, err
	}
	ctx, spClient := sp.MakeServicesProviderServiceClientOrFail()
	sps, err := spClient.List(ctx, &sppb.ListRequest{})
	if err != nil {
		return nil, err
	}
	providers := make(map[string][]string)
	for _, sp := range sps.GetPool() {
		pool := providers[sp.GetType()]
		if pool == nil {
			pool = make([]string, 0)
		}
		pool = append(pool, fmt.Sprintf("%s | %s", sp.GetTitle(), sp.GetUuid()))
		providers[sp.GetType()] = pool
	}

	res = make(map[int32]string)
	for i, group := range service.GetInstancesGroups() {
		p := promptui.Select{
			Label: fmt.Sprintf("Select Service Provider for Instances Group %s (%s)", group.Title, group.GetUuid()),
			Items: providers[group.GetType()],
		}
		_, selected, err := p.Run()
		if err != nil {
			return nil, err
		}
		selected = strings.Split(selected, " | ")[1]
		res[int32(i)] = selected
	}
	return res, nil
}

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:     "create [path to template] [flags]",
	Aliases: []string{"crt", "c"},
	Short:   "Create Service Config",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		if namespace == "" {
			return errors.New("Namespace UUID isn't given")
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
			fmt.Println("Error while parsing template1")
			return err
		}
		var service pb.Service
		err = json.Unmarshal(template, &service)
		if err != nil {
			fmt.Println("Error while parsing template2")
			return err
		}

		ctx, client := MakeServicesServiceClientOrFail()
		req := pb.CreateRequest{Service: &service, Namespace: namespace}

		if rulesJson, _ := cmd.Flags().GetString("rules"); rulesJson != "" {
			fmt.Println("Rules as string given", rulesJson)
			json.Unmarshal([]byte(rulesJson), &req.DeployPolicies)
		} else if rulesFile, _ := cmd.Flags().GetString("rules-file"); rulesFile != "" {
			fmt.Println("Rules as File given", rulesFile)
			rulesJson, err := os.ReadFile(rulesFile)
			if err != nil {
				return err
			}
			json.Unmarshal(rulesJson, &req.DeployPolicies)
		} else {
			fmt.Println("Nothing given, selecting in interactive mode")
			r, err := SelectDeployPoliciesInteractive(ctx, cmd, client, &service)
			if err != nil {
				return err
			}
			req.DeployPolicies = r
		}

		fmt.Println(req)

		res, err := client.Create(ctx, &req)
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
	CreateCmd.Flags().StringP("namespace", "n", "", "Namespace UUID (required)")
	CreateCmd.Flags().StringP("rules", "r", "", "Deploy rules")
	CreateCmd.Flags().StringP("rules-file", "f", "", "Deploy rules")
}
