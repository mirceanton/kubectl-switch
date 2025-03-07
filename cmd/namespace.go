package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/mirceanton/kubectl-switch/pkg/kubeconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var namespaceManager = kubeconfig.NewNamespaceManager(manager)

var namespaceCmd = &cobra.Command{
	Use:     "namespace",
	Aliases: []string{"ns"},
	Short:   "Switch the active Kubernetes namespace",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get all namespaces from the current context
		namespaceNames, err := namespaceManager.GetNamespaces()
		if err != nil {
			log.Fatalf("Error listing namespaces: %v", err)
		}

		// Determine the target namespace
		var selectedNamespace string
		if len(args) == 1 {
			selectedNamespace = args[0]
		} else {
			prompt := &survey.Select{
				Message: "Choose a namespace:",
				Options: namespaceNames,
			}
			err = survey.AskOne(prompt, &selectedNamespace)
			if err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}

		// Switch to the selected namespace
		if err := namespaceManager.SwitchNamespace(selectedNamespace); err != nil {
			log.Fatalf("Failed to switch namespace: %v", err)
		}

		log.Infof("Switched to namespace '%s'", selectedNamespace)
	},
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}
