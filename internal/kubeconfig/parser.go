package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"
)

// ContextDetails holds information about a Kubernetes context.
type ContextDetails struct {
	Name     string // Name of the context
	FilePath string // Path to the kubeconfig file containing the context
}

// ParseKubeconfig reads a kubeconfig file and extracts context details.
func ParseKubeconfig(filePath string) ([]ContextDetails, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	var kubeconfig clientcmdapi.Config
	if err := yaml.Unmarshal(data, &kubeconfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kubeconfig from file %s: %v", filePath, err)
	}

	var contexts []ContextDetails
	for _, context := range kubeconfig.Contexts {
		contexts = append(contexts, ContextDetails{
			Name:     context.Name,
			FilePath: filePath,
		})
	}

	return contexts, nil
}

// ParseKubeconfigs parses all kubeconfig files in the specified directory.
// It returns a list of context details from all files and checks for duplicate context names.
func ParseKubeconfigs(configDir string) ([]ContextDetails, error) {
	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", configDir, err)
	}

	var allContextDetails []ContextDetails
	duplicateCheck := make(map[string]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(configDir, file.Name())
		if filepath.Ext(filePath) != ".yaml" && filepath.Ext(filePath) != ".yml" {
			continue
		}

		contextDetails, err := ParseKubeconfig(filePath)
		if err != nil {
			return nil, fmt.Errorf("error parsing kubeconfig file %s: %v", filePath, err)
		}

		for _, contextDetail := range contextDetails {
			if existingFile, exists := duplicateCheck[contextDetail.Name]; exists {
				return nil, fmt.Errorf("duplicate context name '%s' found in files:\n- %s\n- %s", contextDetail.Name, existingFile, filePath)
			}
			duplicateCheck[contextDetail.Name] = filePath
			allContextDetails = append(allContextDetails, contextDetail)
		}
	}

	return allContextDetails, nil
}

// CopyConfig copies the specified kubeconfig file to the destination path.
func CopyConfig(srcPath, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read kubeconfig file %s: %v", srcPath, err)
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write kubeconfig to %s: %v", destPath, err)
	}

	return nil
}
