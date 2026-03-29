package find

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
)

var Cmd = &cobra.Command{
	Use:          "find [pattern]",
	Short:        "Find secret names in the cache",
	Long:         "Search cached secret names by substring.",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		pattern := strings.TrimSpace(args[0])
		if pattern == "" {
			return fmt.Errorf("search pattern cannot be empty")
		}

		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		cacheData, err := cache.LoadCache(alias)
		if err != nil {
			return fmt.Errorf("failed to load cache for profile '%s': %w.", alias, err)
		}

		normalizedPattern := strings.ToLower(pattern)

		for _, vault := range cacheData.Vaults {
			if !vault.Accessible || vault.Status != "success" {
				continue
			}

			for _, secret := range vault.Secrets {

				if strings.Contains(strings.ToLower(secret.Name), normalizedPattern) {
					fmt.Printf("%s:%s\n", vault.Name, secret.Name)
				}
			}
		}

		return nil
	},
}
