package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/zemelkajakub/lskv/internal/keyvault"
	"github.com/zemelkajakub/lskv/internal/profile"
)

type Cache struct {
	ProfileAlias   string     `json:"profile"`
	SubscriptionID string     `json:"subscription_id"`
	LastRefresh    time.Time  `json:"last_refresh"`
	Vaults         []Vault    `json:"vaults"`
	Statistics     Statistics `json:"statistics"`
}

type Vault struct {
	Name          string            `json:"name"`
	Status        string            `json:"status"`
	ResourceGroup string            `json:"resource_group"`
	Location      string            `json:"location"`
	Accessible    bool              `json:"accessible"`
	SecretsCount  int               `json:"secrets_count"`
	Secrets       []keyvault.Secret `json:"secrets,omitempty"`
	LastRefreshed time.Time         `json:"last_refreshed"`
}

type Statistics struct {
	TotalVaults      int `json:"total_vaults"`
	TotalSecrets     int `json:"total_secrets"`
	AccessibleVaults int `json:"accessible_vaults"`
}

func NewCache(ctx context.Context, p *profile.Profile) (*Cache, error) {

	// Create client for subscription from active profile
	client, err := keyvault.NewClient(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create Key Vault client: %v", err)
	}

	var vaults []keyvault.VaultInfo
	if len(p.Vaults) > 0 {
		vaults = make([]keyvault.VaultInfo, 0, len(p.Vaults))
		for _, vaultName := range p.Vaults {
			vaults = append(vaults, keyvault.VaultInfo{Name: vaultName})
		}
	} else {
		vaults, err = keyvault.ListVaults(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to list Key Vaults: %v", err)
		}
	}

	cacheData := Cache{
		ProfileAlias:   p.Alias,
		SubscriptionID: p.SubscriptionID,
		LastRefresh:    time.Now(),
		Vaults:         make([]Vault, 0, len(vaults)),
	}

	if len(vaults) == 0 {
		cacheData.Statistics = Statistics{
			TotalVaults:      0,
			TotalSecrets:     0,
			AccessibleVaults: 0,
		}
		return &cacheData, nil
	}

	type vaultJob struct {
		vault keyvault.VaultInfo
	}

	type vaultResult struct {
		vaultEntry Vault
		err        error
	}

	workers := len(vaults) / 2
	if workers < 1 {
		workers = 1
	}

	totalSecrets := 0
	accessibleVaults := 0

	jobs := make(chan vaultJob, len(vaults))
	results := make(chan vaultResult, len(vaults))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for job := range jobs {
				vault := job.vault

				vaultEntry := Vault{
					Name:          vault.Name,
					ResourceGroup: vault.ResourceGroup,
					Location:      vault.Location,
					LastRefreshed: time.Now(),
				}

				secrets, status, err := client.ListSecrets(ctx, vault.Name)

				vaultEntry.Status = status
				if status == "success" {
					vaultEntry.Accessible = true
					vaultEntry.Secrets = secrets
					vaultEntry.SecretsCount = len(secrets)
				} else {
					vaultEntry.Accessible = false
				}

				results <- vaultResult{vaultEntry: vaultEntry, err: err}
			}
		}()
	}

	for _, vault := range vaults {
		jobs <- vaultJob{vault: vault}
	}
	close(jobs)

	for i := 0; i < len(vaults); i++ {
		result := <-results

		cacheData.Vaults = append(cacheData.Vaults, result.vaultEntry)

		if result.vaultEntry.Status == "success" {
			accessibleVaults++
			totalSecrets += result.vaultEntry.SecretsCount
		}

		// Log error details for unexpected failures
		if result.err != nil {
			fmt.Printf("  ✗ %-30s %s: %v\n", result.vaultEntry.Name, result.vaultEntry.Status, result.err)
		} else {
			statusSymbol := "✗"
			if result.vaultEntry.Status == "success" {
				statusSymbol = "✓"
			}
			fmt.Printf("  %s %-30s %s\n", statusSymbol, result.vaultEntry.Name, result.vaultEntry.Status)
		}
	}

	wg.Wait()
	close(results)

	cacheData.Statistics = Statistics{
		TotalVaults:      len(vaults),
		TotalSecrets:     totalSecrets,
		AccessibleVaults: accessibleVaults,
	}
	return &cacheData, nil
}

func SaveCache(cache *Cache) error {

	alias := cache.ProfileAlias

	if err := EnsureCacheDir(); err != nil {
		return fmt.Errorf("cannot create cache directory: %w", err)
	}

	cacheDirPath, err := GetCacheDir()
	if err != nil {
		return fmt.Errorf("cannot get cache directory path: %w", err)
	}

	cacheName := fmt.Sprintf("%s.json", alias)
	cacheFile := path.Join(cacheDirPath, cacheName)

	cacheJson, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("Failed marshaling cache data to JSON: %w", err)
	}

	// 0600 to keep private
	if err := os.WriteFile(cacheFile, cacheJson, 0o600); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	return nil
}

func LoadCache(alias string) (*Cache, error) {

	if alias == "" {
		return nil, fmt.Errorf("profile alias cannot be empty")
	}

	cacheDirPath, err := GetCacheDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get cache directory path: %w", err)
	}

	cacheName := fmt.Sprintf("%s.json", alias)
	cacheFile := path.Join(cacheDirPath, cacheName)

	cacheJSON, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache for profile '%s' not found: %w", alias, err)
		}
		return nil, fmt.Errorf("failed to read cache file '%s': %w", cacheFile, err)
	}

	var loadedCache Cache
	if err := json.Unmarshal(cacheJSON, &loadedCache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache file '%s': %w", cacheFile, err)
	}

	return &loadedCache, nil
}

func ClearCache(alias string) error {

	cacheDirPath, err := GetCacheDir()
	if err != nil {
		return fmt.Errorf("cannot get cache directory path: %w", err)
	}

	cacheName := fmt.Sprintf("%s.json", alias)
	cacheFile := path.Join(cacheDirPath, cacheName)

	if err := os.Remove(cacheFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache for profile '%s' not found: %w", alias, err)
		}
		return fmt.Errorf("failed to delete cache file '%s': %w", cacheFile, err)
	}

	return nil
}
