package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ticket-cli/internal/config"
)

// Client HTTP 客户端（API Key + JWT 双轨鉴权）
type Client struct {
	cfg    *config.Config
	HTTP   *http.Client
}

// New 构造客户端
func New(cfg *config.Config) *Client {
	return &Client{
		cfg:  cfg,
		HTTP: &http.Client{Timeout: 30 * time.Second},
	}
}

// setAuth 按需附加 API Key / JWT 头
func (c *Client) setAuth(req *http.Request, useAPIKey bool) {
	if useAPIKey && c.cfg.APIKey != "" {
		req.Header.Set("X-API-Key", c.cfg.APIKey)
		return
	}
	if c.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	}
}

func (c *Client) do(method, path string, body interface{}, query url.Values, useAPIKey bool) ([]byte, int, error) {
	u := c.cfg.APIBase + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuth(req, useAPIKey)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return data, resp.StatusCode, nil
}

// Error 错误响应
type Error struct {
	Status int    `json:"-"`
	Code   string `json:"error_code"`
	Msg    string `json:"error_message"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("HTTP %d: %s — %s", e.Status, e.Msg, e.Detail)
	}
	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Msg)
}

func (c *Client) Get(path string, query url.Values, useAPIKey bool) (json.RawMessage, error) {
	data, status, err := c.do(http.MethodGet, path, nil, query, useAPIKey)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseError(status, data)
	}
	return data, nil
}

func (c *Client) Post(path string, body interface{}, useAPIKey bool) (json.RawMessage, error) {
	return c.sendBody(http.MethodPost, path, body, useAPIKey)
}

func (c *Client) Patch(path string, body interface{}) (json.RawMessage, error) {
	return c.sendBody(http.MethodPatch, path, body, false)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.sendBody(http.MethodDelete, path, nil, false)
}

func (c *Client) sendBody(method, path string, body interface{}, useAPIKey bool) (json.RawMessage, error) {
	data, status, err := c.do(method, path, body, nil, useAPIKey)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, parseError(status, data)
	}
	return data, nil
}

func (c *Client) Download(path string, query url.Values, useAPIKey bool) ([]byte, string, error) {
	data, status, err := c.do(http.MethodGet, path, nil, query, useAPIKey)
	if err != nil {
		return nil, "", err
	}
	if status >= 400 {
		return nil, "", parseError(status, data)
	}
	// filename 简单解析
	disp := ""
	return data, disp, nil
}

func parseError(status int, data []byte) *Error {
	e := &Error{Status: status, Msg: "request failed"}
	_ = json.Unmarshal(data, e)
	if e.Msg == "" {
		e.Msg = string(data)
	}
	return e
}

// BuildQuery 构造 url.Values
func BuildQuery(m map[string]string) url.Values {
	q := url.Values{}
	for k, v := range m {
		if v != "" {
			q.Set(k, v)
		}
	}
	return q
}

// HasAPIKey 是否配置了 API Key
func (c *Client) HasAPIKey() bool { return strings.TrimSpace(c.cfg.APIKey) != "" }
