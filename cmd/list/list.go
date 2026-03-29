package list

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "list",
	Short: "List vaults and secrets",
	Long:  "List accessible vaults or secrets in vault:secret format.",
}

func init() {
	Cmd.AddCommand(listVaultsCmd)
	Cmd.AddCommand(listSecretsCmd)
}
