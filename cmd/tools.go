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
	"crypto/md5"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Tools that might or going to be useful to work with NoCloud",
}

var hashCmd = &cobra.Command{
	Use: "hash",
	Short: "Generate Hash of various things like string, certs etc",
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var data []byte

		if cert, err := cmd.Flags().GetString("cert"); err == nil && cert != "" {
			r, err := ioutil.ReadFile(cert)
			fmt.Println(r)
			if err != nil {
				return err
			}
			block, _ := pem.Decode(r)
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return err
			}
			data = cert.Raw
		} else if str, err := cmd.Flags().GetString("string"); err == nil && str != "" {
			data = []byte(str)
		} else {
			return errors.New("Nothing to do or an Error occured while parsing flags")
		}

		var resultB []byte
		alg, _ := cmd.Flags().GetString("alg")
		switch alg {
		case "sha256":
			hash := sha256.Sum256(data)
			resultB = hash[:]
		case "md5":
			hash := md5.Sum(data)
			resultB = hash[:]
		default:
			return errors.New("Not supported Algorythm")
		}
		result := hex.EncodeToString(resultB)
		if printJson, _ := cmd.Flags().GetBool("json"); printJson {
			data, err := json.Marshal(map[string]string{
				"hash": string(result), "alg": alg,
			})
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			fmt.Println("Hash:", string(result))
		}

		return nil
	},
}

func init() {
	hashCmd.Flags().StringP("cert", "c", "", ".crt Certificate File to hash")
	hashCmd.Flags().StringP("string", "s", "", "String to hash")
	hashCmd.Flags().StringP("alg", "a", "sha256", "Algorythm to use for hashing")
	toolsCmd.AddCommand(hashCmd)

	rootCmd.AddCommand(toolsCmd)
}
