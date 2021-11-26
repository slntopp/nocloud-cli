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
package sp

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/slntopp/nocloud/pkg/api/apipb"
	spb "github.com/slntopp/nocloud/pkg/services_providers/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
)

func MakeServicesProviderServiceClientOrFail() (context.Context, pb.ServicesProvidersServiceClient){
	host := viper.Get("nocloud")
	if host == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Host is unset")
	}

	conn, err := grpc.Dial(host.(string), grpc.WithInsecure())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic(err)
	}

	token := viper.Get("token")
	if token == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Token is unset")
	}

	client := pb.NewServicesProvidersServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer " + token.(string))
	return ctx, client
}

func PrintServicesProvider(s *spb.ServicesProvider) error {
	out, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func PrintServicesProvidersPool(pool []*spb.ServicesProvider) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"UUID", "Title", "Type"})

	rows := make([]table.Row, len(pool))
	for i, s := range pool {
		rows[i] = table.Row{s.GetUuid(), s.GetTitle(), s.GetType()}
	}
	t.AppendRows(rows)

	t.AppendFooter(table.Row{"", "Total Found", len(pool)})
    t.Render()
}