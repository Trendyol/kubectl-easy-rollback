package cmd

import (
	"fmt"
	"github.com/Trendyol/easy-rollback/client"
	"github.com/spf13/cobra"
	"os"
)

var k8sClient *client.K8SClient

var rootCmd = &cobra.Command{
	Use:   "kubectl-easyrollback",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func initConfig() {
	fmt.Println("initializing..")
	k8sClient = client.NewK8sClient()
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
	rootCmd.AddCommand(listCommand, rollbackCommand, versionCommand)
}
