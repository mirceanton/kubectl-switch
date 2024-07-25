package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kube-switcher",
	Short: "A tool to switch Kubernetes contexts",
	Long:  `Kubectl-switcher is a CLI tool to switch Kubernetes contexts from multiple kubeconfig files.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add any global flags here
}
