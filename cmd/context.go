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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Aliases: []string{"ctx"},
	Short: "Print current NoCloud CLI Context",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := make(map[string]interface{})
		data["host"] = viper.Get("nocloud")
		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(data)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		for k, v := range data {
			fmt.Printf("%s: %v\n", strings.Title(k), v)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
