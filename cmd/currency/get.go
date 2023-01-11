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
package currency

import (
	"errors"
	"fmt"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	proto "github.com/slntopp/nocloud-proto/billing"
	"github.com/spf13/cobra"
)

// GetCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get exchange rate",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			from int
			to   int
		)
		from, err := cmd.Flags().GetInt("from")
		if err != nil || from == -1 {
			return errors.New("empty from")
		}
		to, err = cmd.Flags().GetInt("to")
		if err != nil || to == -1 {
			return errors.New("empty to")
		}

		ctx, client := MakeCurrencyServiceClientOrFail()
		res, err := client.GetExchangeRate(ctx, &proto.GetExchangeRateRequest{
			From: proto.Currency(from),
			To:   proto.Currency(to),
		})
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, res)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Successfully fetched ", res)
		}

		return nil
	},
}

func init() {
	GetCmd.Flags().Int("from", -1, "")
	GetCmd.Flags().Int("to", -1, "")
}
