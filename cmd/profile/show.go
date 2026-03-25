package profile

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var profileShowCmd = &cobra.Command{

	Use:   "show [alias]",
	Short: "Show details of current or specified profile",
	Long: `Example:
lskv profile show
lskv profile show DEV
	`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {

		switch argsLength := len(args); argsLength {
		case 0:

			activeProfile, err := config.GetActiveProfile()
			if activeProfile == "" {
				return fmt.Errorf("no profile alias provided and no current profile set in config")
			}

			profileConfigPath, err := profile.GetProfilePath(activeProfile)
			if err != nil {
				return fmt.Errorf("error loading profile '%s': %v", activeProfile, err)
			}

			data, err := os.ReadFile(profileConfigPath)
			if err != nil {
				return fmt.Errorf("error reading profile '%s': %v", activeProfile, err)
			}

			fmt.Printf(string(data))

		case 1:

			alias := args[0]

			if !profile.ExistsProfile(alias) {
				return fmt.Errorf("profile '%s' does not exist", alias)
			}

			profileConfigPath, err := profile.GetProfilePath(alias)
			if err != nil {
				return fmt.Errorf("error loading profile '%s': %v", alias, err)
			}

			data, err := os.ReadFile(profileConfigPath)
			if err != nil {
				return fmt.Errorf("error reading profile '%s': %v", alias, err)
			}

			fmt.Printf(string(data))
		}

		return nil

	},
}
