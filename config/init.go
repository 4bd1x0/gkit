package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func init() {
	configDir, err := findConfigsDir()
	if err != nil {
		panic(fmt.Sprintf("failed to find config directory: %s", err))
	}

	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("error reading config files: %s", err))
	}

	env := os.Getenv(EnvService)
	if env == "" {
		env = viper.GetString("env")
	}

	// Merge environment-specific config file only when env is set
	if env != "" {
		envConfigFile := fmt.Sprintf("application.%s", env)
		viper.SetConfigName(envConfigFile)
		if err := viper.MergeInConfig(); err != nil {
			// Environment config file may not exist, this is normal, just log a warning
			fmt.Printf("Warning: No config file found for env [%s], using base config only\n", env)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("unable to decode config of env [%s] into struct: %v", env, err))
	}
}

// Environment variable constants
const (
	EnvConfigDir = "GO_CONFIG_DIR"  // Configuration file directory
	EnvService   = "GO_SERVICE_ENV" // Service environment (local/test/prod)
)

// findConfigsDir finds the configuration file directory
// 1. Priority: use environment variable GO_CONFIG_DIR
// 2. Check configs in the executable directory
// 3. Search upward for cmd directory, find configs in its parent directory
func findConfigsDir() (string, error) {
	// 1. Priority: use config directory specified by environment variable
	if configDir := os.Getenv(EnvConfigDir); configDir != "" {
		if stat, err := os.Stat(configDir); err == nil && stat.IsDir() {
			return configDir, nil
		}
		return "", fmt.Errorf("GO_CONFIG_DIR [%s] is not a valid directory", configDir)
	}

	// 2. Get the executable file directory
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)

	// 3. First check if configs exists in the executable directory
	configsPath := filepath.Join(exeDir, "configs")
	if stat, err := os.Stat(configsPath); err == nil && stat.IsDir() {
		return configsPath, nil
	}

	// 4. Search upward for cmd directory
	dir := exeDir
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached file system root, cmd directory not found
			break
		}

		// Check if current directory name is "cmd"
		if filepath.Base(dir) == "cmd" {
			// cmd's parent directory is the project root
			projectRoot := filepath.Dir(dir)
			configsPath = filepath.Join(projectRoot, "configs")
			if stat, err := os.Stat(configsPath); err == nil && stat.IsDir() {
				return configsPath, nil
			}
			// Found cmd directory but configs not found in project root, return error
			return "", fmt.Errorf("configs directory not found in project root: %s", projectRoot)
		}

		dir = parent
	}

	return "", fmt.Errorf("config directory not found. Neither found configs in %s nor cmd directory in parent paths", exeDir)
}
