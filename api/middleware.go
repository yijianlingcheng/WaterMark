package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"WaterMark/internal"
)

// 自定义一个结构体,实现 gin.ResponseWriter interface.
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 获取中间件列表.
func getMiddlewareList() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		logMiddleware(),
		corsMiddleware(),
	}
}

// corsMiddleware 处理跨域问题.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)

			return
		}
		c.Next()
	}
}

// 重写 Write([]byte) (int, error) 方法.
func (w responseWriter) Write(b []byte) (int, error) {
	// 向一个bytes.buffer中写一份数据来为获取body使用
	n, err := w.b.Write(b)
	if err != nil {
		return n, err
	}

	// 完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

// logMiddleware 日志中间件,记录接口请求与响应的参数.
func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 请求日期
		requestDate := start.Format(time.RFC3339)
		// 请求接口路由
		requestPath := c.Request.URL
		// 这里面的接口不记录日志
		for _, noLogApi := range getNoLogApis() {
			if strings.Contains(requestPath.String(), noLogApi) {
				c.Next()

				return
			}
		}
		// 请求方式
		requestMethod := c.Request.Method
		// 请求体 body
		var requestBody string
		b, err := c.GetRawData()
		if err != nil {
			requestBody = "failed to get request body"
		} else {
			requestBody = string(b)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer(make([]byte, 0)),
		}
		c.Writer = writer
		c.Next()
		// 请求耗时
		cost := time.Since(start).Milliseconds()
		// 响应状态码
		responseStatus := c.Writer.Status()
		// 响应体 body
		responseBody := writer.b.String()

		// 记录日志
		internal.Log.Debug(fmt.Sprintf("[GIN] %s |%d| %d | %s | %s | %s | %s",
			requestDate,
			responseStatus,
			cost,
			requestMethod,
			requestPath,
			requestBody,
			responseBody))
	}
}
