package cache

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var cacheRefreshCmd = &cobra.Command{

	Use:          "refresh",
	Short:        "Refresh cache",
	Long:         "Refresh the cache for the active profile.",
	SilenceUsage: true,
	Args:         cobra.NoArgs,
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
			return fmt.Errorf("failed to refresh cache: %w", err)
		}

		if err := cache.SaveCache(newCache); err != nil {
			return fmt.Errorf("failed to save cache: %w", err)
		}

		fmt.Println("Cache refreshed")
		fmt.Println("---------------")
		fmt.Printf("Profile: %s\n", newCache.ProfileAlias)
		fmt.Printf("Vaults: %d total, %d accessible\n", newCache.Statistics.TotalVaults, newCache.Statistics.AccessibleVaults)
		fmt.Printf("Secrets: %d total\n", newCache.Statistics.TotalSecrets)

		return nil

	},
}
