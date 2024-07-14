package shared_ops

import (
	"encoding/json"
	"time"

	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/cache"
)

type StatsRequest struct {
	Environment string
	FilterProxy string
	FromTime    time.Time
	ToTime      time.Time
}

func GetTrafficStats(req StatsRequest, useCache bool) (map[string]interface{}, error) {
	var cacheKey string
	if useCache {
		cacheKey := "stats:" + req.FilterProxy + ":" + req.FromTime.Format("01/02/2006 15:04") + ":" + req.ToTime.Format("01/02/2006 15:04")
		cachedData, err := cache.Get(cacheKey)
		if err == nil {
			var response map[string]interface{}
			json.Unmarshal([]byte(cachedData), &response)
			return response, nil
		}
	}

	fromStr := req.FromTime.Format("01/02/2006 15:04")
	toStr := req.ToTime.Format("01/02/2006 15:04")

	headers, data, err := apigee.ListAllTraffic(req.Environment, req.FilterProxy, fromStr, toStr)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"headers": headers,
		"data":    data,
	}

	if useCache {
		jsonResponse, _ := json.Marshal(response)
		cache.Set(cacheKey, jsonResponse, 5*time.Minute)
	}

	return response, nil
}

func CalculateFromTime(toTime time.Time, timeRange string) time.Time {
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
		return toTime.Add(-1 * time.Hour)
	}
}
