package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	answerText string
	answerFile string
)

var answerCmd = &cobra.Command{
	Use:   "answer <ticket_id>",
	Short: "答复工单（可覆盖已有答复）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := answerText
		if text == "" && answerFile != "" {
			data, err := os.ReadFile(answerFile)
			if err != nil {
				return fmt.Errorf("读取文件失败: %w", err)
			}
			text = string(data)
		}
		if text == "" {
			return fmt.Errorf("请提供 --text 或 --file")
		}
		body := map[string]string{"answer": text}
		data, err := cli.Post("/tickets/"+args[0]+"/answer", body, false)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var closeCmd = &cobra.Command{
	Use:   "close <ticket_id>",
	Short: "关闭工单",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string]string{}
		if closeReason != "" {
			body["reason"] = closeReason
		}
		data, err := cli.Post("/tickets/"+args[0]+"/close", body, false)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var closeReason string

var reopenCmd = &cobra.Command{
	Use:   "reopen <ticket_id>",
	Short: "重开工单",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := cli.Post("/tickets/"+args[0]+"/reopen", map[string]string{}, false)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	answerCmd.Flags().StringVar(&answerText, "text", "", "答复内容")
	answerCmd.Flags().StringVar(&answerFile, "file", "", "从文件读取答复")
	closeCmd.Flags().StringVar(&closeReason, "reason", "", "关闭原因")
	rootCmd.AddCommand(answerCmd, closeCmd, reopenCmd)
}
