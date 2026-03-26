package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileDeleteCmd = &cobra.Command{

	Use:          "delete [alias]",
	Short:        "Delete profiles",
	Long:         "Delete one or more profiles.",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deleted := make([]string, 0, len(args))
		skipped := make([]string, 0)

		if argsLength := len(args); argsLength == 0 {
			return fmt.Errorf("no profile aliases provided")
		}

		for _, alias := range args {

			if !profile.ExistsProfile(alias) {
				skipped = append(skipped, alias)
				continue
			}

			if err := profile.Delete(alias); err != nil {
				return fmt.Errorf("failed to delete profile '%s': %w", alias, err)
			}
			deleted = append(deleted, alias)
		}

		if len(deleted) > 0 {
			fmt.Println("Deleted profiles")
			fmt.Println("----------------")
			for _, alias := range deleted {
				fmt.Printf("• %s\n", alias)
			}
		}

		if len(skipped) > 0 {
			fmt.Println("Skipped profiles")
			fmt.Println("----------------")
			for _, alias := range skipped {
				fmt.Printf("• %s (not found)\n", alias)
			}
		}
		return nil
	},
}
