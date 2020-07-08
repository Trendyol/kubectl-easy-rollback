package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-easyrollback",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute(version string) {
	versionCmd.Run = VersionCmdFunc(version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("namespace", "default", "namespace")
	rootCmd.PersistentFlags().String("deployment", "", "deployment")
	rollbackCmd.Flags().String("to-image", "", "to-image")
	rootCmd.AddCommand(listCmd, rollbackCmd, versionCmd)
}
