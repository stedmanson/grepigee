package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var environment string
var regExpression string
var save bool

var rootCmd = &cobra.Command{
	Use:   "github.com/stedmanson/grepigee",
	Short: "A handy tool for finding data in Apigee",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&environment, "env", "e", "", "Specify the environment to search in")
	rootCmd.MarkPersistentFlagRequired("env")

	rootCmd.PersistentFlags().StringVarP(&regExpression, "expr", "x", "", "Specify the regex pattern to search for")
	rootCmd.PersistentFlags().BoolVarP(&save, "save", "s", false, "Save output in a csv file")

}
