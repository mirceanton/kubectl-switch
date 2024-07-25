package prompt

import (
	"github.com/mirceanton/kube-switcher/internal/kubeconfig"

	"github.com/manifoldco/promptui"
)

// SelectContext prompts the user to select a Kubernetes context.
func SelectContext(contexts []kubeconfig.ContextDetails) (string, error) {
	items := make([]string, len(contexts))
	for i, context := range contexts {
		items[i] = context.Name
	}

	prompt := promptui.Select{
		Label: "Select Kubernetes Context",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}
