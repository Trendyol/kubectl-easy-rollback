package cmd

import (
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Print all the previous deployed image versions of deployment",
	Run:   k8sClient.ListPreviousDeployedImages(),
}
