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

	Use:   "init [alias]",
	Short: "Initialize a profile",
	Long: `Example:
lskv profile init DEV --subscription-id 0000-1111 --description "Development Subscription"
lskv profile init PROD --subscription-id 2222-3333 --description "Production Subscription" --vaults kv-prod1-kv, kv-prod2-kv
	`,
	Args: cobra.ExactArgs(1),
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
			return fmt.Errorf("failed to initialize profile: %v", err)
		}

		// Save profile to file
		if err := profile.Save(*initializedProfile); err != nil {
			return fmt.Errorf("failed to save profile: %v", err)
		}

		if err := config.SetActiveProfile(initializedProfile.Alias); err != nil {
			return fmt.Errorf("failed to set active profile: %v", err)
		}

		fmt.Printf("Profile '%s' initialized successfully.\n", initializedProfile.Alias)
		fmt.Printf("Subscription ID: %s\n", initializedProfile.SubscriptionID)
		fmt.Printf("Description: %s\n", initializedProfile.Description)
		if len(initializedProfile.Vaults) > 0 {
			fmt.Printf("Vaults: %s\n", strings.Join(initializedProfile.Vaults, ", "))
		} else {
			fmt.Printf("Vaults: (all vaults will be used)\n")
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
