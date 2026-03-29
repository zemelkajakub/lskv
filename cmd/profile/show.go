package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileShowCmd = &cobra.Command{

	Use:          "show [alias]",
	Short:        "Show profile details",
	Long:         "Show the active profile or a specific profile alias.",
	SilenceUsage: true,
	Args:         cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var alias string
		var profileData *profile.Profile
		var err error

		switch argsLength := len(args); argsLength {
		case 0:

			alias, err = config.GetActiveProfile()
			if err != nil {
				return fmt.Errorf("failed to get active profile: %w", err)
			}
			if alias == "" {
				return fmt.Errorf("no profile specified and no active profile is set")
			}
			profileData, err = profile.Load(alias)
			if err != nil {
				return fmt.Errorf("failed to load profile '%s': %w", alias, err)
			}

		case 1:

			alias = args[0]

			if !profile.ExistsProfile(alias) {
				return fmt.Errorf("profile '%s' not found", alias)
			}
			profileData, err = profile.Load(alias)
			if err != nil {
				return fmt.Errorf("failed to load profile '%s': %w", alias, err)
			}
		}

		activeProfile, err := config.GetActiveProfile()
		if err != nil {
			return fmt.Errorf("failed to get active profile: %w", err)
		}

		fmt.Println("Profile")
		fmt.Println("-------")
		fmt.Printf("Alias: %s", profileData.Alias)
		if profileData.Alias == activeProfile {
			fmt.Print(" (active)")
		}
		fmt.Println()
		fmt.Printf("Subscription ID: %s\n", profileData.SubscriptionID)
		if profileData.Description != "" {
			fmt.Printf("Description: %s\n", profileData.Description)
		}
		if len(profileData.Vaults) == 0 {
			fmt.Println("Vault scope: all accessible vaults")
		} else {
			fmt.Printf("Vault scope: %d explicit vault(s)\n", len(profileData.Vaults))
			fmt.Println("Vaults:")
			for _, vault := range profileData.Vaults {
				fmt.Printf("  • %s\n", vault)
			}
		}

		return nil

	},
}
