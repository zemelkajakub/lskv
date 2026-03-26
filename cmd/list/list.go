package list

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "list",
	Short: "List vaults or secrets",
	Long:  "List accessible vaults or cached secrets.",
}

func init() {
	Cmd.AddCommand(listVaultsCmd)
	Cmd.AddCommand(listSecretsCmd)
}
