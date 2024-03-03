package cmd

import (
	"github.com/keyCat/srcds-manager/cmd/parsers"
	"github.com/keyCat/srcds-manager/config"
	"github.com/keyCat/srcds-manager/lifecycle"
	"github.com/keyCat/srcds-manager/utils"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	startCommand.Flags().BoolP("force", "f", false, "Force start, even if server is already started (and possibly has players on it)")
	rootCmd.AddCommand(startCommand)
}

var startCommand = &cobra.Command{
	Use:   "start [1,2-4,6-] [-f]",
	Short: "Starts stopped servers. No action is taken, if the server is already running",
	Args: func(cmd *cobra.Command, args []string) error {
		// validate server argument
		_, err := parsers.ParseAsServerNum(args)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		var serverNumbers []int
		force, _ := cmd.Flags().GetBool("force")
		serverNumbers, _ = parsers.ParseAsServerNum(args)
		for _, server := range config.Value.Servers {
			if utils.SliceIncludesInt(serverNumbers, server.Number) {
				_, err := lifecycle.StartServer(server, force)
				if err != nil {
					log.Printf("%v\n", err)
				}
			}
		}
	},
}
