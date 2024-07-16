package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stedmanson/grepigee/internal/output"
	"github.com/stedmanson/grepigee/internal/shared_ops"
)

var statsListTrafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Display traffic data information for proxies or developer apps.",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if environment == "" {
			environment = viper.GetString("environment")
		}
		if environment == "" {
			fmt.Println("Error: environment is not set. Use --env flag or set it in the config file.")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		filter, _ := cmd.Flags().GetString("filter")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		fromTime, _ := time.Parse("01/02/2006 15:04", from)
		toTime, _ := time.Parse("01/02/2006 15:04", to)

		req := shared_ops.StatsRequest{
			Environment: environment,
			FilterProxy: filter,
			FromTime:    fromTime,
			ToTime:      toTime,
		}

		response, err := shared_ops.GetTrafficStats(req, false)
		if err != nil {
			fmt.Println(err)
			return
		}

		output.DisplayAsTable(response["headers"].([]string), response["data"].([][]string))
	},
}

func init() {
	currentTime := time.Now()

	statsListTrafficCmd.Flags().StringP("filter", "f", "", "filter the list returned for traffic metrics")
	statsListTrafficCmd.Flags().String("from", currentTime.AddDate(0, 0, -1).Format("01/02/2006 15:04"), "From datetime (MM/DD/YYYY HH:MM)")
	statsListTrafficCmd.Flags().String("to", currentTime.Format("01/02/2006 15:04"), "To datetime (MM/DD/YYYY HH:MM)")

	statsCmd.AddCommand(statsListTrafficCmd)
}
