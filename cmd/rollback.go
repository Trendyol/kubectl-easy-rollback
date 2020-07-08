package cmd

import (
	"github.com/spf13/cobra"
)

var rollbackCommand = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to given image version",
	Run:   k8sClient.RollbackDeployment(),
}
