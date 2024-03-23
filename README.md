# 📜 Log Package 
[![Generic badge](https://img.shields.io/badge/Go-v1.17.0-blue.svg)](https://golang.org/doc/go1.17)
[![Generic badge](https://img.shields.io/badge/status-development-green.svg)](https://shields.io/)
[![Generic badge](https://img.shields.io/badge/release-v1.1.3-yellow.svg)](https://shields.io/)


## 📌 Getting Started
```shell
go get -u github.com/gerins/log
```

### Echo 
```go
package main

import (
	"context"
	"net/http"

	"github.com/gerins/log"
	middlewareLog "github.com/gerins/log/middleware/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Using the default configuration. Please use InitWithConfig() for the production environment.
	log.Init()
	e := echo.New()

	// Initialize logging middleware
	e.Use(middlewareLog.SetLogRequest())                       // Mandatory
	e.Use(middleware.BodyDump(middlewareLog.SaveLogRequest())) // Mandatory

	// Route to simulate logging request
	e.GET("/", func(c echo.Context) error {
		// Get context from echo locals
		ctx := c.Get("ctx").(context.Context)

		// Capture a duration for a function
		defer log.Context(ctx).RecordDuration("handler total process duration").Stop()

		// Assign user id to Log Request
		log.Context(ctx).UserID = 2020

		// Add some extra data
		log.Context(ctx).ExtraData["userData"] = struct {
			Name string
			Age  int
		}{
			Name: "Bob",
			Age:  29,
		}

		// Log Request
		log.Context(ctx).Debug("Testing Log Request Debug")
		log.Context(ctx).Info("Testing Log Request Info")
		log.Context(ctx).Warn("Testing Log Request Warn")
		log.Context(ctx).Error("Testing Log Request Error")

		// Global log
		log.Debug("Testing Global Log Debug")
		log.Info("Testing Global Log Info")
		log.Warn("Testing Global Log Warn")
		log.Error("Testing Global Log Error")

		return c.String(http.StatusOK, "Hello, Log!")
	})

	e.Start("localhost:8080")
}

```

## 🍀 Sample Log Request
```json
{
    "level": "info",
    "time": "2024-03-24T00:12:28.966+0700",
    "caller": "runtime/asm_amd64.s:1650",
    "msg": "REQUEST_LOG",
    "ProcessID": "nRtJ4iAW8YrzcR7ASBM6",
    "UserID": 2020,
    "IP": "127.0.0.1",
    "Method": "GET",
    "URL": "localhost:8080/",
    "RequestHeader": {
        "Accept": [
            "*/*"
        ],
        "User-Agent": [
            "curl/8.1.2"
        ]
    },
    "RequestBody": {},
    "ResponseHeader": {
        "Content-Type": [
            "text/plain; charset=UTF-8"
        ]
    },
    "ResponseBody": null,
    "StatusCode": 200,
    "RequestDuration": 0,
    "ExtraData": {
        "userData": {
            "Name": "Bob",
            "Age": 29
        }
    },
    "SubLog": [
        {
            "level": "[DEBUG] fermi/main.go:43",
            "message": "Testing Log Request Debug"
        },
        {
            "level": "[INFO] fermi/main.go:44",
            "message": "Testing Log Request Info"
        },
        {
            "level": "[WARN] fermi/main.go:45",
            "message": "Testing Log Request Warn"
        },
        {
            "level": "[ERROR] fermi/main.go:46",
            "message": "Testing Log Request Error"
        },
        {
            "level": "[DURATION] fermi/main.go:54",
            "message": "[0.796ms] handler total process duration"
        }
    ]
}
```
