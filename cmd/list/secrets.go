package list

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/cache"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/keyvault"
)

var listSecretsCmd = &cobra.Command{
	Use:          "secrets [all|vault]",
	Short:        "List secrets as vault:secret",
	Long:         "List cached secrets for all vaults or list secrets for one vault.",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		target := strings.TrimSpace(args[0])
		if target == "" {
			return fmt.Errorf("vault name or 'all' is required")
		}

		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		if strings.EqualFold(target, "all") {
			cacheData, err := cache.LoadCache(alias)
			if err != nil {
				return fmt.Errorf("failed to load cache for profile '%s': %w.", alias, err)
			}

			var lines []string
			for _, vault := range cacheData.Vaults {
				if !vault.Accessible || vault.Status != "success" {
					continue
				}

				for _, secret := range vault.Secrets {
					lines = append(lines, fmt.Sprintf("%s:%s", vault.Name, secret.Name))
				}
			}
			sort.Strings(lines)
			for _, line := range lines {
				fmt.Println(line)
			}
			return nil
		}

		cacheData, err := cache.LoadCache(alias)
		if err == nil {
			for _, vault := range cacheData.Vaults {
				if strings.EqualFold(vault.Name, target) {
					if !vault.Accessible || vault.Status != "success" {
						return fmt.Errorf("vault '%s' is not accessible: %s", target, vault.Status)
					}

					for _, secret := range vault.Secrets {
						fmt.Printf("%s:%s\n", vault.Name, secret.Name)
					}
					return nil
				}
			}
		}

		client, err := keyvault.NewDataPlaneClient()
		if err != nil {
			return fmt.Errorf("failed to create key vault client: %w", err)
		}

		secrets, status, err := client.ListSecrets(ctx, target)
		if err != nil {
			return fmt.Errorf("failed to list secrets for vault '%s': %w", target, err)
		}
		if status != "success" {
			return fmt.Errorf("failed to list secrets for vault '%s': %s", target, status)
		}

		for _, secret := range secrets {
			fmt.Printf("%s:%s\n", target, secret.Name)
		}

		return nil
	},
}
