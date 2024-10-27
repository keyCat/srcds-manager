package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/keyCat/srcds-manager/a2s_query"
	"github.com/keyCat/srcds-manager/cmd/parsers"
	"github.com/keyCat/srcds-manager/config"
	"github.com/keyCat/srcds-manager/screen"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status [1,2-4,6-]",
	Short: "Queries servers and displays their status",
	Args: func(cmd *cobra.Command, args []string) error {
		// validate server argument
		_, err := parsers.ParseAsServerNum(args)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 2, Align: text.AlignCenter},
			{Number: 3, Align: text.AlignCenter},
		})
		t.AppendHeader(table.Row{"#", "Running", "Reachable", "Host", "Players", "Map", "Version"})
		for _, server := range config.Value.Servers {
			info, err := a2s_query.GetServerInfo(server)
			reachableLabel := text.FgRed.Sprintf("✗")
			runningLabel := text.FgRed.Sprintf("✗")
			hostname := "<UNREACHABLE>"
			currentMap := "<none>"
			version := "0.0.0.0"
			players := 0
			maxPlayers := 0
			if err == nil {
				reachableLabel = text.FgGreen.Sprintf("✓")
				hostname = info.Name
				currentMap = info.Map
				version = info.Version
				players = int(info.Players)
				maxPlayers = int(info.MaxPlayers)
			}
			if screen.IsRunningForServer(server) {
				runningLabel = text.FgGreen.Sprintf(screen.GetNameForServer(server))
			}
			t.AppendSeparator()
			t.AppendRow([]interface{}{server.Number, runningLabel, reachableLabel, fmt.Sprintf("%s\n%s:%d", text.Bold.Sprintf(hostname), server.Ip, server.Port), fmt.Sprintf("%d/%d", players, maxPlayers), currentMap, version})
		}

		t.Render()
	},
}
