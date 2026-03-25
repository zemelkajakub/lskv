package cache

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
)

var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache status for active profile",
	Long: `Example:
lskv cache status
	`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		cacheData, err := cache.LoadCache(alias)
		if err != nil {
			return err
		}

		cacheDirPath, err := cache.GetCacheDir()
		if err != nil {
			return fmt.Errorf("cannot get cache directory path: %w", err)
		}

		cacheFile := path.Join(cacheDirPath, fmt.Sprintf("%s.json", alias))
		fileInfo, err := os.Stat(cacheFile)
		if err != nil {
			return fmt.Errorf("failed to read cache file metadata '%s': %w", cacheFile, err)
		}

		age := time.Since(cacheData.LastRefresh).Round(time.Second)

		fmt.Printf("Profile: %s\n", cacheData.ProfileAlias)
		fmt.Printf("Subscription ID: %s\n", cacheData.SubscriptionID)
		fmt.Printf("Last refresh: %s\n", cacheData.LastRefresh.Format(time.RFC3339))
		fmt.Printf("Cache age: %s\n", age)
		fmt.Printf("Cache file: %s\n", cacheFile)
		fmt.Printf("Cache size: %d bytes\n", fileInfo.Size())
		fmt.Printf("Vaults: %d total, %d accessible\n", cacheData.Statistics.TotalVaults, cacheData.Statistics.AccessibleVaults)
		fmt.Printf("Secrets: %d total\n", cacheData.Statistics.TotalSecrets)

		return nil
	},
}
