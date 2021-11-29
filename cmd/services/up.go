/*
Copyright © 2021 Nikita Ivanovski info@slnt-opp.xyz

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
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/slntopp/nocloud-cli/cmd/sp"
	"github.com/slntopp/nocloud/pkg/api/apipb"
	pb "github.com/slntopp/nocloud/pkg/services/proto"
	sppb "github.com/slntopp/nocloud/pkg/services_providers/proto"
	"github.com/spf13/cobra"
)

func SelectDeployPoliciesInteractive(ctx context.Context, client apipb.ServicesServiceClient, id string) (res map[string]string, err error) {
	service, err := client.Get(ctx, &pb.GetRequest{Id: id})
	if err != nil {
		return nil, err
	}
	ctx, spClient := sp.MakeServicesProviderServiceClientOrFail()
	sps, err := spClient.List(ctx, &sppb.ListRequest{})
	if err != nil {
		return nil, err
	}
	providers := make(map[string][]string)
	for _, sp := range sps.GetServicesProviders() {
		pool := providers[sp.GetType()]
		if pool == nil {
			pool = make([]string, 0)
		}
		pool = append(pool, fmt.Sprintf("%s | %s", sp.GetTitle(), sp.GetUuid()))
		providers[sp.GetType()] = pool
	}

	res = make(map[string]string)
	for name, group := range service.GetInstancesGroups() {
		p := promptui.Select{
			Label: fmt.Sprintf("Select Service Provider for Instances Group %s (%s)", name, group.GetUuid()),
			Items: providers[group.GetType()],
		}
		_, selected, err := p.Run()
		if err != nil {
			return nil, err
		}
		selected = strings.Split(selected, " | ")[1]
		res[group.GetUuid()] = selected
	}
	return res, nil
}

// createCmd represents the create command
var UpCmd = &cobra.Command{
	Use:   "up [service_id] [flags]",
	Short: "NoCloud Service Up",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx, client := MakeServicesServiceClientOrFail()

		req := pb.UpRequest{Id: args[0]}
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
			r, err := SelectDeployPoliciesInteractive(ctx, client, args[0])
			if err != nil {
				return err
			}
			req.DeployPolicies = r
		}

		_, err = client.Up(ctx, &req)
		return err
	},
}

func init() {
	UpCmd.Flags().StringP("rules", "r", "", "Deploy rules")
	UpCmd.Flags().StringP("rules-file", "f", "", "Deploy rules")
}