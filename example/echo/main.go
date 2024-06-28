package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gerins/log"
	middlewareLog "github.com/gerins/log/middleware/echo"
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
		time.Sleep(100 * time.Millisecond) // Simulate a process

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
