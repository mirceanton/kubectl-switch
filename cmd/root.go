package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Global flag for the kubeconfig directory
var configDir string

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
	rootCmd.PersistentFlags().StringVarP(&configDir, "config-dir", "c", "", "Directory containing kubeconfig files")
}
