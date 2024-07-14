package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/cache"
	"github.com/stedmanson/grepigee/internal/utils"
)

func HandleAPIStats(c echo.Context) error {
	filterProxy := c.QueryParam("proxyName")
	timeRange := c.QueryParam("timeRange")

	if timeRange == "" {
		timeRange = "1h"
	}

	cacheKey := "stats:" + filterProxy + ":" + timeRange

	cachedData, err := cache.Get(cacheKey)
	if err == nil {
		var response map[string]interface{}
		json.Unmarshal([]byte(cachedData), &response)
		return c.JSON(http.StatusOK, response)
	}

	toTime := time.Now().UTC()
	fromTime := utils.CalculateFromTime(toTime, timeRange)

	toStr := toTime.Format("01/02/2006 15:04")
	fromStr := fromTime.Format("01/02/2006 15:04")

	headers, data, err := apigee.ListAllTraffic("prod", filterProxy, fromStr, toStr)
	if err != nil {
		c.Logger().Errorf("Error in ListAllTraffic: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"headers": headers,
		"data":    data,
	}

	jsonResponse, _ := json.Marshal(response)
	cache.Set(cacheKey, jsonResponse, 5*time.Minute)

	return c.JSON(http.StatusOK, response)
}
