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
	stopCommand.Flags().BoolP("force", "f", false, "Force stop, even if server has players on it, or unreachable")
	rootCmd.AddCommand(stopCommand)
}

var stopCommand = &cobra.Command{
	Use:   "stop [1,2-4,6-] [-f]",
	Short: "Restarts servers. By default, does not attempt to stop unreachable servers or servers with players.",
	Args: func(cmd *cobra.Command, args []string) error {
		// validate server argument
		_, err := parsers.ParseAsServerNum(args)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		//config.ReadOrThrow()
		var serverNumbers []int
		force, _ := cmd.Flags().GetBool("force")
		serverNumbers, _ = parsers.ParseAsServerNum(args)
		for _, server := range config.Value.Servers {
			if utils.SliceIncludesInt(serverNumbers, server.Number) {
				_, err := lifecycle.StopServer(server, force)
				if err != nil {
					log.Printf("%v\n", err)
				}
			}
		}
	},
}
