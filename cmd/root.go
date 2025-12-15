package cmd

import (
	"os"

	"github.com/mirceanton/kubectl-switch/v2/internal/config"
	"github.com/mirceanton/kubectl-switch/v2/internal/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version       string
	configManager *manager.Manager
	appConfig     *config.Config
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
		// Load configuration
		var err error
		appConfig, err = config.Load()
		if err != nil {
			return err
		}

		// Set up logging
		log.SetLevel(appConfig.LogLevel)
		log.SetFormatter(appConfig.LogFormat)

		// Create manager with config
		configManager, err = manager.NewManager(appConfig.Kubeconfig, appConfig.KubeconfigDir)
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
	// Initialize Viper
	cobra.OnInitialize(config.Init)

	// Bind flags to Viper
	rootCmd.PersistentFlags().String("kubeconfig-dir", "", "Directory containing kubeconfig files (env: KUBECONFIG_DIR)")
	err := viper.BindPFlag("kubeconfig-dir", rootCmd.PersistentFlags().Lookup("kubeconfig-dir"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}

	rootCmd.PersistentFlags().String("kubeconfig", "", "Currently active kubeconfig file (env: KUBECONFIG)")
	err = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}

	rootCmd.PersistentFlags().String("log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic) (env: LOG_LEVEL)")
	err = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}

	rootCmd.PersistentFlags().String("log-format", "text", "Log format (text, json) (env: LOG_FORMAT)")
	err = viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}

	rootCmd.PersistentFlags().Int("page-size", 10, "Number of items to show per page in selection prompts (env: PAGE_SIZE)")
	err = viper.BindPFlag("page-size", rootCmd.PersistentFlags().Lookup("page-size"))
	if err != nil {
		log.Fatalf("Failed to bind flag: %v", err)
	}
}
