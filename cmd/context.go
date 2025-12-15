package cmd

import (
	"github.com/mirceanton/kubectl-switch/v2/internal/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:               "context",
	Aliases:           []string{"ctx"},
	Short:             "Switch the active Kubernetes context",
	ValidArgsFunction: getContextCompletions,
	Args:              cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := configManager.LoadContexts(); err != nil {
			log.Fatalf("Failed to load contexts: %v", err)
		}

		contextNames := configManager.GetAllContexts()
		if len(contextNames) == 0 {
			log.Fatal("No kubernetes contexts found in the provided directory")
		}

		var selectedContext string
		if len(args) == 1 {
			selectedContext = args[0]
		} else {
			currentContext := configManager.GetCurrentContext()
			selected, err := ui.Select("Choose a context:", contextNames, currentContext, appConfig.PageSize)
			if err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
			selectedContext = selected
		}

		if err := configManager.SwitchToContext(selectedContext); err != nil {
			log.Fatalf("Failed to switch context: %v", err)
		}

		log.Infof("Switched to context '%s'", selectedContext)
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}

func getContextCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if err := configManager.LoadContexts(); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return configManager.GetAllContexts(), cobra.ShellCompDirectiveNoFileComp
}
