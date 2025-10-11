package cmd

import (
	"os"

	"github.com/mirceanton/kubectl-switch/v2/internal/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configDir     string
	version       string
	configManager *manager.Manager
)

var rootCmd = &cobra.Command{
	Use: "kubectl-switch",
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "kubectl switch",
	},
	Short:   "A tool to switch Kubernetes contexts",
	Long:    `kubectl-switch is a CLI tool to switch Kubernetes contexts from multiple kubeconfig files.`,
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		configManager, err = manager.NewManager(configDir)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "-" {
			if err := configManager.Restore(); err != nil {
				log.Fatalf("Failed to switch to previous config: %v", err)
			}
			return nil
		}
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configDir, "kubeconfig-dir", "", "", "Directory containing kubeconfig files")
}
