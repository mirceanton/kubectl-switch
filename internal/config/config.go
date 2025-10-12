package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	KubeconfigDir string
	Kubeconfig    string
	LogLevel      log.Level
	LogFormat     log.Formatter
}

const (
	// Configuration keys
	keyKubeconfigDir = "kubeconfig-dir"
	keyKubeconfig    = "kubeconfig"
	keyLogLevel      = "log-level"
	keyLogFormat     = "log-format"

	// Default values
	defaultLogLevel  = "info"
	defaultLogFormat = "text"
)

var (
	defaultKubeconfigDir = filepath.Join(os.Getenv("HOME"), ".kube", "configs/")
	defaultKubeconfig    = filepath.Join(os.Getenv("HOME"), ".kube", "config")
)

// Init initializes Viper configuration
func Init() {
	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Replace hyphens with underscores in env vars
	// This allows --kubeconfig-dir to map to KUBECONFIG_DIR
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Set default values
	viper.SetDefault(keyKubeconfigDir, defaultKubeconfigDir)
	viper.SetDefault(keyKubeconfig, defaultKubeconfig)
	viper.SetDefault(keyLogLevel, defaultLogLevel)
	viper.SetDefault(keyLogFormat, defaultLogFormat)
}

// Load returns the current configuration
func Load() (*Config, error) {
	cfg := &Config{}

	// Parse log level
	levelStr := viper.GetString(keyLogLevel)
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s", levelStr)
	}
	cfg.LogLevel = level

	// Parse log format
	formatStr := viper.GetString(keyLogFormat)
	switch strings.ToLower(formatStr) {
	case "json":
		cfg.LogFormat = &log.JSONFormatter{}
	case "text":
		cfg.LogFormat = &log.TextFormatter{FullTimestamp: true}
	default:
		return nil, fmt.Errorf("invalid log format: %s", formatStr)
	}

	// Expand and validate kubeconfig dir path
	cfg.KubeconfigDir, err = expandPath(viper.GetString(keyKubeconfigDir))
	if err != nil {
		return nil, fmt.Errorf("failed to expand kubeconfig directory path: %w", err)
	}
	if err := cfg.validateKubeconfigDir(); err != nil {
		return nil, err
	}

	// Expand and validate kubeconfig dir path
	cfg.Kubeconfig, err = expandPath(viper.GetString(keyKubeconfig))
	if err != nil {
		return nil, fmt.Errorf("failed to expand kubeconfig path: %w", err)
	}
	if err := cfg.validateKubeconfig(); err != nil {
		log.Warnf("Kubeconfig file validation failed: %v", err)
		log.Warn("You can still switch contexts, but namespace operations will not be available")
	}

	return cfg, nil
}

// validateKubeconfigDir validates that the kubeconfig directory exists and is a directory
func (c *Config) validateKubeconfigDir() error {
	info, err := os.Stat(c.KubeconfigDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("kubeconfig directory does not exist: %s", c.KubeconfigDir)
		}
		return fmt.Errorf("failed to stat kubeconfig directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("kubeconfig directory path is not a directory: %s", c.KubeconfigDir)
	}

	return nil
}

// validateKubeconfig validates that the kubeconfig file exists and is a file
func (c *Config) validateKubeconfig() error {
	info, err := os.Stat(c.Kubeconfig)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("kubeconfig file does not exist: %s", c.Kubeconfig)
		}
		return fmt.Errorf("failed to stat kubeconfig file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("kubeconfig path is a directory, expected a file: %s", c.Kubeconfig)
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
