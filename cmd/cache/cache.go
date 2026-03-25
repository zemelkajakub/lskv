package cache

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "cache",
	Short: "Manage local cache of vaults and secrets",
	Long:  "refresh/status/clear cache state",
}

func init() {
	Cmd.AddCommand(cacheRefreshCmd)
	Cmd.AddCommand(cacheStatusCmd)
	Cmd.AddCommand(cacheClearCmd)
}
