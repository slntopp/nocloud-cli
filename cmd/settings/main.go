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
package settings

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/slntopp/nocloud-proto/settings"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func MakeSettingsServiceClientOrFail() (context.Context, pb.SettingsServiceClient) {
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

	client := pb.NewSettingsServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token.(string))
	return ctx, client
}

func PrintKeys(res *pb.KeysResponse) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Description", "Public?"})

	rows := make([]table.Row, len(res.GetPool()))
	for i, key := range res.GetPool() {
		rows[i] = table.Row{key.Key, key.Description, key.Public}
	}
	t.AppendRows(rows)
	t.AppendFooter(table.Row{"Total Found", len(rows)})
	t.Render()
}
