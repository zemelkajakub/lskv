package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileSwitchCmd = &cobra.Command{

	Use:   "switch [alias]",
	Short: "Change the active profile",
	Long: `Example:
lskv profile switch DEV
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		alias := args[0]

		if !profile.ExistsProfile(alias) {
			return fmt.Errorf("profile '%s' does not exist", alias)
		}

		if err := config.SetActiveProfile(alias); err != nil {
			return fmt.Errorf("failed to switch profile: %v", err)
		}

		fmt.Printf("Switched to profile '%s'\n", alias)

		return nil
	},
}
