package cmd

import (
	"github.com/spf13/cobra"
)

var rollbackCommand = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to given image version",
	Run: func(cmd *cobra.Command, args []string) {
		deploymentFlag := cmd.Flag("deployment").Value.String()
		namespaceFlag := cmd.Flag("namespace").Value.String()
		toImageFlag := cmd.Flag("to-image").Value.String()
		k8sClient.RollbackDeployment(namespaceFlag, deploymentFlag, toImageFlag)
	},
}
