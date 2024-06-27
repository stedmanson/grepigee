package cmd

import (
	"github.com/spf13/cobra"
)

// statsCmd represents the grep command for proxies
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Functions for displaying statistics on Apigee performance.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
