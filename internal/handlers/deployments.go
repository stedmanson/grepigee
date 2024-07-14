package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/cache"
	"github.com/stedmanson/grepigee/internal/deployments"
)

func HandleAPIDeployments(c echo.Context) error {
	cacheKey := "deployments"

	cachedData, err := cache.Get(cacheKey)
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
	cache.Set(cacheKey, jsonResponse, 15*time.Minute)

	return c.JSON(http.StatusOK, response)
}
