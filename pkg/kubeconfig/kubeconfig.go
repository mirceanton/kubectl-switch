package kubeconfig

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
)

// Manager handles operations related to kubeconfig files
type Manager struct{}

// NewManager creates a new kubeconfig manager
func NewManager() *Manager {
	return &Manager{}
}

// GetPath returns the path to the current kubeconfig file
func (m *Manager) GetPath() string {
	path := os.Getenv("KUBECONFIG")
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to determine home directory: %v", err)
		}
		path = filepath.Join(homeDir, ".kube", "config")
	}
	return path
}

// GetPreviousPath returns the path to the previous kubeconfig file
func (m *Manager) GetPreviousPath() string {
	return filepath.Join(filepath.Dir(m.GetPath()), "config.previous")
}

// SaveCurrent backs up the current kubeconfig file
func (m *Manager) SaveCurrent() error {
	kubeconfigPath := m.GetPath()
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return err
	}

	prevPath := m.GetPreviousPath()
	data, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return err
	}

	return os.WriteFile(prevPath, data, 0o600)
}

// SwitchToPrevious switches back to the previous Kubernetes config
func (m *Manager) SwitchToPrevious() error {
	kubeconfigPath := m.GetPath()
	prevPath := m.GetPreviousPath()

	// Check if previous config exists
	if _, err := os.Stat(prevPath); os.IsNotExist(err) {
		return ErrNoPreviousConfig
	}

	// Swap the current and previous configs
	tempPath := filepath.Join(filepath.Dir(kubeconfigPath), "config.temp")

	// Read current config
	currentConfig, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return err
	}

	// Read previous config
	prevConfig, err := os.ReadFile(prevPath)
	if err != nil {
		return err
	}

	// Write current to temp
	if err := os.WriteFile(tempPath, currentConfig, 0o600); err != nil {
		return err
	}

	// Write previous to current
	if err := os.WriteFile(kubeconfigPath, prevConfig, 0o600); err != nil {
		return err
	}

	// Write temp to previous
	if err := os.WriteFile(prevPath, currentConfig, 0o600); err != nil {
		return err
	}

	// Remove temp file
	if err := os.Remove(tempPath); err != nil {
		log.Warnf("Failed to remove temporary file: %v", err)
	}

	return nil
}

// EnsureDir ensures that the directory for the kubeconfig file exists
func (m *Manager) EnsureDir() error {
	destDir := filepath.Dir(m.GetPath())
	return os.MkdirAll(destDir, 0o755)
}

// GetContextsFromDir gets all contexts from kubeconfig files in a directory
func (m *Manager) GetContextsFromDir(configDir string) (map[string]string, []string, error) {
	contextMap := make(map[string]string)
	var contextNames []string

	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(configDir, file.Name())
		kubeconfig, err := clientcmd.LoadFromFile(path)
		if err != nil {
			log.WithFields(log.Fields{"file": file.Name()}).Warnf("Failed to parse kubeconfig file: %v", err)
			continue
		}

		for contextName := range kubeconfig.Contexts {
			if _, exists := contextMap[contextName]; exists {
				log.Warnf("Duplicate context name '%s' found in files:\n- %s\n- %s", contextName, contextMap[contextName], path)
				continue
			}
			contextMap[contextName] = path
			contextNames = append(contextNames, contextName)
		}
	}

	return contextMap, contextNames, nil
}

// SwitchContext switches to the specified context
func (m *Manager) SwitchContext(contextFilePath string, contextName string) error {
	// Save the current config as the previous config
	if err := m.SaveCurrent(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	// Load the kubeconfig file for the selected context
	kubeconfig, err := clientcmd.LoadFromFile(contextFilePath)
	if err != nil {
		return err
	}

	// Update the current context
	kubeconfig.CurrentContext = contextName

	// Write the updated kubeconfig back to the file
	return clientcmd.WriteToFile(*kubeconfig, m.GetPath())
}

// SwitchNamespace switches the namespace in the current context
func (m *Manager) SwitchNamespace(namespace string) error {
	// Save the current config as the previous config
	if err := m.SaveCurrent(); err != nil {
		log.Warnf("Failed to save current configuration as previous: %v", err)
	}

	// Load the kubeconfig file
	kubeconfig, err := clientcmd.LoadFromFile(m.GetPath())
	if err != nil {
		return err
	}

	// Update the namespace in the current context
	kubeconfig.Contexts[kubeconfig.CurrentContext].Namespace = namespace

	// Write the updated kubeconfig back to the file
	return clientcmd.WriteToFile(*kubeconfig, m.GetPath())
}
