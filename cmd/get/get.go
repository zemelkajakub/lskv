package get

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/keyvault"
)

var Cmd = &cobra.Command{

	Use:          "get [vault:secret|-]",
	Short:        "Get a secret value",
	Long:         "Get a secret value directly from Azure Key Vault.",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := keyvault.NewDataPlaneClient()
		if err != nil {
			return fmt.Errorf("failed to create key vault client: %w", err)
		}

		target := strings.TrimSpace(args[0])
		if target == "-" {
			return runBatchGet(cmd, client)
		}

		vaultName, secretName, err := parseVaultSecret(target)
		if err != nil {
			return err
		}

		value, status, err := client.GetSecretValue(cmd.Context(), vaultName, secretName)
		if err != nil {
			return fmt.Errorf("failed to get secret '%s:%s': %s", vaultName, secretName, simplifyGetError(err))
		}

		if status != "success" {
			return fmt.Errorf("failed to get secret '%s:%s': %s", vaultName, secretName, status)
		}

		fmt.Println(value)
		return nil
	},
}

func runBatchGet(cmd *cobra.Command, client *keyvault.Client) error {
	scanner := bufio.NewScanner(os.Stdin)

	type getJob struct {
		line       string
		vaultName  string
		secretName string
	}

	type getResult struct {
		line   string
		value  string
		status string
		err    error
	}

	hadFailures := false
	jobsList := make([]getJob, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		vaultName, secretName, err := parseVaultSecret(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\t%v\n", line, err)
			hadFailures = true
			continue
		}

		jobsList = append(jobsList, getJob{
			line:       line,
			vaultName:  vaultName,
			secretName: secretName,
		})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	if len(jobsList) > 0 {
		workers := len(jobsList) / 2
		if workers < 1 {
			workers = 1
		}

		jobs := make(chan getJob, len(jobsList))
		results := make(chan getResult, len(jobsList))

		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for job := range jobs {
					value, status, err := client.GetSecretValue(cmd.Context(), job.vaultName, job.secretName)
					results <- getResult{
						line:   job.line,
						value:  value,
						status: status,
						err:    err,
					}
				}
			}()
		}

		for _, job := range jobsList {
			jobs <- job
		}
		close(jobs)

		for i := 0; i < len(jobsList); i++ {
			result := <-results
			if result.err != nil {
				fmt.Fprintf(os.Stderr, "%s\t%s\n", result.line, simplifyGetError(result.err))
				hadFailures = true
				continue
			}

			if result.status != "success" {
				fmt.Fprintf(os.Stderr, "%s\t%s\n", result.line, result.status)
				hadFailures = true
				continue
			}

			fmt.Printf("%s\t%s\n", result.line, result.value)
		}

		wg.Wait()
		close(results)
	}

	if hadFailures {
		return fmt.Errorf("failed to retrieve one or more secrets")
	}

	return nil
}

func parseVaultSecret(input string) (string, string, error) {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid secret reference '%s': expected format vault:secret", input)
	}

	vaultName := strings.TrimSpace(parts[0])
	secretName := strings.TrimSpace(parts[1])

	if vaultName == "" || secretName == "" {
		return "", "", fmt.Errorf("invalid secret reference '%s': both vault and secret must be non-empty", input)
	}

	return vaultName, secretName, nil
}

func simplifyGetError(err error) string {
	if err == nil {
		return "unknown error"
	}

	msg := err.Error()

	if strings.Contains(msg, "DefaultAzureCredential: failed to acquire a token") || strings.Contains(msg, "AADSTS700016") {
		return "authentication failed: this key vault is likely not related to your current tenant/login. Run 'az logout' and 'az login --tenant <tenant-id>'"
	}

	return msg
}
