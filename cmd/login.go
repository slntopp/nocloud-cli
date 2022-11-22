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
package cmd

import (
	"context"
	"crypto/tls"
	"fmt"

	tools "github.com/slntopp/nocloud-cli/pkg/tools"
	regpb "github.com/slntopp/nocloud-proto/registry"
	pb "github.com/slntopp/nocloud-proto/registry/accounts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [host:port] [username] [password]",
	Short: "Authorize in NoCloud Platform API",
	Long:  `Generate Auth Token in NoCloud API and store it in CLI config.`,
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		opt := grpc.WithTransportCredentials(creds)
		_insecure, _ := cmd.Flags().GetBool("insecure")
		if _insecure {
			opt = grpc.WithTransportCredentials(insecure.NewCredentials())
		}
		conn, err := grpc.Dial(args[0], opt)
		if err != nil {
			return err
		}

		client := regpb.NewAccountsServiceClient(conn)
		authType, _ := cmd.Flags().GetString("auth-type")
		req := &pb.TokenRequest{
			Auth: &pb.Credentials{
				Type: authType, Data: args[1:],
			},
		}
		if rootClaim, _ := cmd.Flags().GetBool("root-claim"); rootClaim {
			req.RootClaim = true
		}
		res, err := client.Token(context.Background(), req)
		if err != nil {
			return err
		}
		token := res.GetToken()
		ok, _ := tools.PrintJsonDataQ(cmd, map[string]string{"token": token})
		if !ok {
			fmt.Println(token)
		}

		viper.Set("nocloud", args[0])
		viper.Set("token", token)
		viper.Set("insecure", _insecure)
		err = viper.WriteConfig()
		return err
	},
}

func init() {
	loginCmd.Flags().String("auth-type", "standard", "Type of Credentials to be used")
	loginCmd.Flags().Bool("print-token", false, "")
	loginCmd.Flags().Bool("root-claim", true, "")
	loginCmd.Flags().Bool("insecure", false, "Use WithInsecure instead of TLS")

	rootCmd.AddCommand(loginCmd)
}
