package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stedmanson/grepigee/internal/apigee"
	"github.com/stedmanson/grepigee/internal/output"
	"github.com/stedmanson/grepigee/internal/shared_ops"
)

var proxyDeploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "List all deployments in Apigee proxies across environments.",
	Long:  `Lists all proxy deployments across specified Apigee environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		environments, _ := apigee.GetEnvironments()
		allDeployments := shared_ops.ProcessAllEnvironments(environments)

		headers, data := shared_ops.FormatDeploymentData(allDeployments, environments)
		output.DisplayAsTable(headers, data)
	},
}

func init() {
	proxyCmd.AddCommand(proxyDeploymentsCmd)
}
