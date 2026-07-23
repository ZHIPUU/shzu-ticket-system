package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"ticket-cli/internal/config"
)

var (
	loginUser string
	loginPwd  string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录获取 JWT",
	Long:  "用用户名/密码登录保存 JWT 到配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		user := loginUser
		if user == "" {
			user = cfg.Username
		}
		pwd := loginPwd
		if pwd == "" {
			pwd = cfg.Password
		}
		if user == "" || pwd == "" {
			return fmt.Errorf("请提供 --username 和 --password，或先在配置文件中设置")
		}
		body := map[string]string{"username": user, "password": pwd}
		data, err := cli.Post("/auth/login", body, false)
		if err != nil {
			return err
		}
		var resp struct {
			Token     string    `json:"token"`
			ExpiresAt time.Time `json:"expires_at"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
		cfg.Token = resp.Token
		cfg.TokenExpiresAt = resp.ExpiresAt.Format(time.RFC3339)
		cfg.Username = user
		cfg.Password = pwd
		path := cfgFile
		if path == "" {
			path = config.DefaultConfigPath()
		}
		if err := cfg.Save(path); err != nil {
			return fmt.Errorf("保存配置失败: %w", err)
		}
		fmt.Printf("✓ 登录成功，token 已保存到 %s\n", path)
		fmt.Printf("  过期时间: %s\n", cfg.TokenExpiresAt)
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginUser, "username", "u", "", "用户名")
	loginCmd.Flags().StringVarP(&loginPwd, "password", "p", "", "密码")
	rootCmd.AddCommand(loginCmd)
}
