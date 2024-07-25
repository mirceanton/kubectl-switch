package cmd

import (
	"fmt"
	"os"

	"github.com/mirceanton/kube-switcher/internal/kubeconfig"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command to change Kubernetes contexts.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available Kubernetes context",
	Run: func(cmd *cobra.Command, args []string) {
		// Determine the kubeconfig directory
		if configDir == "" {
			configDir = os.Getenv("KUBESWITCHER_CONFIG_DIR")
			if configDir == "" {
				log.Fatal("KUBESWITCHER_CONFIG_DIR environment variable is not set")
			}
		}

		// Parse all kubeconfigs in the specified directory
		contexts, err := kubeconfig.ParseKubeconfigs(configDir)
		if err != nil {
			log.Fatalf("Error parsing kubeconfig files: %v", err)
		}

		for _, context := range contexts {
			fmt.Printf("- %s [%s]\n", context.Name, context.FilePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
