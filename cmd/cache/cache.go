package cache

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "cache",
	Short: "Manage local cache",
	Long:  "Refresh, view status, or clear the local secret cache.",
}

func init() {
	Cmd.AddCommand(cacheRefreshCmd)
	Cmd.AddCommand(cacheStatusCmd)
	Cmd.AddCommand(cacheClearCmd)
}
