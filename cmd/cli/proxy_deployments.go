package cmd

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/output"
)

type ProxyDeploymentInfo struct {
	Name        string
	Version     string
	State       string
	Environment string
}

var proxyDeploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "List all deployments in Apigee proxies across environments.",
	Long:  `Lists all proxy deployments across specified Apigee environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		environments, _ := apigee.GetEnvironments()
		allDeployments := processAllEnvironments(environments)

		headers, data := formatDeploymentData(allDeployments, environments)
		output.DisplayAsTable(headers, data)
	},
}

func init() {
	proxyCmd.AddCommand(proxyDeploymentsCmd)
}

func processAllEnvironments(environments []string) []ProxyDeploymentInfo {
	var allDeployments []ProxyDeploymentInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, env := range environments {
		wg.Add(1)
		go func(environment string) {
			defer wg.Done()
			deployments := processProxyDeployments(environment)

			mu.Lock()
			allDeployments = append(allDeployments, deployments...)
			mu.Unlock()
		}(env)
	}

	// Use a timeout for the entire process
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:

	case <-time.After(10 * time.Minute): // Adjust timeout as needed
		fmt.Println("Warning: Timeout while processing environments")
	}

	return allDeployments
}

func processProxyDeployments(environment string) []ProxyDeploymentInfo {
	proxyList, err := apigee.GetProxyList()
	if err != nil {
		fmt.Printf("Error getting proxy list: %v\n", err)
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

	// Use a timeout to prevent indefinite hanging
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All goroutines finished successfully
	case <-time.After(2 * time.Minute): // Adjust timeout as needed
		fmt.Printf("Warning: Timeout while processing environment %s\n", environment)
	}

	return deployments
}

func formatDeploymentData(deployments []ProxyDeploymentInfo, environments []string) ([]string, [][]string) {
	headers := append([]string{"Name"}, environments...)

	// Create a map to hold deployment info for each proxy
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

	// Sort proxy names for consistent output
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
