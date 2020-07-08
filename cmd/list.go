package cmd

import (
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Print all the previous deployed image versions of deployment",
	Run: func(cmd *cobra.Command, args []string) {
		deploymentFlag := cmd.Flag("deployment").Value.String()
		namespaceFlag := cmd.Flag("namespace").Value.String()
		k8sClient.ListPreviousDeployedImages(deploymentFlag, namespaceFlag)
	},
}
