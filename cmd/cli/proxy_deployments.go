package cmd

import (
	"fmt"
	"os"

	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/output"
	"github.com/stedmanson/grepigee/internal/searcher"

	"github.com/spf13/cobra"
)

// proxyFindCmd represents the find command for proxies
var proxyFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Search for regex patterns in Apigee proxies.",
	Long: `Finds and reports occurrences of a specified regex pattern within a proxy. 
	This tool scans through the proxy configurations in your Apigee environment, helping you quickly identify and locate usage of specific patterns. 
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
		foundProxyItems := processProxies(environment, regExpression)
		headers, data := output.FormatFoundData(foundProxyItems)

		output.DisplayAsTable(headers, data)

		if save {
			output.SaveAsCSV(foundProxyItems, "proxy-find-"+environment+"-"+regExpression+".csv")
		}

		cleanupDirectory(environment)
	},
}

func init() {
	proxyCmd.AddCommand(proxyFindCmd)
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

	cleanupDirectory(environment)

	return foundProxyItems
}
