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
package account

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/slntopp/nocloud/pkg/accounting/accountspb"
	pb "github.com/slntopp/nocloud/pkg/api/apipb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func MakeAccountsServiceClientOrFail(cmd *cobra.Command) (context.Context, pb.AccountsServiceClient){
	host := viper.Get("nocloud")
	if host == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Host is unset")
	}

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	opt := grpc.WithTransportCredentials(creds)
	if r, _ := cmd.Flags().GetBool("insecure"); r {
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

	client := pb.NewAccountsServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer " + token.(string))
	return ctx, client
}

func PrintAccount(acc *accountspb.Account) {
	fmt.Println()

	fmt.Println("ID:", acc.GetId())
	fmt.Println("Title:", acc.GetTitle())
}

func PrintAccountsPool(pool []*accountspb.Account) {
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