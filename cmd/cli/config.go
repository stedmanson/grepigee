package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	cfgFile       string
	environment   string
	organisation  string
	regExpression string
	save          bool
)

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".grepagee" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".grepigee")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Bind the flags to viper
	viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("env"))
	// Bind other flags here

	// Set default values
	viper.SetDefault("environment", "")
	// Set other default values here
}

func saveConfig() error {
	viper.Set("environment", environment)
	// Set other variables here

	configHome, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configHome, ".grepigee.yaml")
	return viper.WriteConfigAs(configPath)
}
