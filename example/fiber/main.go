package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/gerins/log"
	middlewareLog "github.com/gerins/log/middleware/fiber"
)

func main() {
	log.Init() // Using default configuration
	f := fiber.New()

	// Init logging middleware
	f.Use(middlewareLog.SaveLogRequest()) // Mandatory

	f.Get("", func(c *fiber.Ctx) error {
		// Get user context from fiber.
		ctx := c.UserContext()

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

		return c.Status(http.StatusOK).JSON("Hello, Log!")
	})

	f.Listen("localhost:8080")
}
