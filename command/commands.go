package command

import (
	"fmt"
	"github.com/developer-guy/easy-rollback/kubernetes"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Print all the previous deployed image versions of deployment",
	Run:   kubernetes.ListPreviousDeployedImages(),
}

var rollbackCommand = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to given image version",
	Run:   kubernetes.RollbackDeployment(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("namespace", "default", "namespace")
	rootCmd.PersistentFlags().String("deployment", "", "deployment")
	rollbackCommand.Flags().String("to-image", "", "to-image")
	rootCmd.AddCommand(listCommand, rollbackCommand)
}
