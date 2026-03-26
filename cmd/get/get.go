package get

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zemelkajakub/lskv/internal/keyvault"
)

var Cmd = &cobra.Command{

	Use:   "get [vault:secret|-]",
	Short: "Get secret value",
	Long:  "Get value of provided secret in vault:secret format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := keyvault.NewDataPlaneClient()
		if err != nil {
			return fmt.Errorf("failed to create Key Vault data-plane client: %w", err)
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
			return fmt.Errorf("failed to get secret '%s:%s': %w", vaultName, secretName, err)
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

	hadFailures := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		vaultName, secretName, err := parseVaultSecret(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\terror\t%v\n", line, err)
			hadFailures = true
			continue
		}

		value, status, err := client.GetSecretValue(cmd.Context(), vaultName, secretName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\terror\t%v\n", line, err)
			hadFailures = true
			continue
		}

		if status != "success" {
			fmt.Fprintf(os.Stderr, "%s\t%s\n", line, status)
			hadFailures = true
			continue
		}

		fmt.Printf("%s\t%s\n", line, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	if hadFailures {
		return fmt.Errorf("one or more secrets failed to retrieve")
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
