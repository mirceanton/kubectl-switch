package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var configDir string // Global flag for the kubeconfig directory
var version string   // The version of the tool, set at build time

var rootCmd = &cobra.Command{
	Use:     "kube-switcher",
	Short:   "A tool to switch Kubernetes contexts",
	Long:    `kube-switcher is a CLI tool to switch Kubernetes contexts from multiple kubeconfig files.`,
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
