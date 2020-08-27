package cmd

import (
	"github.com/Trendyol/kubectl-easy-rollback/client"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Print all the previous deployed image versions of deployment",
		Run: func(cmd *cobra.Command, args []string) {
			k8sClient := client.NewK8sClient(cmd.Flag("kubeconfig").Value.String(), cmd.Flag("context").Value.String())
			deploymentFlag := cmd.Flag("deployment").Value.String()
			namespaceFlag := cmd.Flag("namespace").Value.String()
			k8sClient.ListPreviousDeployedImages(deploymentFlag, namespaceFlag)
		},
	}
	return cmd
}
