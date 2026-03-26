package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zemelkajakub/lskv/internal/config"
	"gopkg.in/yaml.v3"
)

const (
	fileExt string = ".yaml"
)

// Profile struct represents a profile configuration
type Profile struct {
	Alias          string   `yaml:"alias"`
	SubscriptionID string   `yaml:"subscription_id"`
	Description    string   `yaml:"description"`
	Vaults         []string `yaml:"vaults,omitempty"`
}

// New create a new Profile instance
func New(alias string, subscriptionID string, description string, vaults []string) (*Profile, error) {

	if alias == "" {
		return nil, fmt.Errorf("alias cannot be empty")
	}
	if subscriptionID == "" {
		return nil, fmt.Errorf("subscription ID cannot be empty")
	}

	return &Profile{
		Alias:          alias,
		SubscriptionID: subscriptionID,
		Description:    description,
		Vaults:         vaults,
	}, nil
}

// Save saves the profile to a YAML file
func Save(p Profile) error {

	// Ensure profiles directory exists
	if err := config.EnsureProfilesDir(); err != nil {
		return fmt.Errorf("cannot create a profiles directory: %w", err)
	}

	path, err := GetProfilePath(p.Alias)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	// 0600 to keep private
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	return nil

}

// Load loads a profile by its alias and shows its details
func Load(alias string) (*Profile, error) {

	if !ExistsProfile(alias) {
		return nil, fmt.Errorf("%w: %s", ErrProfileNotFound, alias)
	}

	path, err := GetProfilePath(alias)
	if err != nil {
		return nil, err
	}

	// Read profile file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile file: %w", err)
	}

	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile data: %w", err)
	}

	return &p, nil
}

// Delete deletes the profile with the given alias
func Delete(alias string) error {

	if !ExistsProfile(alias) {
		return fmt.Errorf("%w: %s", ErrProfileNotFound, alias)
	}

	path, err := GetProfilePath(alias)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil

}

func List() ([]string, error) {

	var profiles []string

	profilesDir, err := config.GetProfilesDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(profilesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {

			fileName := file.Name()

			if strings.HasSuffix(fileName, fileExt) {
				alias := strings.TrimSuffix(fileName, fileExt)

				profiles = append(profiles, alias)
			}
		}

	}
	sort.Strings(profiles)
	return profiles, nil
}

// ExistsProfile checks if a profile with the given alias exists
func ExistsProfile(alias string) bool {
	path, err := GetProfilePath(alias)
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}

// GetProfilePath returns the file path for the given profile alias
func GetProfilePath(alias string) (string, error) {
	profilesDir, err := config.GetProfilesDir()
	if err != nil {
		return "", err
	}

	profileFileName := filepath.Join(profilesDir, alias+fileExt)
	return profileFileName, nil
}
