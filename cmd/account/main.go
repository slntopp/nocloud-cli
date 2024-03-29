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
package account

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	regpb "github.com/slntopp/nocloud-proto/registry"
	pb "github.com/slntopp/nocloud-proto/registry/accounts"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func MakeAccountsServiceClientOrFail() (context.Context, regpb.AccountsServiceClient) {
	host := viper.Get("nocloud")
	if host == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Host is unset")
	}

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	insec := viper.GetBool("insecure")
	if insec {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.Dial(host.(string), grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic(err)
	}

	token := viper.Get("token")
	if token == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Token is unset")
	}

	client := regpb.NewAccountsServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token.(string))
	return ctx, client
}

func PrintAccount(acc *pb.Account) {
	fmt.Println("UUID:", acc.GetUuid())
	fmt.Println("Title:", acc.GetTitle())
	if acc.Balance != nil {
		fmt.Printf("Balance: %.2f NCU\n", acc.GetBalance())
	}
}

func PrintAccountsPool(pool []*pb.Account) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"UUID", "Title", "Balance NCU", "Access", "Role", "Namespace"})

	rows := make([]table.Row, len(pool))
	for i, acc := range pool {
		balance := "-"
		if acc.Balance != nil {
			balance = fmt.Sprintf("%.2f", *acc.Balance)
		}
		a, r, n := "READ", "-", "-"
		if acc.Access != nil {
			a = acc.Access.Level.Enum().String()
			r = acc.Access.Role
			if r == "" {
				r = "-"
			}
			if acc.Access.Namespace != nil {
				n = *acc.Access.Namespace
			}
		}
		rows[i] = table.Row{acc.Uuid, acc.Title, balance, a, r, n}
	}
	t.AppendRows(rows)

	t.SortBy([]table.SortBy{
		{Name: "UUID", Mode: table.Asc},
	})

	t.AppendFooter(table.Row{"Total Found", len(pool)})
	t.Render()
}
