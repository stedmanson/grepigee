package cmd

import (
	"github.com/spf13/cobra"
)

// sharedflowCmd represents the command for sharedflow
var sharedflowCmd = &cobra.Command{
	Use:   "sharedflow",
	Short: "Functions for finding information within sharedflows.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(sharedflowCmd)
}
