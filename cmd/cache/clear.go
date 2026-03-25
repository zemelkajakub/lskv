package cache

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
)

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache for active profile",
	Long: `Example:
lskv cache clear
	`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		err = cache.ClearCache(alias)
		if err != nil {
			return err
		}
		fmt.Printf("Cache cleared for profile '%s'\n", alias)

		return nil
	},
}
