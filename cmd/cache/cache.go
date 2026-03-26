package cache

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "cache",
	Short: "Manage the local secret cache",
	Long:  "Refresh, inspect, or clear the local cache.",
}

func init() {
	Cmd.AddCommand(cacheRefreshCmd)
	Cmd.AddCommand(cacheStatusCmd)
	Cmd.AddCommand(cacheClearCmd)
}
