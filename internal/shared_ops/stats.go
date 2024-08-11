package shared_ops

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
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
type ProxyStatus struct {
	Name     string
	Deployed bool
}

func GetTrafficStats(req StatsRequest, useCache bool) (map[string]interface{}, error) {
	var cacheKey string
	var allProxies []ProxyStatus
	var err error

	if useCache {
		// Round times to the nearest hour for more effective caching
		roundedFromTime := req.FromTime.Truncate(time.Hour)
		roundedToTime := req.ToTime.Truncate(time.Hour)

		// Calculate the duration and round it to common intervals
		duration := roundedToTime.Sub(roundedFromTime)
		roundedDuration := roundDuration(duration)

		cacheKey = fmt.Sprintf("stats:%s:%s:%s",
			req.FilterProxy,
			roundedFromTime.Format("2006-01-02T15:00:00Z"),
			roundedDuration.String())

		cachedData, err := cache.Get(cacheKey)
		if err == nil {
			var response map[string]interface{}
			json.Unmarshal([]byte(cachedData), &response)
			return response, nil
		}
	}

	fromStr := req.FromTime.Format("01/02/2006 15:04")
	toStr := req.ToTime.Format("01/02/2006 15:04")

	if req.FilterProxy == "" {
		// Get all proxies and their deployment status
		allProxies, err = GetAllProxiesStatus(req.Environment)
		if err != nil {
			return nil, fmt.Errorf("error getting proxy list: %v", err)
		}
	}

	headers, data, err := apigee.ListAllTraffic(req.Environment, req.FilterProxy, fromStr, toStr)
	if err != nil {
		return nil, err
	}

	// Create a map to easily look up traffic data
	trafficMap := make(map[string][]string)
	for _, row := range data {
		trafficMap[row[0]] = row[1:] // row[0] is the proxy name
	}

	// Prepare the final data including all proxies
	finalData := make([][]string, 0, len(allProxies))
	for _, proxy := range allProxies {
		var row []string
		if trafficData, exists := trafficMap[proxy.Name]; exists {
			row = append([]string{proxy.Name}, trafficData...)
		} else if proxy.Deployed {
			row = []string{proxy.Name, "0", "-", "-", "-", "-"}
		} else {
			row = []string{proxy.Name, "-", "-", "-", "-", "-"}
		}
		finalData = append(finalData, row)
	}

	// Sort the finalData
	sort.Slice(finalData, func(i, j int) bool {
		countI := parseTrafficCount(finalData[i][1])
		countJ := parseTrafficCount(finalData[j][1])

		if countI != countJ {
			return countI > countJ // Sort by traffic count descending
		}

		// If traffic counts are equal, sort by proxy name
		return finalData[i][0] < finalData[j][0]
	})

	response := map[string]interface{}{
		"headers": headers,
		"data":    finalData,
	}

	if useCache {
		jsonResponse, _ := json.Marshal(response)
		cache.Set(cacheKey, jsonResponse, 5*time.Minute)
	}

	return response, nil
}

// GetAllProxiesStatus returns a list of all proxies and their deployment status
func GetAllProxiesStatus(environment string) ([]ProxyStatus, error) {
	proxyList, err := apigee.GetProxyList()
	if err != nil {
		return nil, fmt.Errorf("error getting proxy list: %v", err)
	}

	deployedChan, undeployedChan := apigee.StreamProxyDeployments(proxyList, environment)

	proxyStatus := make([]ProxyStatus, 0, len(proxyList))
	deployedMap := make(map[string]bool)

	// Use a WaitGroup to ensure we process all items from both channels
	var wg sync.WaitGroup
	wg.Add(2)

	// Process deployed proxies
	go func() {
		defer wg.Done()
		for deployment := range deployedChan {
			deployedMap[deployment.Name] = true
		}
	}()

	// Process undeployed proxies
	go func() {
		defer wg.Done()
		for undeployed := range undeployedChan {
			deployedMap[undeployed.Name] = false
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Create final list of proxy statuses
	for _, proxyName := range proxyList {
		deployed, exists := deployedMap[proxyName]
		if !exists {
			deployed = false // Consider as undeployed if not found in the deployment streams
		}
		proxyStatus = append(proxyStatus, ProxyStatus{Name: proxyName, Deployed: deployed})
	}

	return proxyStatus, nil
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
	case "14d":
		return toTime.AddDate(0, 0, -14)
	case "30d":
		return toTime.AddDate(0, 0, -30)
	default:
		return toTime.Add(-1 * time.Hour)
	}
}

// parseTrafficCount converts the traffic count string to a number for sorting
func parseTrafficCount(count string) int64 {
	if count == "-" {
		return -1 // Undeployed proxies will be at the end
	}

	// Remove any commas from the count string
	count = strings.ReplaceAll(count, ",", "")

	n, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return 0 // If parsing fails, treat as zero
	}
	return n
}

// roundDuration rounds the given duration to common intervals
func roundDuration(d time.Duration) time.Duration {
	switch {
	case d <= time.Hour:
		return time.Hour
	case d <= 3*time.Hour:
		return 3 * time.Hour
	case d <= 6*time.Hour:
		return 6 * time.Hour
	case d <= 12*time.Hour:
		return 12 * time.Hour
	case d <= 24*time.Hour:
		return 24 * time.Hour
	default:
		return d.Round(24 * time.Hour)
	}
}
