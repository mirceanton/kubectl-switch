package kubeconfig

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NamespaceManager handles operations related to Kubernetes namespaces
type NamespaceManager struct {
	kubeconfigManager *Manager
}

// NewNamespaceManager creates a new namespace manager
func NewNamespaceManager(km *Manager) *NamespaceManager {
	return &NamespaceManager{
		kubeconfigManager: km,
	}
}

// GetNamespaces gets all namespaces from the current context
func (nm *NamespaceManager) GetNamespaces() ([]string, error) {
	kubeconfigPath := nm.kubeconfigManager.GetPath()

	// Build the Kubernetes client configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Get all namespaces
	var namespaceNames []string
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ns := range namespaces.Items {
		namespaceNames = append(namespaceNames, ns.Name)
	}

	return namespaceNames, nil
}

// SwitchNamespace switches the namespace in the current context
func (nm *NamespaceManager) SwitchNamespace(namespace string) error {
	return nm.kubeconfigManager.SwitchNamespace(namespace)
}
