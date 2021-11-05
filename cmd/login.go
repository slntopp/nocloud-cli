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
	"fmt"

	"github.com/slntopp/nocloud/pkg/accounting/accountspb"
	pb "github.com/slntopp/nocloud/pkg/api/apipb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [host:port] [username] [password]",
	Short: "Authorize in NoCloud Platform API",
	Long: `Generate Auth Token in NoCloud API and store it in CLI config.`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(args[0], grpc.WithInsecure())
		if err != nil {
			return err
		}

		client := pb.NewAccountsServiceClient(conn)
		res, err := client.Token(context.Background(), &accountspb.TokenRequest{
			Auth: &accountspb.Credentials{
				Type: "standard", Data: []string{args[1], args[2]},
			},
		})
		if err != nil {
			return err
		}
		token := res.GetToken()
		printToken, _ := cmd.Flags().GetBool("print-token")
		if printToken {
			fmt.Println(token)
		}

		viper.Set("nocloud", args[0])
		viper.Set("token", token)
		err = viper.WriteConfig()
		return err
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().Bool("print-token", false, "")
}
