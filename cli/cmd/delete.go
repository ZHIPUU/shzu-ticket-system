package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	deleteHard bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <ticket_id>",
	Short: "删除工单（软删）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "/tickets/" + args[0]
		if deleteHard {
			path += "?hard=true"
		}
		data, err := cli.Delete(path)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var batchDeleteCmd = &cobra.Command{
	Use:   "batch-delete <id1> <id2> ...",
	Short: "批量删除工单",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string][]string{"ticket_ids": args}
		data, err := cli.Post("/tickets/batch-delete", body, false)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

var exportFormat string
var exportOutput string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出工单（CSV / JSON）",
	Long: `按当前筛选条件导出工单。

示例:
  ticket export --format csv --output tickets.csv
  ticket export --format csv --status closed --output closed.csv
  ticket export --format json --output tickets.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		q := url.Values{}
		if exportFormat == "" {
			exportFormat = "csv"
		}
		q.Set("format", exportFormat)
		if listStatus != "" {
			q.Set("status", listStatus)
		}
		if listCategory != "" {
			q.Set("category", listCategory)
		}
		if listArchived != "" {
			q.Set("archived", listArchived)
		}
		if listStart != "" {
			q.Set("start_date", listStart)
		}
		if listEnd != "" {
			q.Set("end_date", listEnd)
		}
		data, err := cli.Get("/tickets/export", q, true)
		if err != nil {
			return err
		}
		if exportOutput != "" {
			if err := writeExport(exportOutput, data); err != nil {
				return err
			}
			fmt.Printf("✓ 已导出到 %s（%d 字节）\n", exportOutput, len(data))
			return nil
		}
		if exportFormat == "json" {
			var pretty interface{}
			_ = json.Unmarshal(data, &pretty)
			out, _ := json.MarshalIndent(pretty, "", "  ")
			fmt.Println(string(out))
		} else {
			fmt.Println(string(data))
		}
		return nil
	},
}

func writeExport(path string, data []byte) error {
	if strings.HasSuffix(path, ".json") {
		var pretty interface{}
		_ = json.Unmarshal(data, &pretty)
		out, _ := json.MarshalIndent(pretty, "", "  ")
		data = out
	}
	return os.WriteFile(path, data, 0644)
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteHard, "hard", false, "硬删（不可恢复）")
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "格式 csv/json")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "输出文件路径（默认输出到 stdout）")
	// 共享 list 的筛选 flag
	exportCmd.Flags().StringVar(&listStatus, "status", "", "状态筛选")
	exportCmd.Flags().StringVar(&listCategory, "category", "", "分类筛选")
	exportCmd.Flags().StringVar(&listArchived, "archived", "", "归档筛选")
	exportCmd.Flags().StringVar(&listStart, "start", "", "开始日期")
	exportCmd.Flags().StringVar(&listEnd, "end", "", "结束日期")
	rootCmd.AddCommand(deleteCmd, batchDeleteCmd, exportCmd)
}

