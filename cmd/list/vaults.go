package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/keyvault"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var listVaultsCmd = &cobra.Command{

	Use:     "vaults",
	Aliases: []string{"vaults", "v"},
	Short:   "List vaults included in the active profile",
	Long: `Example:
lskv list vaults
	`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context()

		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		// Load profile config based on alias
		p, err := profile.Load(alias)
		if err != nil {
			return err
		}

		// Create client for subscription from active profile
		client, err := keyvault.NewClient(p)
		if err != nil {
			return fmt.Errorf("failed to create Key Vault client: %v", err)
		}

		vaults, err := keyvault.ListVaults(ctx, client)

		for _, vault := range vaults {
			fmt.Printf("\nName: %s\nLocation: %s\nResource Group: %s\nID: %s\n",
				vault.Name, vault.Location, vault.ResourceGroup, vault.ID)
		}

		return nil
	},
}
