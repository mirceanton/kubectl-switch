package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// contextCmd represents the switch command to change Kubernetes contexts.
var contextCmd = &cobra.Command{
	Use:     "context",
	Aliases: []string{"ctx"}, // Add this line to define the alias
	Short:   "Switch the active Kubernetes context",
	Args:    cobra.MaximumNArgs(1), // Accept at most one argument (the context name)
	Run: func(cmd *cobra.Command, args []string) {
		// Determine the kubeconfig directory
		if configDir == "" {
			configDir = os.Getenv("KUBECONFIG_DIR")
			if configDir == "" {
				log.Fatal("kubeconfig directory not provided.")
				log.Fatal("Please provide the directory containing kubeconfig files via the --config-dir flag or KUBECONFIG_DIR environment variable")
			}
		}
		if strings.HasPrefix(configDir, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Failed to determine home directory: %v", err)
			}
			configDir = filepath.Join(homeDir, configDir[2:])
		}

		// Get all kubeconfig files in the config directory
		files, err := os.ReadDir(configDir)
		if err != nil {
			log.Fatalf("Failed to read directory: %v", err)
		}

		// Parse all kubeconfig files in the directory
		contextMap := make(map[string]string) // Map storing the context name (key) and the corresponding kubeconfig file (value)
		var contextNames []string             // List of context names for the interactive prompt

		for _, file := range files {
			// Skip directories
			if file.IsDir() {
				continue
			}

			// Skip files that are not YAML
			if filepath.Ext(file.Name()) != ".yaml" && filepath.Ext(file.Name()) != ".yml" {
				continue
			}

			// Parse the kubeconfig file
			path := filepath.Join(configDir, file.Name())
			kubeconfig, err := clientcmd.LoadFromFile(path)
			if err != nil {
				log.WithFields(log.Fields{"file": file.Name()}).Warnf("Failed to parse kubeconfig file: %v", err)
				continue
			}

			// Add context details to the map
			for contextName := range kubeconfig.Contexts {
				if _, exists := contextMap[contextName]; exists {
					log.Fatalf("Duplicate context name '%s' found in files:\n- %s\n- %s", contextName, contextMap[contextName], path)
				}
				contextMap[contextName] = path
				contextNames = append(contextNames, contextName)
			}
		}

		// Check if any contexts were found
		if len(contextMap) == 0 {
			log.Fatal("No kubernetes contexts found in the provided directory: ", configDir)
		}

		// Determine the target context
		var selectedContext string
		if len(args) == 1 {
			// Non-interactive mode: use the provided cluster name
			selectedContext = args[0]
			if _, exists := contextMap[selectedContext]; !exists {
				log.Fatalf("Context '%s' not found", selectedContext)
			}
		} else {
			// Interactive mode: show list of clusters
			prompt := &survey.Select{
				Message: "Choose a context:",
				Options: contextNames,
			}
			err = survey.AskOne(prompt, &selectedContext)
			if err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}

		// Determine the target location for copying the file
		destPath := os.Getenv("KUBECONFIG")
		if destPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Failed to determine home directory: %v", err)
			}
			destPath = filepath.Join(homeDir, ".kube", "config")
		}

		// Ensure the destination directory exists
		destDir := filepath.Dir(destPath)
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", destDir, err)
		}

		// Load the kubeconfig file for the selected context
		kubeconfig, err := clientcmd.LoadFromFile(contextMap[selectedContext])
		if err != nil {
			log.WithFields(log.Fields{"source": contextMap[selectedContext]}).Fatalf("Failed to parse kubeconfig file: %v", err)
		}

		// Update the current context
		kubeconfig.CurrentContext = selectedContext

		// Write the updated kubeconfig back to the file
		err = clientcmd.WriteToFile(*kubeconfig, destPath)
		if err != nil {
			log.Fatalf("Error writing kubeconfig file: %v", err)
		}

		log.Infof("Switched to context '%s'", selectedContext)
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
