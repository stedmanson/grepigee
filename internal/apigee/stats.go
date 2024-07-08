package apigee

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func ListAllTraffic(environment string, filterProxy string, from string, to string) ([]string, [][]string, error) {
	var dimension = "apiproxy"
	if filterProxy != "" {
		dimension = "apiproxy,request_path"
	}

	var trafficList = new(StatsTrafficList)

	params := url.Values{}
	params.Add("select", "sum(message_count),max(request_processing_latency),max(target_response_time),max(response_processing_latency),max(total_response_time)")
	params.Add("timeRange", from+"~"+to)
	params.Add("sortby", "sum(message_count)")
	if filterProxy != "" {
		params.Add("filter", "(apiproxy in '"+filterProxy+"')")
	}

	url := baseURL + "organizations/woolworths/environments/" + environment + "/stats/" + dimension + "?" + params.Encode()

	body, err := Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("error calling %s: %s", url, err)
	}

	err = json.Unmarshal(body, trafficList)
	if err != nil {
		return nil, nil, err
	}

	return processListAllTraffic(trafficList, environment)
}

func processListAllTraffic(list *StatsTrafficList, environment string) ([]string, [][]string, error) {
	var headers = []string{"Proxy Name", "Traffic Count", "request_processing_latency", "target_response_time", "response_processing_latency", "total_response_time"}
	var data = [][]string{}

	if len(list.Environments) == 0 {
		return headers, data, fmt.Errorf("no data found for any environment")
	}

	var found bool
	for _, env := range list.Environments {
		if env.Name == environment {
			found = true

			for _, proxy := range env.Dimensions {
				var count string
				var reqestLatency string
				var targetResponseTime string
				var responseLatency string
				var totalResponseTime string

				count, _ = humanReadable(proxy.Metrics[0].Values[0], false)
				reqestLatency, _ = humanReadable(proxy.Metrics[1].Values[0], true)
				targetResponseTime, _ = humanReadable(proxy.Metrics[2].Values[0], true)
				responseLatency, _ = humanReadable(proxy.Metrics[3].Values[0], true)
				totalResponseTime, _ = humanReadable(proxy.Metrics[4].Values[0], true)

				data = append(data, []string{proxy.Name, count, reqestLatency, targetResponseTime, responseLatency, totalResponseTime})
			}
		}
	}

	if !found {
		return headers, data, fmt.Errorf("no data found for environment: %s", environment)
	}

	return headers, data, nil
}
