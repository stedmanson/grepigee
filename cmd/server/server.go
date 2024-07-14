package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/deployments"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func main() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Update this to your Redis server address
	})

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // Enable CORS for all routes

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.DEBUG,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Printf("Error: %v\n%s", err, stack)
			return nil
		},
	}))

	e.Logger.SetLevel(log.DEBUG)

	// Routes
	api := e.Group("/api")
	api.GET("/stats", handleAPIStats)
	api.GET("/deployments", handleAPIDeployments)

	//Serve React App
	e.Static("/", "frontend/build")

	//Handle client-side routing
	e.GET("/*", func(c echo.Context) error {
		return c.File("frontend/build/index.html")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func handleAPIStats(c echo.Context) error {
	// Get query parameters
	filterProxy := c.QueryParam("proxyName")
	timeRange := c.QueryParam("timeRange")

	// Use default values if not provided
	if timeRange == "" {
		timeRange = "1h" // Default to 1 hour if not specified
	}

	cacheKey := "stats:" + filterProxy + ":" + timeRange

	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var response map[string]interface{}
		json.Unmarshal([]byte(cachedData), &response)
		return c.JSON(http.StatusOK, response)
	}

	toTime := time.Now().UTC()
	fromTime := calculateFromTime(toTime, timeRange)

	toStr := toTime.Format("01/02/2006 15:04")
	fromStr := fromTime.Format("01/02/2006 15:04")

	headers, data, err := apigee.ListAllTraffic("prod", filterProxy, fromStr, toStr)
	if err != nil {
		c.Logger().Errorf("Error in ListAllTraffic: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Prepare the response
	response := map[string]interface{}{
		"headers": headers,
		"data":    data,
	}

	jsonResponse, _ := json.Marshal(response)
	rdb.Set(ctx, cacheKey, jsonResponse, 5*time.Minute)

	return c.JSON(http.StatusOK, response)
}

func calculateFromTime(toTime time.Time, timeRange string) time.Time {
	switch timeRange {
	case "1h":
		return toTime.Add(-1 * time.Hour)
	case "6h":
		return toTime.Add(-6 * time.Hour)
	case "12h":
		return toTime.Add(-12 * time.Hour)
	case "1d":
		return toTime.AddDate(0, 0, -1)
	case "7d":
		return toTime.AddDate(0, 0, -7)
	default:
		return toTime.Add(-1 * time.Hour) // Default to 1 hour if invalid input
	}
}

func handleAPIDeployments(c echo.Context) error {
	cacheKey := "deployments"

	// Try to get from cache
	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var response map[string]interface{}
		json.Unmarshal([]byte(cachedData), &response)
		return c.JSON(http.StatusOK, response)
	}

	environments, err := apigee.GetEnvironments()
	if err != nil {
		c.Logger().Errorf("Error getting environments: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get environments"})
	}

	allDeployments := deployments.ProcessAllEnvironments(environments)
	headers, data := deployments.FormatDeploymentData(allDeployments, environments)
	response := map[string]interface{}{
		"headers": headers,
		"data":    data,
	}

	jsonResponse, _ := json.Marshal(response)
	rdb.Set(ctx, cacheKey, jsonResponse, 15*time.Minute)

	return c.JSON(http.StatusOK, response)
}
