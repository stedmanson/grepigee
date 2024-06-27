package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/output"

	"github.com/spf13/cobra"
)

// statsListTrafficCmd represents the find command for proxies
var statsListTrafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Display traffic data information for proxies or developer 	apps.",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Check if the environment flag was set by the user
		if environment == "" {
			fmt.Println("Error: --env flag is required")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		filter, _ := cmd.Flags().GetString("filter")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		parsedFromDt := parseDateTime(from)
		parsedToDt := parseDateTime(to)

		headers, data, err := apigee.ListAllTraffic(environment, filter, parsedFromDt, parsedToDt)
		if err != nil {
			fmt.Println(err)
			return
		}

		output.DisplayAsTable(headers, data)
	},
}

func init() {
	currentTime := time.Now()

	statsListTrafficCmd.Flags().StringP("filter", "f", "", "filter the list returned for traffic metrics")
	statsListTrafficCmd.Flags().String("from", currentTime.AddDate(0, 0, -1).Format("01/02/2006 15:04"), "From datetime (MM/DD/YYYY HH:MM)")
	statsListTrafficCmd.Flags().String("to", currentTime.Format("01/02/2006 15:04"), "To datetime (MM/DD/YYYY HH:MM)")

	statsCmd.AddCommand(statsListTrafficCmd)
}
