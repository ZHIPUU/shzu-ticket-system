package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"ticket-cli/internal/output"
)

var getCmd = &cobra.Command{
	Use:   "get <ticket_id>",
	Short: "工单详情",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := cli.Get("/tickets/"+args[0], nil, true)
		if err != nil {
			return err
		}
		if listJSON {
			fmt.Println(string(data))
			return nil
		}
		var t output.Ticket
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		output.Pretty(t)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
