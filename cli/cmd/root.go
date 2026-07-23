package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ticket-cli/internal/client"
	"ticket-cli/internal/config"
)

var (
	cfgFile    string
	flagAPIBase string
	flagAPIKey  string
	cfg        *config.Config
	cli        *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "ticket",
	Short: "工单系统命令行工具",
	Long: `ticket 是石河子大学 AI 迎新助手工单系统的 CLI 客户端。

常用命令:
  ticket login          登录获取 JWT
  ticket list           工单列表（支持筛选）
  ticket get <id>       工单详情
  ticket submit         提交工单（用 API Key 鉴权）
  ticket answer <id>    答复工单
  ticket close <id>     关闭工单
  ticket reopen <id>    重开工单
  ticket archive <id>   归档工单
  ticket delete <id>    删除工单
  ticket export         导出工单
  ticket batch-delete   批量删除`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}
		if flagAPIBase != "" {
			cfg.APIBase = flagAPIBase
		}
		if flagAPIKey != "" {
			cfg.APIKey = flagAPIKey
		}
		cli = client.New(cfg)
		return nil
	},
}

// Execute 入口
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "✗", err)
		return err
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径（默认 ~/.ticket-cli.yaml）")
	rootCmd.PersistentFlags().StringVar(&flagAPIBase, "api-base", "", "API 地址（覆盖配置）")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "API Key（覆盖配置）")
}

