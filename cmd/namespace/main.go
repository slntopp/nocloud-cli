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
package namespace

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	regpb "github.com/slntopp/nocloud/pkg/registry/proto"
	pb "github.com/slntopp/nocloud/pkg/registry/proto/namespaces"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func MakeNamespacesServiceClientOrFail() (context.Context, regpb.NamespacesServiceClient){
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

	client := regpb.NewNamespacesServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer " + token.(string))
	return ctx, client
}

func PrintNamespacesPool(pool []*pb.Namespace) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title"})
	
	rows := make([]table.Row, len(pool))
	for i, acc := range pool {
		rows[i] = table.Row{acc.Id, acc.Title}
	}
	t.AppendRows(rows)

	t.SortBy([]table.SortBy{
		{Name: "ID", Mode: table.Asc},
	})

	t.AppendFooter(table.Row{"Total Found", len(pool)})
    t.Render()
}