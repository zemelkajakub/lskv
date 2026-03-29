package list

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/config"
	"github.com/zemelkajakub/lskv/internal/keyvault"
	"github.com/zemelkajakub/lskv/internal/profile"
)

var listVaultsCmd = &cobra.Command{

	Use:          "vaults",
	Aliases:      []string{"vault", "v"},
	Short:        "List accessible vaults",
	Long:         "List vaults that are accessible for secret operations.",
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := cmd.Context()

		alias, err := config.GetActiveProfile()
		if err != nil {
			return err
		}

		// Load profile config based on alias
		p, err := profile.Load(alias)
		if err != nil {
			return err
		}

		// Create client for subscription from active profile
		client, err := keyvault.NewClient(p)
		if err != nil {
			return fmt.Errorf("failed to create key vault client: %w", err)
		}

		var vaults []keyvault.VaultInfo
		if len(p.Vaults) > 0 {
			vaults = make([]keyvault.VaultInfo, 0, len(p.Vaults))
			for _, name := range p.Vaults {
				vaults = append(vaults, keyvault.VaultInfo{Name: name})
			}
		} else {
			vaults, err = keyvault.ListVaults(ctx, client)
			if err != nil {
				return fmt.Errorf("failed to list key vaults: %w", err)
			}
		}

		accessibleCount := 0
		totalCount := len(vaults)
		if totalCount == 0 {
			fmt.Println("Accessible vaults")
			fmt.Println("-----------------")
			fmt.Println("No accessible vaults found.")
			fmt.Fprintf(os.Stderr, "accessible vaults: 0/0\n")
			return nil
		}

		type vaultJob struct {
			name string
		}

		type vaultResult struct {
			name   string
			status string
			err    error
		}

		workers := min(len(vaults)/2, 40)
		if workers < 1 {
			workers = 1
		}

		jobs := make(chan vaultJob, totalCount)
		results := make(chan vaultResult, totalCount)
		accessibleVaults := make([]string, 0, totalCount)

		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for job := range jobs {
					_, status, err := client.ListSecrets(ctx, job.name)

					results <- vaultResult{
						name:   job.name,
						status: status,
						err:    err,
					}
				}
			}()
		}

		for _, vault := range vaults {
			jobs <- vaultJob{name: vault.Name}
		}
		close(jobs)

		for i := 0; i < totalCount; i++ {
			result := <-results
			if result.err != nil {
				fmt.Fprintf(os.Stderr, "%s\terror\t%v\n", result.name, result.err)
				continue
			}

			if result.status != "success" {
				continue
			}

			accessibleVaults = append(accessibleVaults, result.name)
			accessibleCount++
		}

		wg.Wait()
		sort.Strings(accessibleVaults)

		fmt.Println("Accessible vaults")
		fmt.Println("-----------------")
		for _, vaultName := range accessibleVaults {
			fmt.Printf("• %s\n", vaultName)
		}

		fmt.Fprintf(os.Stderr, "accessible vaults: %d/%d\n", accessibleCount, totalCount)

		return nil
	},
}
