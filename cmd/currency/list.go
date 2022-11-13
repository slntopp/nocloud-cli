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
	"fmt"

	"github.com/slntopp/nocloud-cli/pkg/tools"
	"github.com/slntopp/nocloud/pkg/billing/proto"
	"github.com/spf13/cobra"
)

// ListCurrenciesCmd represents the list currencies command
var ListCurrenciesCmd = &cobra.Command{
	Use:     "currencies",
	Aliases: []string{"curs"},
	Short:   "List currencies",
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeCurrencyServiceClientOrFail()
		curs, err := client.GetCurrencies(ctx, &proto.GetCurrenciesRequest{})
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, curs)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Successfully fetched ", curs)
		}
		return nil
	},
}

// ListRatesCmd represents the list rates command
var ListRatesCmd = &cobra.Command{
	Use:     "rates",
	Aliases: []string{"exchangerates"},
	Short:   "List exchange rates",
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, client := MakeCurrencyServiceClientOrFail()
		rates, err := client.GetExchangeRates(ctx, &proto.GetExchangeRatesRequest{})
		if err != nil {
			return err
		}

		ok, err := tools.PrintJsonDataQ(cmd, rates)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("Successfully fetched ", rates)
		}
		return nil
	},
}
