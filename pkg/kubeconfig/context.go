package kubeconfig

import (
	"os"
	"path/filepath"
	"strings"
)

// ContextManager handles operations related to Kubernetes contexts
type ContextManager struct {
	kubeconfigManager *Manager
}

// NewContextManager creates a new context manager
func NewContextManager(km *Manager) *ContextManager {
	return &ContextManager{
		kubeconfigManager: km,
	}
}

// GetContextsFromDir gets all contexts from kubeconfig files in a directory
func (cm *ContextManager) GetContextsFromDir(configDir string) (map[string]string, []string, error) {
	// Handle ~/ in path
	if strings.HasPrefix(configDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, err
		}
		configDir = filepath.Join(homeDir, configDir[2:])
	}

	return cm.kubeconfigManager.GetContextsFromDir(configDir)
}

// SwitchContext switches to the specified context
func (cm *ContextManager) SwitchContext(contextFilePath string, contextName string) error {
	// Ensure the destination directory exists
	if err := cm.kubeconfigManager.EnsureDir(); err != nil {
		return err
	}

	// Switch to the selected context
	if err := cm.kubeconfigManager.SwitchContext(contextFilePath, contextName); err != nil {
		return err
	}

	return nil
}

// ValidateConfigDir ensures the kubeconfig directory is valid and available
func (cm *ContextManager) ValidateConfigDir(configDir string) (string, error) {
	if configDir == "" {
		configDir = os.Getenv("KUBECONFIG_DIR")
		if configDir == "" {
			return "", ErrNoConfigDir
		}
	}

	// Handle ~/ in path
	if strings.HasPrefix(configDir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, configDir[2:])
	}

	// Check if the directory exists
	info, err := os.Stat(configDir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", ErrNotADirectory
	}

	return configDir, nil
}
