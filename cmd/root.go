package cmd

import (
	"github.com/keyCat/srcds-manager/config"
	"github.com/spf13/cobra"
	"log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "servers",
	Short: "Manage srcds servers",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.ReadOrThrow()
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Long: `This command helps to manage the state of multiple running servers by facilitating tasks of starting, restarting and updating servers.

Version: 1.1.0  Author: 0x0c  Discord: @choree`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.Path, "config", "c", config.GetDefaultLocation(), "config file path")
}
