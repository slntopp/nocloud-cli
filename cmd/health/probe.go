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
package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/slntopp/nocloud/pkg/health/proto"
	"github.com/spf13/cobra"
)

// ProbeCmd represents the Probe command
var ProbeCmd = &cobra.Command{
	Use:   "probe [probe_type]",
	Short: "Do health probe",
	Long: `Available probe types:
	* ping — Make PING probe. NoCloud should return PONG
	* services - Check if NoCloud microservices are up, resolvable and responding
	`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeHealthServiceClientOrFail()
		var res *pb.ProbeResponse
		var err error
		switch args[0] {
		case "ping":
			res, err = client.Probe(ctx, &pb.ProbeRequest{ProbeType: "ping"})
		case "services":
			return CheckServices(cmd, ctx, client)
		case "routines":
			return CheckRoutines(cmd, ctx, client)
		default:
			err = errors.New("Probe type " + args[0] + " not declared")
		}
		if err != nil {
			return err
		}

		fmt.Println("Probe Result:", res.Response)
		return nil
	},
}

func CheckServices(cmd *cobra.Command, ctx context.Context, client pb.HealthServiceClient) error {

	res, err := client.Probe(ctx, &pb.ProbeRequest{ProbeType: "services"})
	if err != nil {
		return err
	}

	if printJson, _ := cmd.Flags().GetBool("json"); printJson {
		data, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println("Probe Result: ", res.GetStatus().String())
	
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Service", "Status", "Error"})

	for _, service := range res.GetServing() {
		t.AppendRow(table.Row{service.GetService(), service.GetStatus().String(), service.GetError()})
	}

    t.Render()

	return nil
}

func CheckRoutines(cmd *cobra.Command, ctx context.Context, client pb.HealthServiceClient) error {

	res, err := client.Probe(ctx, &pb.ProbeRequest{ProbeType: "routines"})
	if err != nil {
		return err
	}

	if printJson, _ := cmd.Flags().GetBool("json"); printJson {
		data, err := json.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println("Probe Result: ", res.GetStatus().String())
	
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Service", "Routine", "Status", "Error"})

	for _, service := range res.GetRoutines() {
		t.AppendRow(table.Row{service.GetStatus().GetService(), service.GetRoutine(), service.GetStatus().GetStatus().String(), service.GetStatus().GetError()})
	}

    t.Render()

	return nil
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ProbeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ProbeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
