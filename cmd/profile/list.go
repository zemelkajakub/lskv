package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileListCmd = &cobra.Command{

	Use:          "list",
	Short:        "List configured profiles",
	Long:         "List all saved profile aliases.",
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		profiles, err := profile.List()
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		activeProfile, err := config.GetActiveProfile()
		if err != nil {
			return fmt.Errorf("failed to get active profile: %w", err)
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles configured.")
			return nil
		}

		fmt.Println("Profiles")
		fmt.Println("--------")
		for _, profileAlias := range profiles {
			if profileAlias == activeProfile {
				fmt.Printf("• %s (active)\n", profileAlias)
				continue
			}
			fmt.Printf("• %s\n", profileAlias)
		}
		return nil
	},
}
