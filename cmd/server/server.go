package main

import (
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stedmanson/grepigee/internal/apigee"
)

type ProxyStats struct {
	ProxyName                 string `json:"proxyName"`
	TrafficCount              string `json:"trafficCount"`
	RequestProcessingLatency  string `json:"requestProcessingLatency"`
	TargetResponseTime        string `json:"targetResponseTime"`
	ResponseProcessingLatency string `json:"responseProcessingLatency"`
	TotalResponseTime         string `json:"totalResponseTime"`
}

type PageData struct {
	Title string
	Stats []ProxyStats
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.DEBUG,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Printf("Error: %v\n%s", err, stack)
			return nil
		},
	}))

	e.Logger.SetLevel(log.DEBUG)

	// Initialize templates
	templatesDir := filepath.Join("web", "templates") // Adjust this path as needed
	t := &Template{
		templates: template.Must(template.ParseGlob(filepath.Join(templatesDir, "*.html"))),
	}
	e.Renderer = t

	// Routes
	e.GET("/", handleHome)
	e.GET("/stats", handleStats)
	e.GET("/grep", handleGrep)
	e.GET("/api/stats", handleAPIStats)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func handleHome(c echo.Context) error {
	return c.Render(http.StatusOK, "layout", PageData{Title: "Home"})
}

func handleStats(c echo.Context) error {
	return c.Render(http.StatusOK, "layout", PageData{Title: "Stats"})
}

func handleGrep(c echo.Context) error {
	return c.Render(http.StatusOK, "layout", PageData{Title: "Grep"})
}

func handleAPIStats(c echo.Context) error {
	// Get query parameters
	filterProxy := c.QueryParam("proxyName")
	timeRange := c.QueryParam("timeRange")

	// Use default values if not provided
	if timeRange == "" {
		timeRange = "1h" // Default to 1 hour if not specified
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
