package kubeconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create a temporary file with given content
func createTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestParseKubeconfig(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Config
contexts:
  - name: test-context
    context:
      cluster: test-cluster
      user: test-user
clusters:
  - name: test-cluster
    cluster:
      server: https://localhost:6443
users:
  - name: test-user
    user:
      token: test-token
`

	// Create a temporary kubeconfig file
	filePath := createTempFile(t, yamlContent)

	// Test the ParseKubeconfig function
	contexts, err := ParseKubeconfig(filePath)
	if err != nil {
		t.Fatalf("ParseKubeconfig() returned an error: %v", err)
	}

	if len(contexts) != 1 {
		t.Fatalf("Expected 1 context, got %d", len(contexts))
	}

	if contexts[0].Name != "test-context" {
		t.Fatalf("Expected context name 'test-context', got %s", contexts[0].Name)
	}
}

func TestParseKubeconfigs(t *testing.T) {
	// Create temporary kubeconfig files
	file1Content := `
apiVersion: v1
kind: Config
contexts:
  - name: context-1
    context:
      cluster: cluster-1
      user: user-1
clusters:
  - name: cluster-1
    cluster:
      server: https://192.168.1.1:6443
users:
  - name: user-1
    user:
      token: token-1
`
	file1Path := createTempFile(t, file1Content)

	file2Content := `
apiVersion: v1
kind: Config
contexts:
  - name: context-2
    context:
      cluster: cluster-2
      user: user-2
clusters:
  - name: cluster-2
    cluster:
      server: https://192.168.1.2:6443
users:
  - name: user-2
    user:
      token: token-2
`
	file2Path := createTempFile(t, file2Content)

	// Create a temporary directory and move files there
	dirPath := filepath.Join(os.TempDir(), "kubeconfigs")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dirPath) // Clean up the directory

	// Move temporary files to the directory
	if err := os.Rename(file1Path, filepath.Join(dirPath, "file1.yaml")); err != nil {
		t.Fatalf("Failed to move file1: %v", err)
	}
	if err := os.Rename(file2Path, filepath.Join(dirPath, "file2.yaml")); err != nil {
		t.Fatalf("Failed to move file2: %v", err)
	}

	// Test the ParseKubeconfigs function
	contexts, err := ParseKubeconfigs(dirPath)
	if err != nil {
		t.Fatalf("ParseKubeconfigs() returned an error: %v", err)
	}

	if len(contexts) != 2 {
		t.Fatalf("Expected 2 contexts, got %d", len(contexts))
	}

	expectedNames := []string{"context-1", "context-2"}
	for _, context := range contexts {
		if !strings.Contains(strings.Join(expectedNames, ","), context.Name) {
			t.Fatalf("Unexpected context name: %s", context.Name)
		}
	}
}

func TestParseKubeconfigs_DuplicateContexts(t *testing.T) {
	// Create temporary kubeconfig files with duplicate context names
	duplicateContextContent := `
apiVersion: v1
kind: Config
contexts:
  - name: duplicate-context
    context:
      cluster: cluster-a
      user: user-a
clusters:
  - name: cluster-a
    cluster:
      server: https://192.168.1.1:6443
users:
  - name: user-a
    user:
      token: token-a
`
	file1Path := createTempFile(t, duplicateContextContent)

	// Creating another file with the same context name
	file2Path := createTempFile(t, duplicateContextContent)

	// Create a temporary directory and move files there
	dirPath := filepath.Join(os.TempDir(), "kubeconfigs-duplicates")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(dirPath) // Clean up the directory

	// Move temporary files to the directory
	if err := os.Rename(file1Path, filepath.Join(dirPath, "file1.yaml")); err != nil {
		t.Fatalf("Failed to move file1: %v", err)
	}
	if err := os.Rename(file2Path, filepath.Join(dirPath, "file2.yaml")); err != nil {
		t.Fatalf("Failed to move file2: %v", err)
	}

	// Test that ParseKubeconfigs returns an error due to duplicate context names
	_, err := ParseKubeconfigs(dirPath)
	if err == nil {
		t.Fatal("Expected error due to duplicate context names, but got none")
	}

	expectedErrMsg := "duplicate context name 'duplicate-context' found"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Fatalf("Expected error message to contain '%s', but got: %v", expectedErrMsg, err)
	}
}

func TestCopyConfig(t *testing.T) {
	srcContent := `
apiVersion: v1
kind: Config
contexts:
  - name: test-context
    context:
      cluster: test-cluster
      user: test-user
clusters:
  - name: test-cluster
    cluster:
      server: https://localhost:6443
users:
  - name: test-user
    user:
      token: test-token
`
	srcPath := createTempFile(t, srcContent)

	destPath := filepath.Join(os.TempDir(), "copied-config.yaml")

	// Test CopyConfig function
	if err := CopyConfig(srcPath, destPath); err != nil {
		t.Fatalf("CopyConfig() returned an error: %v", err)
	}

	// Verify the content of the copied file
	destContent, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if !strings.Contains(string(destContent), "test-context") {
		t.Fatalf("Copied file content does not contain expected context")
	}
}
