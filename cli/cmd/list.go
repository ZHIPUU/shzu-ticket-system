package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"ticket-cli/internal/client"
	"ticket-cli/internal/output"
)

var (
	listStatus   string
	listCategory string
	listArchived string
	listQuery    string
	listStart    string
	listEnd      string
	listPage     int
	listSize     int
	listAll      bool
	listJSON     bool
	listVerbose  bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "工单列表（支持筛选）",
	Long: `查询工单列表，支持按状态/分类/归档/关键字/时间筛选。

示例:
  ticket list
  ticket list --status pending
  ticket list --category 宿舍
  ticket list --archived true
  ticket list --q "宿舍" --size 20
  ticket list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		q := client.BuildQuery(map[string]string{
			"status":   listStatus,
			"category": listCategory,
			"archived": listArchived,
			"q":        listQuery,
			"start":    listStart,
			"end":      listEnd,
		})
		if listAll {
			// 翻所有页
			page := listPage
			if page < 1 {
				page = 1
			}
			var allItems []output.Ticket
			for {
				q.Set("page", fmt.Sprintf("%d", page))
				q.Set("size", "100")
				data, err := cli.Get("/tickets", q, true)
				if err != nil {
					return err
				}
				var resp struct {
					Total    int64           `json:"total"`
					Items    []output.Ticket `json:"items"`
					Page     int             `json:"page"`
					PageSize int             `json:"page_size"`
				}
				if err := json.Unmarshal(data, &resp); err != nil {
					return err
				}
				allItems = append(allItems, resp.Items...)
				if page*100 >= int(resp.Total) || len(resp.Items) == 0 {
					break
				}
				page++
			}
			if listJSON {
				output.JSON(allItems)
			} else {
				fmt.Printf("共 %d 条\n", len(allItems))
				output.TableList(allItems)
			}
			return nil
		}

		q.Set("page", fmt.Sprintf("%d", listPage))
		q.Set("size", fmt.Sprintf("%d", listSize))
		data, err := cli.Get("/tickets", q, true)
		if err != nil {
			return err
		}
		if listJSON {
			fmt.Println(string(data))
			return nil
		}
		var resp struct {
			Total    int64           `json:"total"`
			Items    []output.Ticket `json:"items"`
			Page     int             `json:"page"`
			PageSize int             `json:"page_size"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}
		fmt.Printf("共 %d 条，当前第 %d 页（每页 %d 条）\n", resp.Total, resp.Page, resp.PageSize)
		if listVerbose {
			for _, t := range resp.Items {
				output.Pretty(t)
			}
		} else {
			output.TableList(resp.Items)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&listStatus, "status", "", "状态筛选 pending/answered/closed")
	listCmd.Flags().StringVar(&listCategory, "category", "", "分类筛选")
	listCmd.Flags().StringVar(&listArchived, "archived", "", "归档 true/false")
	listCmd.Flags().StringVarP(&listQuery, "q", "q", "", "关键字搜索工单号/问题")
	listCmd.Flags().StringVar(&listStart, "start", "", "开始日期 YYYY-MM-DD")
	listCmd.Flags().StringVar(&listEnd, "end", "", "结束日期 YYYY-MM-DD")
	listCmd.Flags().IntVar(&listPage, "page", 1, "页码")
	listCmd.Flags().IntVar(&listSize, "size", 20, "每页条数（最大 200）")
	listCmd.Flags().BoolVar(&listAll, "all", false, "翻所有页（自动分页）")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "JSON 输出")
	listCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "详细输出（每个工单完整字段）")
	rootCmd.AddCommand(listCmd)
}
