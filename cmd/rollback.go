package cmd

import (
	"github.com/Trendyol/easy-rollback/client"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to given image version",
	Run: func(cmd *cobra.Command, args []string) {
		k8sClient := client.NewK8sClient(cmd.Flag("kubeconfig").Value.String(), cmd.Flag("context").Value.String())
		deploymentFlag := cmd.Flag("deployment").Value.String()
		namespaceFlag := cmd.Flag("namespace").Value.String()
		toImageFlag := cmd.Flag("to-image").Value.String()
		k8sClient.RollbackDeployment(namespaceFlag, deploymentFlag, toImageFlag)
	},
}