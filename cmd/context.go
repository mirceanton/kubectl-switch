package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/mirceanton/kubectl-switch/pkg/kubeconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var contextManager = kubeconfig.NewContextManager(manager)

var contextCmd = &cobra.Command{
	Use:               "context",
	Aliases:           []string{"ctx"},
	Short:             "Switch the active Kubernetes context",
	ValidArgsFunction: getContextCompletions,
	Args:              cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Validate and get the config directory
		validConfigDir, err := contextManager.ValidateConfigDir(configDir)
		if err != nil {
			log.Fatal(err)
		}

		// Get all contexts from kubeconfig files
		contextMap, contextNames, err := contextManager.GetContextsFromDir(validConfigDir)
		if err != nil {
			log.Fatalf("Failed to read kubeconfig files: %v", err)
		}

		if len(contextMap) == 0 {
			log.Fatal("No kubernetes contexts found in the provided directory: ", validConfigDir)
		}

		// Determine the target context
		var selectedContext string
		if len(args) == 1 {
			selectedContext = args[0]
			if _, exists := contextMap[selectedContext]; !exists {
				log.Fatalf("Context '%s' not found", selectedContext)
			}
		} else {
			prompt := &survey.Select{
				Message: "Choose a context:",
				Options: contextNames,
			}
			err = survey.AskOne(prompt, &selectedContext)
			if err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}

		// Switch to the selected context
		if err := contextManager.SwitchContext(contextMap[selectedContext], selectedContext); err != nil {
			log.Fatalf("Failed to switch context: %v", err)
		}

		log.Infof("Switched to context '%s'", selectedContext)
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}

// getContextCompletions provides bash completion for available contexts
func getContextCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	configDirFlag, _ := cmd.Flags().GetString("kubeconfig-dir")
	validConfigDir, err := contextManager.ValidateConfigDir(configDirFlag)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	_, contextNames, err := contextManager.GetContextsFromDir(validConfigDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return contextNames, cobra.ShellCompDirectiveNoFileComp
}
