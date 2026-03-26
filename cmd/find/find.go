package find

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
)

var findRegex bool

var Cmd = &cobra.Command{
	Use:          "find [pattern]",
	Short:        "Find secret names in the cache",
	Long:         "Search cached secret names by substring or regex.",
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
			return fmt.Errorf("failed to load cache for profile '%s': %w. Hint: run 'lskv cache refresh' first", alias, err)
		}

		normalizedPattern := strings.ToLower(pattern)

		var regex *regexp.Regexp
		if findRegex {
			regex, err = regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern '%s': %w", pattern, err)
			}
		}

		for _, vault := range cacheData.Vaults {
			if !vault.Accessible || vault.Status != "success" {
				continue
			}

			for _, secret := range vault.Secrets {
				if findRegex {
					if regex.MatchString(secret.Name) {
						fmt.Printf("%s:%s\n", vault.Name, secret.Name)
					}
					continue
				}

				if strings.Contains(strings.ToLower(secret.Name), normalizedPattern) {
					fmt.Printf("%s:%s\n", vault.Name, secret.Name)
				}
			}
		}

		return nil
	},
}

func init() {
	Cmd.Flags().BoolVarP(&findRegex, "regex", "e", false, "Use regex pattern matching")
}
