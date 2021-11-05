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
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/slntopp/nocloud/pkg/accounting/accountspb"
	pb "github.com/slntopp/nocloud/pkg/api/apipb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts, prints info about current by default",
	Long: ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client, err := MakeAccountsServiceClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error setting connection up")
			panic(err)
		}

		res, err := client.Get(ctx, &accountspb.GetRequest{
			Id: "me",
		})
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(res)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			PrintAccount(res)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}

func MakeAccountsServiceClient() (context.Context, pb.AccountsServiceClient, error){
	host := viper.Get("nocloud")
	if host == nil {
		return nil, nil, errors.New("Host is unset")
	}

	conn, err := grpc.Dial(host.(string), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	token := viper.Get("token")
	if token == nil {
		return nil, nil, errors.New("Token is unset")
	}

	client := pb.NewAccountsServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer " + token.(string))
	return ctx, client, nil
}

func PrintAccount(acc *accountspb.Account) {
	fmt.Println()

	fmt.Println("ID:", acc.GetId())
	fmt.Println("Title:", acc.GetTitle())
}