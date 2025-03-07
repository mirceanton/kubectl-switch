package cmd

import (
	"os"

	"github.com/mirceanton/kubectl-switch/pkg/kubeconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configDir string // Global flag for the kubeconfig directory
	version   string // The version of the tool, set at build time
	manager   = kubeconfig.NewManager()
)

var rootCmd = &cobra.Command{
	Use:     "kubectl switch",
	Aliases: []string{"kubectl swtich"},
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "kubectl switch",
	},
	Short:   "A tool to switch Kubernetes contexts",
	Long:    `kubectl-switch is a CLI tool to switch Kubernetes contexts from multiple kubeconfig files.`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if the "-" argument is provided to switch to the previous config
		if len(args) == 1 && args[0] == "-" {
			if err := manager.SwitchToPrevious(); err != nil {
				if err == kubeconfig.ErrNoPreviousConfig {
					log.Fatal("No previous configuration found")
				} else {
					log.Fatalf("Failed to switch to previous config: %v", err)
				}
			}
			return nil
		}
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configDir, "kubeconfig-dir", "", "", "Directory containing kubeconfig files")
}
