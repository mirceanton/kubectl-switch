package cmd

import (
	"github.com/AlecAivazis/survey/v2"
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
		contextNames := configManager.GetAllContexts()
		if len(contextNames) == 0 {
			log.Fatal("No kubernetes contexts found in the provided directory")
		}

		var selectedContext string
		if len(args) == 1 {
			selectedContext = args[0]
		} else {
			prompt := &survey.Select{
				Message:  "Choose a context:",
				Options:  contextNames,
				PageSize: appConfig.PageSize,
			}
			if err := survey.AskOne(prompt, &selectedContext); err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
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
	return configManager.GetAllContexts(), cobra.ShellCompDirectiveNoFileComp
}
