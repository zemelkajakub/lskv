package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileListCmd = &cobra.Command{

	Use:   "list",
	Short: "List available profiles",
	Long: `Example:
lskv profile list
	`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		profiles, err := profile.List()
		if err != nil {
			return fmt.Errorf("failed to list profiles: %v", err)
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles found.")
			return nil
		}

		fmt.Println("Available profiles:")
		for _, profile := range profiles {
			fmt.Println(profile)
		}
		return nil
	},
}
