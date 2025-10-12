package manager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Manager handles kubeconfig file operations and Kubernetes context switching.
type Manager struct {
	kubeconfigPath string
	backupPath     string
	kubeconfigDir  string
	contextMap     map[string]string
	contextNames   []string
	namespaceNames []string
}

// NewManager creates a new kubeconfig Manager instance.
// It takes the validated configuration paths and loads available contexts.
func NewManager(kubeconfigPath, kubeconfigDir string) (*Manager, error) {
	m := &Manager{
		kubeconfigPath: kubeconfigPath,
		kubeconfigDir:  kubeconfigDir,
		backupPath:     kubeconfigPath + ".previous",
	}

	// Load available contexts from the config directory
	if err := m.loadContexts(); err != nil {
		return nil, fmt.Errorf("failed to load contexts: %w", err)
	}

	// Load available namespaces from the current cluster
	// This may fail if no kubeconfig is active, which is OK
	if err := m.loadNamespaces(); err != nil {
		log.Warnf("Failed to load namespaces: %v", err)
		log.Warn("Namespace operations will not be available until a valid context is selected")
		m.namespaceNames = []string{} // Initialize to empty slice
	}

	return m, nil
}

// GetAllContexts returns the available context names.
func (m *Manager) GetAllContexts() []string {
	return m.contextNames
}

// GetAllNamespaces retrieves all namespaces from the current Kubernetes cluster.
func (m *Manager) GetAllNamespaces() []string {
	return m.namespaceNames
}

// SwitchToContext switches to the specified Kubernetes context.
func (m *Manager) SwitchToContext(contextName string) error {
	// Find the kubeconfig file containing the desired context
	contextFilePath, exists := m.contextMap[contextName]
	if !exists {
		return fmt.Errorf("context '%s' not found", contextName)
	}

	// Load the kubeconfig file containing the desired context
	kubeconfig, err := clientcmd.LoadFromFile(contextFilePath)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig from %s: %w", contextFilePath, err)
	}

	// Update the current context in the loaded kubeconfig
	kubeconfig.CurrentContext = contextName

	// Backup current config
	if err := m.backup(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	// Write updated kubeconfig back to the main kubeconfig file
	if err := clientcmd.WriteToFile(*kubeconfig, m.kubeconfigPath); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return nil
}

// SwitchToNamespace switches the namespace for the current context.
func (m *Manager) SwitchToNamespace(namespace string) error {
	// Parse current kubeconfig
	kubeconfig, err := clientcmd.LoadFromFile(m.kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to load current kubeconfig: %w", err)
	}

	// Update namespace for current context
	kubeconfig.Contexts[kubeconfig.CurrentContext].Namespace = namespace

	// Backup current config
	if err := m.backup(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	// Write updated kubeconfig back to file
	if err := clientcmd.WriteToFile(*kubeconfig, m.kubeconfigPath); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return nil
}

// Restore swaps the current kubeconfig with the previous backup.
func (m *Manager) Restore() error {
	// Read current kubeconfig
	currentConfig, err := os.ReadFile(m.kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current config: %w", err)
	}

	// Read backup kubeconfig
	prevConfig, err := os.ReadFile(m.backupPath)
	if err != nil {
		return fmt.Errorf("failed to read previous config: %w", err)
	}

	// Swap the files
	if err := os.WriteFile(m.kubeconfigPath, prevConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write current config: %w", err)
	}
	if err := os.WriteFile(m.backupPath, currentConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write previous config: %w", err)
	}

	return nil
}

// ================================================================================================
// Helper functions
// ================================================================================================

// backup backs up the current kubeconfig to config.previous.
func (m *Manager) backup() error {
	data, err := os.ReadFile(m.kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current kubeconfig: %w", err)
	}

	if err := os.WriteFile(m.backupPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write previous kubeconfig: %w", err)
	}

	return nil
}

// loadContexts scans the config directory for kubeconfig files and loads all available contexts.
func (m *Manager) loadContexts() error {
	m.contextMap = make(map[string]string)
	m.contextNames = nil

	files, err := os.ReadDir(m.kubeconfigDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(m.kubeconfigDir, file.Name())
		kubeconfig, err := clientcmd.LoadFromFile(path)
		if err != nil {
			log.WithField("file", file.Name()).Warnf("Failed to parse kubeconfig file: %v", err)
			continue
		}

		for contextName := range kubeconfig.Contexts {
			if existingPath, exists := m.contextMap[contextName]; exists {
				log.Warnf("Duplicate context name '%s' found in files:\n  - %s\n  - %s",
					contextName, existingPath, path)
				continue
			}
			m.contextMap[contextName] = path
			m.contextNames = append(m.contextNames, contextName)
		}
	}

	return nil
}

// loadNamespaces loads all namespaces from the current Kubernetes cluster.
func (m *Manager) loadNamespaces() error {
	config, err := clientcmd.BuildConfigFromFlags("", m.kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	m.namespaceNames = make([]string, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		m.namespaceNames = append(m.namespaceNames, ns.Name)
	}

	return nil
}
