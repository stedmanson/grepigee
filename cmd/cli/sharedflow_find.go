package cli

import (
	"fmt"
	"os"

	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/output"
	"github.com/stedmanson/grepigee/internal/searcher"

	"github.com/spf13/cobra"
)

// sharedflowFindCmd represents the find command for proxies
var sharedflowFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Search for regex patterns in Apigee proxies.",
	Long: `Finds and reports occurrences of a specified regex pattern within a sharedflow. 
	This tool scans through the sharedflow configurations in your Apigee environment, helping you quickly identify and locate usage of specific patterns. 
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

		headers, data := output.FormatFoundData(foundSharedflowItems)

		output.DisplayAsTable(headers, data)

		if save {
			output.SaveAsCSV(foundSharedflowItems, "sharedflow-find-"+environment+"-"+regExpression+".csv")
		}

		cleanupDirectory(environment)
	},
}

func init() {
	sharedflowCmd.AddCommand(sharedflowFindCmd)
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
