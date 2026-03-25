package profile

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{

	Use:   "profile",
	Short: "Manage profiles",
	Long:  "init/list/show/switch/delete profiles",
}

func init() {
	Cmd.AddCommand(profileInitCmd)
	Cmd.AddCommand(profileShowCmd)
	Cmd.AddCommand(profileDeleteCmd)
	Cmd.AddCommand(profileListCmd)
	Cmd.AddCommand(profileSwitchCmd)
}
