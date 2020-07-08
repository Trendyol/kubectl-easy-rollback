package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kubectl-easyrollback",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kubectl-easyrollback - version: %s", version)
	},
}
