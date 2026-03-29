package profile

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "profile",
	Short: "Manage Azure profiles",
	Long:  "Create, list, show, switch, and delete profiles.",
}

func init() {
	Cmd.AddCommand(profileInitCmd)
	Cmd.AddCommand(profileShowCmd)
	Cmd.AddCommand(profileDeleteCmd)
	Cmd.AddCommand(profileListCmd)
	Cmd.AddCommand(profileSwitchCmd)
}
