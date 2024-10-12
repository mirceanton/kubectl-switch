package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	configDir string // Global flag for the kubeconfig directory
	version   string // The version of the tool, set at build time
)

var rootCmd = &cobra.Command{
	Use:     "kubectl-switch",
	Short:   "A tool to switch Kubernetes contexts",
	Long:    `kubectl-switch is a CLI tool to switch Kubernetes contexts from multiple kubeconfig files.`,
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add any global flags here
	rootCmd.PersistentFlags().StringVarP(&configDir, "kubeconfig-dir", "", "", "Directory containing kubeconfig files")
}
