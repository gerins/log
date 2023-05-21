package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gerins/log"
	middlewareLog "github.com/gerins/log/middleware/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Init() // Using default configuration
	e := echo.New()

	// Init logging middleware
	e.Use(middlewareLog.SetLogRequest())                       // Mandatory
	e.Use(middleware.BodyDump(middlewareLog.SaveLogRequest())) // Mandatory

	// Init handler
	e.GET("/", func(c echo.Context) error {
		// Get context from echo locals.
		ctx := c.Get("ctx").(context.Context)

		// Assign user id to Log Request model
		// So wen can know who make the request to the server
		log.Context(ctx).UserID = 2020

		// Log Request, support multiple arguments
		log.Context(ctx).Debug(1, "Testing Log Request Debug")
		log.Context(ctx).Info(2, "Testing Log Request Info")
		log.Context(ctx).Warn(3, "Testing Log Request Warn")
		log.Context(ctx).Error(4, "Testing Log Request Error")

		func() {
			defer log.Context(ctx).RecordDuration("Try to sleep 50 Milliseconds").Stop()
			time.Sleep(50 * time.Millisecond)
		}()

		// Global log
		log.Debug("Testing Global Log Debug")
		log.Info("Testing Global Log Info")
		log.Warn("Testing Global Log Warn")
		log.Error("Testing Global Log Error")

		return c.String(http.StatusOK, "Hello, Log!")
	})

	e.Start("localhost:8080")
}
