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
package billing

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	pb "github.com/slntopp/nocloud/pkg/billing/proto"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
)

var TransactionsCmd = &cobra.Command{
	Use:     "transaction",
	Aliases: []string{"t", "transactions", "tr", "trx"},
	Short:   "Get Records for Transaction",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, client := MakeBillingServiceClientOrFail()

		res, err := client.GetRecords(ctx, &pb.Transaction{
			Uuid: args[0],
		})

		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			PrintRecords(res.Pool)
		}

		return nil
	},
}

// ListCmd represents the list command
var ListTransactionsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List NoCloud Transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeBillingServiceClientOrFail()
		request := pb.GetTransactionsRequest{}

		account, _ := cmd.Flags().GetString("account")
		if account != "" {
			request.Account = &account
		}

		service, _ := cmd.Flags().GetString("service")
		if service != "" {
			request.Service = &service
		}

		res, err := client.GetTransactions(ctx, &request)

		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			meta, _ := cmd.Flags().GetBool("meta")
			PrintTransactions(res.GetPool(), meta)
		}

		return nil
	},
}

var CreateTransactionCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c", "new", "crt", "add"},
	Short:   "Create the new Transaction",
	RunE: func(cmd *cobra.Command, args []string) error {
		acc, _ := cmd.Flags().GetString("account")
		if acc == "" {
			return fmt.Errorf("Account is required")
		}

		total, _ := cmd.Flags().GetFloat64("total")
		if total == 0 {
			return fmt.Errorf("Total is required and must be not null")
		}

		meta := make(map[string]*structpb.Value)
		raw_meta, _ := cmd.Flags().GetString("meta")
		if raw_meta != "" {
			err := json.Unmarshal([]byte(raw_meta), &meta)
			if err != nil {
				return err
			}
		}

		exec, _ := cmd.Flags().GetInt64("exec")
		if delta, _ := cmd.Flags().GetString("delta"); delta != "" {
			re := regexp.MustCompile(`(?P<sign>^-?)(?P<num>\d+)(?P<mult>[smhdw]{1})`)
			matches := re.FindStringSubmatch(delta)
			if len(matches) == 0 {
				return errors.New("Invalid delta format")
			}
			multiplicators := map[string]int{
				"s": 1,
				"m": 60,
				"h": 3600,
				"d": 86400,
				"w": 604800,
			}
			t, err := strconv.Atoi(matches[2])
			if err != nil {
				return err
			}
			d := multiplicators[matches[3]] * t
			if matches[1] == "-" {
				d *= -1
			}

			exec += int64(d)
		}

		ctx, client := MakeBillingServiceClientOrFail()
		r, err := client.CreateTransaction(ctx, &pb.Transaction{
			Account: acc,
			Total:   total,
			Meta:    meta,
			Exec:    exec,
		})
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(r)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			meta, _ := cmd.Flags().GetBool("meta")
			PrintTransactions([]*pb.Transaction{r}, meta)
		}

		return nil
	},
}

var ReprocessTransactionsCmd = &cobra.Command{
	Use:     "reprocess",
	Aliases: []string{"r", "re", "rep", "repr", "repro"},
	Short:   "Reprocess transactions for Account",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeBillingServiceClientOrFail()
		r, err := client.Reprocess(ctx, &pb.ReprocessTransactionsRequest{
			Account: args[0],
		})
		if err != nil {
			return err
		}

		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(r)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			meta, _ := cmd.Flags().GetBool("meta")
			PrintTransactions(r.GetPool(), meta)
		}

		return nil
	},
}

func init() {
	ListTransactionsCmd.Flags().StringP("account", "a", "", "Account to get transactions for")
	ListTransactionsCmd.Flags().StringP("service", "s", "", "Service to get transactions for")
	ListTransactionsCmd.Flags().Bool("meta", false, "Show Transactions metadata")
	TransactionsCmd.AddCommand(ListTransactionsCmd)

	CreateTransactionCmd.Flags().Int64P("exec", "e", time.Now().Unix(), "Transaction Planned Execution time")
	CreateTransactionCmd.Flags().StringP("delta", "d", "", "Transaction Planned Execution time in form of delta from now, like 5m, 1h, 1d")
	CreateTransactionCmd.Flags().StringP("account", "a", "", "Account to make Account for")
	CreateTransactionCmd.Flags().Float64P("total", "t", 0.0, "Transaction Total, positive to be charged, negative to be refunded")
	CreateTransactionCmd.Flags().StringP("meta", "m", "", "Transaction metadata")
	TransactionsCmd.AddCommand(CreateTransactionCmd)

	ReprocessTransactionsCmd.Flags().Bool("meta", false, "Show Transactions metadata")
	TransactionsCmd.AddCommand(ReprocessTransactionsCmd)
}
