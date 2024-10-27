package cmd

import (
	"fmt"
	"github.com/keyCat/srcds-manager/config"
	"github.com/keyCat/srcds-manager/screen"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	rootCmd.AddCommand(idlermonitorCmd)
}

var idlermonitorCmd = &cobra.Command{
	Use:   "idlermonitor",
	Short: "Runs idler monitor daemon that starts and stops CPU idlers for running servers. Started automatically with start command.",
	Run: func(cmd *cobra.Command, args []string) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		shutdown := make(chan bool, 1)

		log.Println("Idler monitor started")

		go func() {
			sig := <-signals
			log.Println("Received signal:", sig)
			shutdown <- true
		}()

		for {
			select {
			case <-shutdown:
				log.Println("Shutting down")
				return
			default:
				time.Sleep(3 * time.Second)
				config.ReadOrThrow()
				for _, server := range config.Value.Servers {
					runIdlerForServer(server)
				}
			}
		}
	},
}

func runIdlerForServer(server config.Server) {
	idlerEnabled := false
	if server.Idler.Enabled != nil {
		idlerEnabled = *server.Idler.Enabled
	}
	idlerEnabled = idlerEnabled && screen.IsRunningForServer(server)
	idlerScreenName := screen.GetIdlerNameForServer(server)
	if idlerEnabled {
		if screen.IsRunning(idlerScreenName) {
			return
		}
		var cpusched string
		if server.Core != nil {
			// use taskset only if the core was set for the server
			cpusched = fmt.Sprintf("taskset -c %d ", *server.Core)
		}
		if server.Idler.Niceness != nil && *server.Idler.Niceness != 0 {
			cpusched += fmt.Sprintf("nice -%d ", *server.Idler.Niceness)
		}
		log.Printf("Starting idler for server #%02d (%s:%d)\n", server.Number, server.Ip, server.Port)
		exe, _ := os.Executable()
		err := screen.Start(idlerScreenName, fmt.Sprintf("%s%s idler", cpusched, exe))
		if err != nil {
			log.Printf("Error starting idler for server #%02d (%s:%d) %v\n", server.Number, server.Ip, server.Port, err)
		}
	} else {
		if screen.IsRunning(idlerScreenName) {
			log.Printf("Killing idler for server #%02d (%s:%d)\n", server.Number, server.Ip, server.Port)
			screen.Kill(idlerScreenName)
		}
		return
	}
}
