package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zemelkajakub/lskv/cmd/cache"
	"github.com/zemelkajakub/lskv/cmd/find"
	"github.com/zemelkajakub/lskv/cmd/get"
	"github.com/zemelkajakub/lskv/cmd/list"
	"github.com/zemelkajakub/lskv/cmd/profile"
	"github.com/zemelkajakub/lskv/internal/config"
)

var rootCmd = &cobra.Command{

	Use:   "lskv",
	Short: "Azure Key Vault secret helper",
	Long:  "Find and retrieve Azure Key Vault secrets with profiles and local cache.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lskv - Azure Key Vault secret helper")
		fmt.Println()
		fmt.Println("Common commands:")
		fmt.Println("  lskv profile init DEV --subscription-id <subscription-id>")
		fmt.Println("  lskv cache refresh")
		fmt.Println("  lskv find traefik")
		fmt.Println("  lskv list secrets all")
		fmt.Println("  lskv get <vault:secret>")
		fmt.Println("  lskv find traefik | grep -E \"dev-\" | lskv get -")
		fmt.Println()
		fmt.Println("Use 'lskv --help' to see all commands.")
		fmt.Println()
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initializeDefaultConfig)
	rootCmd.AddCommand(profile.Cmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(cache.Cmd)
	rootCmd.AddCommand(find.Cmd)
	rootCmd.AddCommand(get.Cmd)

}

func initializeDefaultConfig() {

	// Get path to config file
	configFile, err := config.GetConfigFile()
	if err != nil {
		// Unable to get config file path; gently move forward
		return
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	viper.SetDefault("version", 1)

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}

}
