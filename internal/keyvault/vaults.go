package keyvault

import (
	"context"
	"strings"
)

type VaultInfo struct {
	Name          string
	Location      string
	ResourceGroup string
	ID            string
}

func getResourceGroupFromID(vaultID string) string {

	var resourceGroup string
	parts := strings.Split(vaultID, "/")

	for i, part := range parts {
		if strings.EqualFold(part, "resourceGroups") && i+1 < len(parts) {
			resourceGroup = parts[i+1]
			return resourceGroup
		}
	}
	return "" // Return empty string if not found
}

func ListVaults(ctx context.Context, client *Client) ([]VaultInfo, error) {

	var results []VaultInfo

	switch len(client.Vaults) {

	case 0:
		pager := client.VaultsClient.NewListBySubscriptionPager(nil)

		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, vault := range page.Value {
				vaultInfo := VaultInfo{
					Name:          *vault.Name,
					Location:      *vault.Location,
					ResourceGroup: getResourceGroupFromID(*vault.ID),
					ID:            *vault.ID,
				}
				results = append(results, vaultInfo)
			}
		}
		return results, nil

	default:
		// Build a map for O(1) lookup (case-insensitive)
		targetVaults := make(map[string]bool)
		for _, vaultName := range client.Vaults {
			targetVaults[strings.ToLower(vaultName)] = true
		}

		pager := client.VaultsClient.NewListBySubscriptionPager(nil)

		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, vault := range page.Value {
				// Check if vault name is in the target list (case-insensitive)
				if targetVaults[strings.ToLower(*vault.Name)] {
					vaultInfo := VaultInfo{
						Name:          *vault.Name,
						Location:      *vault.Location,
						ResourceGroup: getResourceGroupFromID(*vault.ID),
						ID:            *vault.ID,
					}
					results = append(results, vaultInfo)
				}
			}
		}

		return results, nil
	}
}
