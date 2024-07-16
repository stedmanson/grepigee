package shared_ops

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/stedmanson/grepigee/internal/apigee"
)

type ProxyDeploymentInfo struct {
	Name        string
	Version     string
	State       string
	Environment string
}

func ProcessAllEnvironments(environments []string) []ProxyDeploymentInfo {
	var allDeployments []ProxyDeploymentInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, env := range environments {
		wg.Add(1)
		go func(environment string) {
			defer wg.Done()
			deployments := ProcessProxyDeployments(environment)

			mu.Lock()
			allDeployments = append(allDeployments, deployments...)
			mu.Unlock()
		}(env)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All environments processed successfully
	case <-time.After(10 * time.Minute): // Adjust timeout as needed
		// Log timeout warning
	}

	return allDeployments
}

func ProcessProxyDeployments(environment string) []ProxyDeploymentInfo {
	proxyList, err := apigee.GetProxyList()
	if err != nil {
		// Log error
		return nil
	}

	deployedProxyList, undeployedEntities := apigee.StreamProxyDeployments(proxyList, environment)

	var deployments []ProxyDeploymentInfo
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for dep := range deployedProxyList {
			for _, rev := range dep.Revision {
				if rev.State == "deployed" {
					mu.Lock()
					deployments = append(deployments, ProxyDeploymentInfo{
						Name:        dep.Name,
						Version:     rev.Name,
						State:       "deployed",
						Environment: environment,
					})
					mu.Unlock()
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for dep := range undeployedEntities {
			mu.Lock()
			deployments = append(deployments, ProxyDeploymentInfo{
				Name:        dep.Name,
				Version:     "",
				State:       "undeployed",
				Environment: environment,
			})
			mu.Unlock()
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines finished successfully
	case <-time.After(2 * time.Minute): // Adjust timeout as needed
		// Log timeout warning
	}

	return deployments
}

func FormatDeploymentData(deployments []ProxyDeploymentInfo, environments []string) ([]string, [][]string) {
	headers := append([]string{"Name"}, environments...)

	proxyMap := make(map[string]map[string][]string)

	for _, dep := range deployments {
		if _, exists := proxyMap[dep.Name]; !exists {
			proxyMap[dep.Name] = make(map[string][]string)
		}
		if dep.State == "deployed" {
			proxyMap[dep.Name][dep.Environment] = append(proxyMap[dep.Name][dep.Environment], dep.Version)
		} else {
			proxyMap[dep.Name][dep.Environment] = append(proxyMap[dep.Name][dep.Environment], "undeployed")
		}
	}

	var data [][]string
	var proxyNames []string
	for name := range proxyMap {
		proxyNames = append(proxyNames, name)
	}
	sort.Strings(proxyNames)

	for _, name := range proxyNames {
		row := []string{name}
		for _, env := range environments {
			versions, exists := proxyMap[name][env]
			if !exists {
				row = append(row, "-")
			} else {
				row = append(row, strings.Join(versions, ","))
			}
		}
		data = append(data, row)
	}

	return headers, data
}
