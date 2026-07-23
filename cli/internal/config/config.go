package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIBase        string `yaml:"api_base"`
	APIKey         string `yaml:"api_key"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	Token          string `yaml:"token"`
	TokenExpiresAt string `yaml:"token_expires_at"`
	OutputFormat   string `yaml:"output_format"` // table | json | pretty
}

// DefaultConfigPath 默认配置路径：~/.ticket-cli.yaml
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".ticket-cli.yaml"
	}
	return filepath.Join(home, ".ticket-cli.yaml")
}

// Load 从指定路径加载配置（不存在时返回默认配置）
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath()
	}
	cfg := &Config{
		APIBase:      "http://localhost:8000/api/v1",
		OutputFormat: "table",
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "table"
	}
	if cfg.APIBase == "" {
		cfg.APIBase = "http://localhost:8000/api/v1"
	}
	return cfg, nil
}

// Save 写入配置到指定路径
func (c *Config) Save(path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// CreateExample 在指定路径写入示例配置（不覆盖已有）
func CreateExample(path string) error {
	if path == "" {
		path = DefaultConfigPath()
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	cfg := &Config{
		APIBase:      "http://120.48.98.164:8080/api/v1",
		APIKey:       "sk-ticket-prod-shzu-2026",
		Username:     "admin",
		Password:     "admin123",
		OutputFormat: "table",
	}
	return cfg.Save(path)
}
