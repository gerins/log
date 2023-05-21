package echo

import (
	"context"
	"encoding/json"

	"github.com/gerins/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetLogRequest is used for save log request model to echo locals as context.
func SetLogRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get parent context from Echo Locals
			ctx, ok := c.Get("ctx").(context.Context)
			if !ok {
				ctx = context.Background()
			}
			ctx = log.NewRequest().SaveToContext(ctx)
			c.Set("ctx", ctx)
			return next(c)
		}
	}
}

// Save log request to file, dont forget to initiate SetLogRequest middleware before using this middleware.
// Use echo body dump when initiate this middleware.
// echo.Use(middleware.BodyDump(logMiddleware.SaveLogRequest())).
func SaveLogRequest() middleware.BodyDumpHandler {
	return func(c echo.Context, req []byte, resp []byte) {
		// Get parent context from Echo Locals
		ctx, ok := c.Get("ctx").(context.Context)
		if !ok {
			ctx = context.Background()
		}

		extractRequestData(ctx, c, req, resp)
		log.Context(ctx).Save() // Save log request
	}
}

func extractRequestData(ctx context.Context, c echo.Context, req, resp []byte) {
	requestLog := log.Context(ctx) // Get log request from context

	requestLog.IP = c.RealIP()
	requestLog.Method = c.Request().Method
	requestLog.URL = c.Request().Host + c.Request().URL.String()
	requestLog.ReqHeader = getHeader(c, "REQ")
	requestLog.RespHeader = getHeader(c, "RESP")
	requestLog.StatusCode = c.Response().Status
	json.Unmarshal(resp, &requestLog.RespBody) // Get response body

	// Extract Query Args if using GET or DELETE Method
	if requestLog.Method == "GET" || requestLog.Method == "DELETE" {
		queryArgs := make(map[string][]string)
		for k, v := range c.Request().URL.Query() {
			queryArgs[string(k)] = v
		}
		requestLog.ReqBody = queryArgs
	} else {
		json.Unmarshal(req, &requestLog.ReqBody) // Get request body
	}
}

// Get header from request or response
func getHeader(c echo.Context, status string) map[string][]string {
	header := make(map[string][]string)
	if status == "REQ" {
		for k, v := range c.Request().Header {
			header[string(k)] = v
		}
	} else if status == "RESP" {
		for k, v := range c.Response().Header() {
			header[string(k)] = v
		}
	}
	return header
}
