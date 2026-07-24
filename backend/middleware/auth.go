package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth API Key 鉴权中间件
// 支持三种来源（按优先级）：
//  1. HTTP header: X-API-Key
//  2. Query string: ?api_key=xxx
//  3. Request body (JSON): {"x_api_key": "xxx", ...}
//
// 设计原因：不同 LLM 平台对 API Key 鉴权的实现方式不一致，
// 部分平台无法将密钥作为 HTTP header 注入，仅支持 body 字段或 query 参数。
// 后端做兼容处理可大幅降低集成门槛。
func APIKeyAuth(expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := extractAPIKey(c)
		if key == "" || key != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":    "UNAUTHORIZED",
				"error_message": "Invalid or missing API Key",
				"detail":        "X-API-Key is required via header (X-API-Key), query (?api_key=), or body field (x_api_key)",
			})
			return
		}
		c.Next()
	}
}

// extractAPIKey 按优先级从 header > query > body 提取 API Key
func extractAPIKey(c *gin.Context) string {
	// 1. Header
	if k := c.GetHeader("X-API-Key"); k != "" {
		return k
	}
	// 兼容大小写变体
	if k := c.GetHeader("x-api-key"); k != "" {
		return k
	}

	// 2. Query string
	if k := c.Query("api_key"); k != "" {
		return k
	}
	if k := c.Query("X-API-Key"); k != "" {
		return k
	}

	// 3. Body（仅对 JSON 请求体做一次轻量解析，不破坏原 body）
	if isJSONRequest(c) {
		if k := extractKeyFromBody(c); k != "" {
			return k
		}
	}
	return ""
}

// isJSONRequest 判断是否为 JSON 请求
func isJSONRequest(c *gin.Context) bool {
	ct := c.GetHeader("Content-Type")
	if len(ct) < 16 {
		return false
	}
	// 不区分大小写，简单截取前 16 字符
	prefix := ct
	if len(prefix) > 16 {
		prefix = prefix[:16]
	}
	return prefix == "application/json" || (len(ct) > 16 && (ct[:16] == "application/jso"))
}

// extractKeyFromBody 读取并恢复 request body，从中提取 x_api_key 字段
func extractKeyFromBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}
	// 读出 body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	// 恢复 body 供后续 handler 继续读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bodyBytes) == 0 {
		return ""
	}

	// 尝试解析为 JSON object
	var payload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return ""
	}

	// 兼容多种字段名
	for _, field := range []string{"x_api_key", "X-API-Key", "api_key", "apikey", "apiKey"} {
		if v, ok := payload[field]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}
