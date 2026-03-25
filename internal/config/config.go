package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// GetActiveProfile returns the active profile alias from the config file
func GetActiveProfile() (string, error) {

	if viper.ConfigFileUsed() == "" {
		return "", ErrConfigFileNotFound
	}

	activeProfile := viper.GetString("active_profile")
	if activeProfile == "" {
		return "", ErrNoActiveProfile
	}

	return activeProfile, nil
}

// SetActiveProfile sets the active profile in the config file
func SetActiveProfile(alias string) error {

	if err := EnsureAppDir(); err != nil {
		return fmt.Errorf("cannot create app directory: %w", err)
	}

	viper.Set("active_profile", alias)

	configFile, err := GetConfigFile()
	if err != nil {
		return fmt.Errorf("error getting config file path: %w", err)
	}

	if err := viper.WriteConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.WriteConfigAs(configFile)
		}
		return fmt.Errorf("error writing config file: %v", err)
	}
	return nil
}
