package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <ticket_id>",
	Short: "归档工单",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string]bool{"archive": true}
		data, err := cli.Post("/tickets/"+args[0]+"/archive", body, false)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var unarchiveCmd = &cobra.Command{
	Use:   "unarchive <ticket_id>",
	Short: "取消归档工单",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := cli.Post("/tickets/"+args[0]+"/unarchive", nil, false)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd, unarchiveCmd)
}
