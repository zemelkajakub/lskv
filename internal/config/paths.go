package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	configFile    string = "config.yaml"
	AppDir        string = ".lskv"
	profileSubDir string = "profiles"
)

// GetAppDir returns the application directory name
func GetAppDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	appDirPath := filepath.Join(homeDir, AppDir)

	return appDirPath, nil
}

// GetProfilesDir returns the profiles directory path
func GetProfilesDir() (string, error) {
	appDirPath, err := GetAppDir()
	if err != nil {
		return "", err
	}
	profilesDir := filepath.Join(appDirPath, profileSubDir)

	return profilesDir, nil
}

// GetConfigFile returns path to the main configuration file
func GetConfigFile() (string, error) {
	configDir, err := GetAppDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(configDir, configFile)

	return configPath, nil
}

// EnsureAppDir ensures that application directory exists and returns its path
func EnsureAppDir() error {
	appDirPath, err := GetAppDir()
	if err != nil {
		return err
	}

	// Ensure .lskv directory exists
	err = os.MkdirAll(appDirPath, 0o700)
	if err != nil {
		return fmt.Errorf("failed to create '%s' directory: %w", appDirPath, err)
	}
	return nil
}

// EnsureProfileDir ensures that the profiles directory exists and returns its path
func EnsureProfilesDir() error {
	profilesDir, err := GetProfilesDir()
	if err != nil {
		return err
	}

	// Ensure profiles directory exists
	err = os.MkdirAll(profilesDir, 0o700)
	if err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	return nil
}
