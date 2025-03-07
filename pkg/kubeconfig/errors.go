package kubeconfig

import (
	"errors"
)

// Common errors for kubeconfig operations
var (
	ErrNoConfigDir      = errors.New("kubeconfig directory not provided, please provide the directory containing kubeconfig files via the --config-dir flag or KUBECONFIG_DIR environment variable")
	ErrNotADirectory    = errors.New("the provided path is not a directory")
	ErrNoPreviousConfig = errors.New("no previous configuration found")
)
