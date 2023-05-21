package fiber

import (
	"encoding/json"

	"github.com/gerins/log"
	"github.com/gofiber/fiber/v2"
)

func SaveLogRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user context from Fiber Locals
		ctx := c.UserContext()
		requestLog := log.NewRequest()
		c.SetUserContext(requestLog.SaveToContext(ctx))

		// Proceed the request first, we capture the data later
		if err := c.Next(); err != nil {
			requestLog.Errorf("saveLogRequest next Middleware failed, %s, %s", err.Error(), c.OriginalURL())
		}

		requestLog.IP = getClientIPAdress(c)
		requestLog.Method = string(c.Request().Header.Method())
		requestLog.URL = c.Hostname() + string(c.OriginalURL())
		requestLog.ReqHeader = getHeader(c, "REQ")
		requestLog.RespHeader = getHeader(c, "RESP")
		requestLog.StatusCode = c.Response().StatusCode()
		json.Unmarshal(c.Response().Body(), &requestLog.RespBody) // Get Response body

		// Extract Query Args if using GET or DELETE Method
		if requestLog.Method == fiber.MethodGet || requestLog.Method == fiber.MethodDelete {
			queryArgs := make(map[string]string)
			c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
				queryArgs[string(key)] = string(value)
			})
			requestLog.ReqBody = queryArgs
		} else {
			json.Unmarshal(c.Request().Body(), &requestLog.ReqBody) // Get Request body
		}

		requestLog.Save()
		return nil
	}
}

// Get header from request or response
func getHeader(f *fiber.Ctx, status string) map[string]string {
	header := make(map[string]string)
	if status == "REQ" {
		f.Request().Header.VisitAll(func(key, value []byte) {
			header[string(key)] = string(value)
		})
	} else if status == "RESP" {
		f.Response().Header.VisitAll(func(key, value []byte) {
			header[string(key)] = string(value)
		})
	}
	return header
}

// getClientIPAdress is for get User IP Address
func getClientIPAdress(f *fiber.Ctx) string {
	if f.Get("X-Real-Ip") == "" {
		return f.IP()
	}
	return f.Get("X-Real-Ip")
}
