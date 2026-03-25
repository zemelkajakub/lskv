package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileDeleteCmd = &cobra.Command{

	Use:   "delete [alias]",
	Short: "Delete specified profile or profiles",
	Long: `Example:
lskv profile delete DEV
lskv profile delete PROD TEST
	`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if argsLength := len(args); argsLength == 0 {
			return fmt.Errorf("no profile alias provided for deletion")
		}

		for _, alias := range args {

			if !profile.ExistsProfile(alias) {
				fmt.Printf("profile '%s' does not exist, skipping deletion\n", alias)
				continue
			}

			if err := profile.Delete(alias); err != nil {
				return fmt.Errorf("failed to delete profile '%s': %v", alias, err)
			}
		}
		return nil
	},
}
