package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	submitQuestion string
	submitUserID   string
	submitRAG      string
	submitSource   string
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "提交工单（API Key 鉴权）",
	Long: `模拟智能体提交工单。需要先在配置文件中设置 api_key。

示例:
  ticket submit --question "宿舍分配标准" --user-id sess_001
  ticket submit --question "..." --rag "检索到的内容" --source hiagent_chat`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if submitQuestion == "" {
			return fmt.Errorf("请提供 --question")
		}
		if !cli.HasAPIKey() {
			return fmt.Errorf("未配置 API Key，请先在 ~/.ticket-cli.yaml 设置 api_key，或用 --api-key 指定")
		}
		body := map[string]interface{}{
			"question": submitQuestion,
			"user_id":  submitUserID,
			"source":   submitSource,
			"rag_result": submitRAG,
			"timestamp": time.Now().Format(time.RFC3339),
		}
		data, err := cli.Post("/tickets", body, true)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	submitCmd.Flags().StringVarP(&submitQuestion, "question", "Q", "", "工单问题（必填）")
	submitCmd.Flags().StringVar(&submitUserID, "user-id", "", "用户 ID（可选）")
	submitCmd.Flags().StringVar(&submitRAG, "rag", "", "RAG 检索结果（可选）")
	submitCmd.Flags().StringVar(&submitSource, "source", "hiagent_chat", "来源（默认 hiagent_chat）")
	_ = json.RawMessage{} // keep import
	rootCmd.AddCommand(submitCmd)
}
