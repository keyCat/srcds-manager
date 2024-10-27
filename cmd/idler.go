package cmd

import (
	"github.com/spf13/cobra"
	"runtime"
)

func init() {
	rootCmd.AddCommand(idler)
}

var idler = &cobra.Command{
	Use:   "idler",
	Short: "Runs CPU idler (called by idlermonitor with necessary params, no need to run manually)",
	Run: func(cmd *cobra.Command, args []string) {
		runtime.GOMAXPROCS(1)
		for true {
			// infinite loop to stress CPU
		}
	},
}
