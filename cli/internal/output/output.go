package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

// Ticket 工单统一展示结构
type Ticket struct {
	TicketID   string    `json:"ticket_id"`
	Question   string    `json:"question"`
	UserID     string    `json:"user_id"`
	Source     string    `json:"source"`
	Status     string    `json:"status"`
	Category   string    `json:"category"`
	Archived   bool      `json:"archived"`
	Answer     *string   `json:"answer,omitempty"`
	AnsweredAt *string   `json:"answered_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Table 打印工单列表
func TableList(items []Ticket) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "工单号\t状态\t分类\t归档\t问题\t创建时间")
	for _, t := range items {
		q := t.Question
		if len(q) > 30 {
			q = q[:30] + "..."
		}
		cat := t.Category
		if cat == "" {
			cat = "-"
		}
		arch := "否"
		if t.Archived {
			arch = "是"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			t.TicketID, statusText(t.Status), cat, arch, q,
			t.CreatedAt.Format("2006-01-02 15:04"))
	}
}

// JSON 原始输出
func JSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

// Pretty 完整单工单详情
func Pretty(t Ticket) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  工单号:    %s\n", t.TicketID)
	fmt.Printf("  状态:      %s\n", statusText(t.Status))
	fmt.Printf("  分类:      %s\n", orDash(t.Category))
	fmt.Printf("  归档:      %s\n", yesno(t.Archived))
	fmt.Printf("  来源:      %s\n", orDash(t.Source))
	fmt.Printf("  提交人:    %s\n", orDash(t.UserID))
	fmt.Printf("  创建时间:  %s\n", t.CreatedAt.Format("2006-01-02 15:04:05"))
	if !t.UpdatedAt.IsZero() {
		fmt.Printf("  更新时间:  %s\n", t.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	if t.Answer != nil {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("  答复:")
		fmt.Println("  " + *t.Answer)
		if t.AnsweredAt != nil {
			fmt.Printf("  答复时间:  %s\n", *t.AnsweredAt)
		}
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  问题: %s\n", t.Question)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func statusText(s string) string {
	switch s {
	case "pending":
		return "待处理"
	case "answered":
		return "已答复"
	case "closed":
		return "已关闭"
	case "processing":
		return "处理中"
	default:
		return s
	}
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func yesno(b bool) string {
	if b {
		return "是"
	}
	return "否"
}
