package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var namespaceCmd = &cobra.Command{
	Use:               "namespace",
	Aliases:           []string{"ns"},
	Short:             "Switch the active Kubernetes namespace",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: getNamespaceCompletions,
	Run: func(cmd *cobra.Command, args []string) {
		namespaceNames := configManager.GetAllNamespaces()
		if len(namespaceNames) == 0 {
			log.Fatal("No kubernetes namespaces found in the current cluster")
		}

		var selectedNamespace string
		if len(args) == 1 {
			selectedNamespace = args[0]
		} else {
			prompt := &survey.Select{
				Message: "Choose a namespace:",
				Options: namespaceNames,
			}
			if err := survey.AskOne(prompt, &selectedNamespace); err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}

		if err := configManager.SwitchToNamespace(selectedNamespace); err != nil {
			log.Fatalf("Failed to switch namespace: %v", err)
		}

		log.Infof("Switched to namespace '%s'", selectedNamespace)
	},
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}

func getNamespaceCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return configManager.GetAllNamespaces(), cobra.ShellCompDirectiveNoFileComp
}
