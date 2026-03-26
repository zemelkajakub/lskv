package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileSwitchCmd = &cobra.Command{

	Use:          "switch [alias]",
	Short:        "Switch the active profile",
	Long:         "Set the active profile.",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		alias := args[0]

		if !profile.ExistsProfile(alias) {
			return fmt.Errorf("profile '%s' not found", alias)
		}

		if err := config.SetActiveProfile(alias); err != nil {
			return fmt.Errorf("failed to switch profile: %w", err)
		}

		fmt.Println("Active profile updated")
		fmt.Printf("• %s\n", alias)

		return nil
	},
}
