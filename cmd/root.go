package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
)

type RootOptions struct {
	configFlags *genericclioptions.ConfigFlags

	resultingContext     *api.Context
	resultingContextName string

	userSpecifiedCluster   string
	userSpecifiedContext   string
	userSpecifiedAuthInfo  string
	userSpecifiedNamespace string

	rawConfig      api.Config
	listNamespaces bool
	args           []string

	genericclioptions.IOStreams
}

// NewNamespaceOptions provides an instance of RootOptions with default values
func NewRootOptions(streams genericclioptions.IOStreams) *RootOptions {
	return &RootOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

func NewCmdRoot(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewRootOptions(streams)

	cmd := &cobra.Command{
		Use:   "kubectl-easyrollback",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	o.configFlags.AddFlags(cmd.PersistentFlags())

	cmd.PersistentFlags().String("deployment", "", "deployment")
	rollbackCmd.Flags().String("to-image", "", "to-image")
	cmd.AddCommand(NewListCmd(), rollbackCmd, versionCmd)

	return cmd
}

func Execute(version string) {
	versionCmd.Run = VersionCmdFunc(version)
	rootCmd := NewCmdRoot(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
