package tools

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func PrintJsonDataQ(cmd *cobra.Command, data interface{}) (ok bool, err error) {
	if printJson, _ := cmd.Flags().GetBool("json"); !printJson {
		return false, nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	fmt.Println(string(b))
	return true, nil
}
