package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// namespaceCmd represents the command to change Kubernetes namespaces.
var namespaceCmd = &cobra.Command{
	Use:     "namespace",
	Aliases: []string{"ns"}, // Add this line to define the alias
	Short:   "Switch the active Kubernetes namespace",
	Args:    cobra.MaximumNArgs(1), // Accept at most one argument (the namespace name)
	Run: func(cmd *cobra.Command, args []string) {
		// Determine the location of the kubeconfig file
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Failed to determine home directory: %v", err)
			}
			kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
		}

		// Build the Kubernetes client configuration
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %v", err)
		}

		// Create the Kubernetes clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating Kubernetes client: %v", err)
		}

		// Get all namespaces
		var namespaceNames []string
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		for _, ns := range namespaces.Items {
			namespaceNames = append(namespaceNames, ns.Name)
		}
		if err != nil {
			log.Fatalf("Error listing namespaces: %v", err)
		}

		// Determine the target namespace
		var selectedNamespace string
		if len(args) == 1 {
			// Non-interactive mode: use the provided ma,es[ace] name
			selectedNamespace = args[0]
		} else {
			// Interactive mode: show list of clusters
			prompt := &survey.Select{
				Message: "Choose a namespace:",
				Options: namespaceNames,
			}
			err = survey.AskOne(prompt, &selectedNamespace)
			if err != nil {
				log.Fatalf("Failed to get user input: %v", err)
			}
		}

		// Load the kubeconfig file
		kubeconfig, err := clientcmd.LoadFromFile(kubeconfigPath)
		if err != nil {
			log.WithFields(log.Fields{"source": kubeconfigPath}).Fatalf("Failed to parse kubeconfig file: %v", err)
		}

		// Update the kubeconfig
		kubeconfig.Contexts[kubeconfig.CurrentContext].Namespace = selectedNamespace

		// Write the updated kubeconfig back to the file
		err = clientcmd.WriteToFile(*kubeconfig, kubeconfigPath)
		if err != nil {
			log.Fatalf("Error writing kubeconfig file: %v", err)
		}

		log.Infof("Switched to namespace '%s'", selectedNamespace)
	},
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}
