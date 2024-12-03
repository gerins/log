package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/gerins/log"
)

// ResponseBodyWriter wraps Gin's ResponseWriter to capture the response body.
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write writes data to the buffer and the original ResponseWriter.
func (w *ResponseBodyWriter) Write(data []byte) (int, error) {
	w.body.Write(data)                  // Write to buffer
	return w.ResponseWriter.Write(data) // Write to original ResponseWriter
}

// SetLogRequest sets up the logging request in the Gin context.
func SetLogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the existing context from the request
		ctx := c.Request.Context()

		// Add the log request to the context
		newCtxWithLog := log.NewRequest().SaveToContext(ctx)

		// Replace the request with the new context
		c.Request = c.Request.WithContext(newCtxWithLog)

		// Proceed to the next handler
		c.Next()
	}
}

// SaveLogRequest handles logging of request and response data.
func SaveLogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get context
		ctx := c.Request.Context()

		// Capture request body
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(io.NopCloser(bytes.NewBuffer(reqBody))) // Re-read buffer for next handler

		// Create a ResponseBodyWriter
		bodyWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}

		// Replace Gin's ResponseWriter with our custom writer
		c.Writer = bodyWriter

		// Process request in the main handler
		c.Next()

		// Save log request
		extractRequestData(ctx, c, reqBody, bodyWriter.body.Bytes())
		log.Context(ctx).Save()
	}
}

func extractRequestData(ctx context.Context, c *gin.Context, req, resp []byte) {
	requestLog := log.Context(ctx) // Get log request from context

	// Populate log information
	requestLog.IP = c.ClientIP()
	requestLog.Method = c.Request.Method
	requestLog.URL = c.Request.Host + c.Request.URL.String()
	requestLog.ReqHeader, requestLog.RespHeader = getHeader(c)
	requestLog.StatusCode = c.Writer.Status()

	// Set request body based on HTTP method
	if requestLog.Method == http.MethodGet || requestLog.Method == http.MethodDelete {
		requestLog.ReqBody = c.Request.URL.Query()
	} else if requestLog.ReqBody == nil {
		if err := json.Unmarshal(req, &requestLog.ReqBody); err != nil {
			requestLog.ReqBody = string(req)
		}
	}

	// Set response body
	if requestLog.RespBody == nil {
		// Only capture file names if the response body is a file
		if headers, ok := c.Writer.Header()["Content-Disposition"]; ok {
			for _, value := range headers {
				if strings.Contains(value, "attachment") {
					requestLog.RespBody = value
					return
				}
			}
		}

		if err := json.Unmarshal(resp, &requestLog.RespBody); err != nil {
			requestLog.RespBody = string(resp)
		}
	}
}

// getHeader extracts headers from the request or response.
func getHeader(c *gin.Context) (map[string][]string, map[string][]string) {
	var (
		req  = make(map[string][]string)
		resp = make(map[string][]string)
	)

	for k, v := range c.Request.Header {
		req[k] = v
	}
	for k, v := range c.Writer.Header() {
		resp[k] = v
	}

	return req, resp
}
