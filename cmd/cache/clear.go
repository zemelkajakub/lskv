package cache

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
)

var cacheClearCmd = &cobra.Command{
	Use:          "clear",
	Short:        "Clear the cache",
	Long:         "Clear the cache for the active profile.",
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		err = cache.ClearCache(alias)
		if err != nil {
			return err
		}
		fmt.Println("Cache cleared")
		fmt.Printf("Profile: %s\n", alias)

		return nil
	},
}
