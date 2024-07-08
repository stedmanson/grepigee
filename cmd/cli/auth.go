package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication related commands",
	Long:  `Commands for authentication and token management.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login and get APIGEE_BEARER_TOKEN",
	Long:  `Login to Apigee and get the APIGEE_BEARER_TOKEN. Run this command with 'eval $(grepagee auth login)' to set the environment variable.`,
	Run:   runLogin,
}

func runLogin(cmd *cobra.Command, args []string) {
	// Check if get_token exists
	_, err := exec.LookPath("get_token")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: 'get_token' command not found.")
		fmt.Fprintln(os.Stderr, "Please install the 'get_token' command. You can find instructions at: [URL to instructions]")
		os.Exit(1)
	}

	if runtime.GOOS == "windows" {
		fmt.Println("$env:APIGEE_BEARER_TOKEN = (get_token); setx APIGEE_BEARER_TOKEN $env:APIGEE_BEARER_TOKEN")
		fmt.Fprintln(os.Stderr, "Run the above command in PowerShell to set the APIGEE_BEARER_TOKEN.")
	} else {
		fmt.Println("export APIGEE_BEARER_TOKEN=$(get_token)")
		fmt.Fprintln(os.Stderr, "Run the above command to set the APIGEE_BEARER_TOKEN.")
	}
}
