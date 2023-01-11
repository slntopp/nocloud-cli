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
	"github.com/slntopp/nocloud-cli/cmd/billing"
	"github.com/slntopp/nocloud-cli/cmd/currency"
	"github.com/spf13/cobra"
)

// billingCmd represents the billing command
var billingCmd = &cobra.Command{
	Use:     "billing",
	Aliases: []string{"bill"},
	Short:   "Manage Billing Plans",
}

// CurrenciesCmd represents the currencies command
var currenciesCmd = &cobra.Command{
	Use:     "currencies",
	Aliases: []string{"cur", "curs"},
	Short:   "Manage Currencies",
}

func init() {
	currenciesCmd.AddCommand(currency.GetCmd)
	currenciesCmd.AddCommand(currency.ListCurrenciesCmd)
	currenciesCmd.AddCommand(currency.ListRatesCmd)
	currenciesCmd.AddCommand(currency.CreateCmd)
	currenciesCmd.AddCommand(currency.DeleteCmd)
	currenciesCmd.AddCommand(currency.UpdateCmd)
	billingCmd.AddCommand(currenciesCmd)

	billingCmd.AddCommand(billing.CreateCmd)
	billingCmd.AddCommand(billing.GetCmd)
	billingCmd.AddCommand(billing.ListCmd)
	billingCmd.AddCommand(billing.DeleteCmd)

	billingCmd.AddCommand(billing.TransactionsCmd)

	rootCmd.AddCommand(billingCmd)
}
