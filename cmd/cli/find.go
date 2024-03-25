package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/output"
	"github.com/stedmanson/grepigee/internal/searcher"

	"github.com/spf13/cobra"
)

var environment string   // To store the environment name
var regExpression string // To store the regex pattern to search for

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Search for regex patterns in Apigee sharedflows and proxies.",
	Long: `Finds and reports occurrences of a specified regex pattern within Apigee sharedflows and proxies. 
	This tool scans through the sharedflow and proxy configurations in your Apigee environment, helping you quickly identify and locate usage of specific patterns. 
	It's particularly useful for auditing, troubleshooting, or ensuring consistency across your API configurations.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Check if the environment flag was set by the user
		if environment == "" {
			fmt.Println("Error: --env flag is required")
			os.Exit(1)
		}

		if regExpression == "" {
			fmt.Println("Error: --expr flag is required")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		foundSharedflowItems := processSharedFlows(environment, regExpression)
		foundProxyItems := processProxies(environment, regExpression)

		combinedItems := append(foundSharedflowItems, foundProxyItems...)

		output.DisplayAsTable(combinedItems)
		output.SaveAsCSV(combinedItems, environment+"-"+regExpression+"-output.csv")

		cleanupDirectory(environment)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().StringVarP(&environment, "env", "e", "", "Specify the environment to search in")
	findCmd.Flags().StringVarP(&regExpression, "expr", "x", "", "Specify the regex pattern to search for")

}

func processSharedFlows(environment string, regExpression string) []searcher.Found {
	sharedflowList, err := apigee.GetSharedFlowList()
	if err != nil {
		fmt.Println("Error getting shared flow list:", err)
		return nil
	}

	deployedSharedflowList := apigee.GetSharedflowDeployments(sharedflowList, environment)

	apigee.DownloadSharedflowRevision(deployedSharedflowList, environment)

	removeZipFiles(environment + "/sharedflows")

	foundSharedflowItems, err := searcher.SearchInDirectory(environment+"/sharedflows", regExpression)
	if err != nil {
		fmt.Println("Error occurred while searching shared flows:", err)
		return nil
	}

	return foundSharedflowItems
}

func processProxies(environment string, regExpression string) []searcher.Found {
	proxyList, err := apigee.GetProxyList()
	if err != nil {
		fmt.Println("Error getting proxy list:", err)
		return nil
	}

	deployedProxyList := apigee.GetProxyDeployments(proxyList, environment)

	apigee.DownloadProxyRevision(deployedProxyList, environment)

	removeZipFiles(environment + "/proxies")

	foundProxyItems, err := searcher.SearchInDirectory(environment+"/proxies", regExpression)
	if err != nil {
		fmt.Println("Error occurred while searching proxies:", err)
		return nil
	}

	return foundProxyItems
}

func cleanupDirectory(directory string) {
	err := os.RemoveAll(directory)
	if err != nil {
		fmt.Printf("Error removing directory %s: %v\n", directory, err)
	}
}

func removeZipFiles(directory string) {
	// Find and remove zip files
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zip") {
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("Error removing file %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory %s: %v\n", directory, err)
	}
}
