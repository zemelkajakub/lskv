package list

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "list",
	Short: "List vaults or secrets",
	Long:  "vaults/secrets based on the active profile",
}

func init() {
	Cmd.AddCommand(listVaultsCmd)
}
