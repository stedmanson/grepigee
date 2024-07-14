package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stedmanson/grepigee/internal/cache"
	"github.com/stedmanson/grepigee/internal/handlers"
)

func main() {
	// Initialize Redis client
	cache.InitRedis("localhost:6379")

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.DEBUG,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Printf("Error: %v\n%s", err, stack)
			return nil
		},
	}))

	e.Logger.SetLevel(log.DEBUG)

	// Routes
	api := e.Group("/api")
	api.GET("/stats", handlers.HandleAPIStats)
	api.GET("/deployments", handlers.HandleAPIDeployments)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
