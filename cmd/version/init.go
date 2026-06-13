package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "unknown"

var Cmd = &cobra.Command{
	Use:          "version",
	Short:        "Print version information",
	Long:         "Print the version of lskv.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("lskv version: %s\n", version)
		return nil
	},
}
