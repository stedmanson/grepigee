package cmd

import (
	"github.com/spf13/cobra"
)

// proxyCmd represents the grep command for proxies
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Functions for finding information within proxies.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
