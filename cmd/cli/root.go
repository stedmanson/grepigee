package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "grepigee",
	Short: "A handy tool for finding data in Apigee",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if environment == "" {
			environment = viper.GetString("environment")
		}
		if organisation == "" {
			organisation = viper.GetString("organisation")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&environment, "env", "e", "", "Specify the environment to search in")
	rootCmd.PersistentFlags().StringVarP(&organisation, "org", "o", "", "Specify the organisation")

	rootCmd.PersistentFlags().StringVarP(&regExpression, "expr", "x", "", "Specify the regex pattern to search for")
	rootCmd.PersistentFlags().BoolVarP(&save, "save", "s", false, "Save output in a csv file")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "save-config",
		Short: "Save the current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := saveConfig(); err != nil {
				fmt.Println("Error saving config:", err)
				os.Exit(1)
			}
			fmt.Println("Configuration saved successfully")
		},
	})

}
