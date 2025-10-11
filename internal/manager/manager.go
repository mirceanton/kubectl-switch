package manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Manager handles kubeconfig file operations and Kubernetes context switching.
type Manager struct {
	kubeconfigPath string
	kubeconfigDir  string
	contextMap     map[string]string
	contextNames   []string
}

// Custom errors
var (
	ErrNoConfigDir      = errors.New("kubeconfig directory not provided, please provide the directory containing kubeconfig files via the --config-dir flag or KUBECONFIG_DIR environment variable")
	ErrNotADirectory    = errors.New("the provided path is not a directory")
	ErrNoPreviousConfig = errors.New("no previous configuration found")
	ErrContextNotFound  = errors.New("context not found")
)

// NewManager creates a new kubeconfig Manager instance.
// It validates and loads the config directory if provided.
func NewManager(configDir string) (*Manager, error) {
	m := &Manager{
		kubeconfigPath: getKubeconfigPath(),
	}

	// Validate and set config directory if provided
	if configDir != "" || os.Getenv("KUBECONFIG_DIR") != "" {
		validatedDir, err := m.validateConfigDir(configDir)
		if err != nil {
			return nil, err
		}
		m.kubeconfigDir = validatedDir

		// Load contexts from directory
		if err := m.loadContexts(); err != nil {
			return nil, fmt.Errorf("failed to load contexts: %w", err)
		}
	}

	return m, nil
}

// GetAllContexts returns the available context names.
func (m *Manager) GetAllContexts() []string {
	return m.contextNames
}

// GetNamespacesForCurrentContext retrieves all namespaces from the current Kubernetes cluster.
func (m *Manager) GetNamespacesForCurrentContext() ([]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", m.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	namespaceNames := make([]string, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}

	return namespaceNames, nil
}

// SwitchToContext switches to the specified Kubernetes context.
func (m *Manager) SwitchToContext(contextName string) error {
	contextFilePath, exists := m.contextMap[contextName]
	if !exists {
		return fmt.Errorf("%w: %s", ErrContextNotFound, contextName)
	}

	if err := m.backup(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	kubeconfig, err := clientcmd.LoadFromFile(contextFilePath)
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig from %s: %w", contextFilePath, err)
	}

	kubeconfig.CurrentContext = contextName

	if err := clientcmd.WriteToFile(*kubeconfig, m.kubeconfigPath); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return nil
}

// SwitchToNamespace switches the namespace for the current context.
func (m *Manager) SwitchToNamespace(namespace string) error {
	if err := m.backup(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	kubeconfig, err := clientcmd.LoadFromFile(m.kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to load current kubeconfig: %w", err)
	}

	if kubeconfig.Contexts[kubeconfig.CurrentContext] == nil {
		return fmt.Errorf("current context '%s' not found", kubeconfig.CurrentContext)
	}

	kubeconfig.Contexts[kubeconfig.CurrentContext].Namespace = namespace

	if err := clientcmd.WriteToFile(*kubeconfig, m.kubeconfigPath); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return nil
}

// Restore swaps the current kubeconfig with the previous backup.
func (m *Manager) Restore() error {
	prevPath := m.getPreviousPath()

	if _, err := os.Stat(prevPath); os.IsNotExist(err) {
		return ErrNoPreviousConfig
	}

	currentConfig, err := os.ReadFile(m.kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current config: %w", err)
	}

	prevConfig, err := os.ReadFile(prevPath)
	if err != nil {
		return fmt.Errorf("failed to read previous config: %w", err)
	}

	if err := os.WriteFile(m.kubeconfigPath, prevConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write current config: %w", err)
	}

	if err := os.WriteFile(prevPath, currentConfig, 0o600); err != nil {
		return fmt.Errorf("failed to write previous config: %w", err)
	}

	return nil
}

// ================================================================================================
// Helper functions
// ================================================================================================
// getKubeconfigPath returns the path to the current kubeconfig file.
func getKubeconfigPath() string {
	path := os.Getenv("KUBECONFIG")
	if path == "" {
		path = "~/.kube/config"
	}

	expandedPath, err := expandPath(path)
	if err != nil {
		log.Fatalf("Failed to expand KUBECONFIG path: %v", err)
	}
	return expandedPath
}

// getPreviousPath returns the path where the previous kubeconfig backup is stored.
func (m *Manager) getPreviousPath() string {
	return filepath.Join(filepath.Dir(m.kubeconfigPath), "config.previous")
}

// backup backs up the current kubeconfig to config.previous.
func (m *Manager) backup() error {
	if _, err := os.Stat(m.kubeconfigPath); os.IsNotExist(err) {
		return fmt.Errorf("kubeconfig file does not exist: %w", err)
	}

	data, err := os.ReadFile(m.kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read current kubeconfig: %w", err)
	}

	if err := os.WriteFile(m.getPreviousPath(), data, 0o600); err != nil {
		return fmt.Errorf("failed to write previous kubeconfig: %w", err)
	}

	return nil
}

// expandPath expands a path starting with ~/ to the full home directory path.
func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, path[2:]), nil
}

// validateConfigDir validates and expands the provided config directory path.
func (m *Manager) validateConfigDir(configDir string) (string, error) {
	if configDir == "" {
		configDir = os.Getenv("KUBECONFIG_DIR")
		if configDir == "" {
			return "", ErrNoConfigDir
		}
	}

	expandedPath, err := expandPath(configDir)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat config directory: %w", err)
	}

	if !info.IsDir() {
		return "", ErrNotADirectory
	}

	return expandedPath, nil
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
