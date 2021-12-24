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
	"crypto/tls"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/slntopp/nocloud/pkg/services/proto"
	"gopkg.in/yaml.v2"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func MakeServicesServiceClientOrFail() (context.Context, pb.ServicesServiceClient){
	host := viper.Get("nocloud")
	if host == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Host is unset")
	}

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	opt := grpc.WithTransportCredentials(creds)
	insecure := viper.GetBool("insecure")
	if insecure {
		opt = grpc.WithInsecure()
	}
	conn, err := grpc.Dial(host.(string), opt)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic(err)
	}

	token := viper.Get("token")
	if token == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Token is unset")
	}

	client := pb.NewServicesServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer " + token.(string))
	return ctx, client
}

func PrintTestErrors(pool []*pb.TestConfigError) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Group", "Instance", "Error"})
	
	rows := make([]table.Row, len(pool))
	for i, err := range pool {
		rows[i] = table.Row{err.GetInstanceGroup(), err.GetInstance(), err.GetError()}
	}
	t.AppendRows(rows)

	t.AppendFooter(table.Row{"", "Total Found", len(pool)})
    t.Render()
}

func PrintService(s *pb.Service) error {
	out, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func PrintServicesPool(pool []*pb.Service) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"UUID", "Title", "Status"})

	rows := make([]table.Row, len(pool))
	for i, s := range pool {
		rows[i] = table.Row{s.GetUuid(), s.GetTitle(), s.GetStatus()}
	}
	t.AppendRows(rows)

	t.AppendFooter(table.Row{"", "Total Found", len(pool)})
    t.Render()
}