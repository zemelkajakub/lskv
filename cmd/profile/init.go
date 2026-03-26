package profile

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var (
	//Flags
	initSubscriptionID string
	initDescription    string
	initVaults         string
)

var profileInitCmd = &cobra.Command{

	Use:          "init [alias]",
	Short:        "Create a profile",
	Long:         "Create a profile and make it active.",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Get alias as a first argument
		alias := args[0]

		var vaults []string

		if strings.TrimSpace(initVaults) != "" {
			for _, v := range strings.Split(initVaults, ",") {
				s := strings.TrimSpace(v)
				if s != "" {
					vaults = append(vaults, s)
				}
			}
		}

		initializedProfile, err := profile.New(alias, initSubscriptionID, initDescription, vaults)
		if err != nil {
			return fmt.Errorf("failed to initialize profile: %w", err)
		}

		// Save profile to file
		if err := profile.Save(*initializedProfile); err != nil {
			return fmt.Errorf("failed to save profile: %w", err)
		}

		if err := config.SetActiveProfile(initializedProfile.Alias); err != nil {
			return fmt.Errorf("failed to set active profile: %w", err)
		}

		fmt.Println("Profile initialized")
		fmt.Println("-------------------")
		fmt.Printf("Alias: %s (active)\n", initializedProfile.Alias)
		fmt.Printf("Subscription ID: %s\n", initializedProfile.SubscriptionID)
		if initializedProfile.Description != "" {
			fmt.Printf("Description: %s\n", initializedProfile.Description)
		}
		if len(initializedProfile.Vaults) > 0 {
			fmt.Printf("Vault scope: %d explicit vault(s)\n", len(initializedProfile.Vaults))
			fmt.Printf("Vaults: %s\n", strings.Join(initializedProfile.Vaults, ", "))
		} else {
			fmt.Println("Vault scope: all accessible vaults")
		}

		return nil
	},
}

func init() {
	profileInitCmd.Flags().StringVar(&initSubscriptionID, "subscription-id", "", "Azure Subscription ID (required)")
	profileInitCmd.Flags().StringVar(&initDescription, "description", "", "Description for the profile")
	profileInitCmd.Flags().StringVar(&initVaults, "vaults", "", "Comma-separated list of Key Vault names. If omitted, profile will cover all Key Vaults.")
	profileInitCmd.MarkFlagsOneRequired("subscription-id", "description")

}
