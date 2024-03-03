package cmd

import (
	"github.com/keyCat/srcds-manager/cmd/parsers"
	"github.com/keyCat/srcds-manager/config"
	"github.com/keyCat/srcds-manager/lifecycle"
	"github.com/keyCat/srcds-manager/utils"
	"github.com/spf13/cobra"
	"log"
	"time"
)

func init() {
	restartCommand.Flags().BoolP("force", "f", false, "Force restart, even if server has players on it, or unreachable")
	rootCmd.AddCommand(restartCommand)
}

var restartCommand = &cobra.Command{
	Use:   "restart [1,2-4,6-] [-f]",
	Short: "Restarts servers. By default, does not attempt to restart unreachable servers or servers with players.",
	Args: func(cmd *cobra.Command, args []string) error {
		// validate server argument
		_, err := parsers.ParseAsServerNum(args)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		var serverNumbers []int
		force, _ := cmd.Flags().GetBool("force")
		serverNumbers, _ = parsers.ParseAsServerNum(args)
		// stop servers
		for _, server := range config.Value.Servers {
			if utils.SliceIncludesInt(serverNumbers, server.Number) {
				_, err := lifecycle.StopServer(server, force)
				if err != nil {
					log.Printf("%v\n", err)
				}
			}
		}
		// wait safe amount of time before starting servers
		time.Sleep(5 * time.Second)
		// start servers
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
