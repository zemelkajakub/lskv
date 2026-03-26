package cache

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var cacheRefreshCmd = &cobra.Command{

	Use:     "refresh",
	Aliases: []string{"refresh"},
	Short:   "Refresh cache of vaults included in the active profile",
	Long: `Example:
lskv cache refresh
	`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context()

		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		fmt.Printf("Refreshing cache for profile '%s'...\n", alias)

		// Load profile config based on alias
		p, err := profile.Load(alias)
		if err != nil {
			return err
		}

		newCache, err := cache.NewCache(ctx, p)
		if err != nil {
			return fmt.Errorf("failed to create cache: %v", err)
		}

		if err := cache.SaveCache(newCache); err != nil {
			return fmt.Errorf("failed to save cache: %v", err)
		}

		return nil

	},
}
