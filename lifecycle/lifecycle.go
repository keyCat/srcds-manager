package lifecycle

import (
	"errors"
	"fmt"
	"github.com/keyCat/srcds-manager/a2s_query"
	"github.com/keyCat/srcds-manager/config"
	"github.com/keyCat/srcds-manager/screen"
	"github.com/keyCat/srcds-manager/utils"
	"log"
)

func StartServer(server config.Server, force bool) (bool, error) {
	var err error
	var started = false
	if force {
		utils.WriteStartScriptForServer(server)
		screen.KillForServer(server)
		err = screen.StartForServer(server)
		started = err == nil
	} else if screen.IsRunningForServer(server) {
		log.Printf("server %02d (%s:%d) is already running\n", server.Number, server.Ip, server.Port)
	} else {
		utils.WriteStartScriptForServer(server)
		err = screen.StartForServer(server)
		started = err == nil
	}
	if err != nil {
		return false, errors.New(fmt.Sprintf("error starting server #%02d (%s:%d, force=%t): %v", server.Number, server.Ip, server.Port, force, err))
	}
	if started {
		log.Printf("started server #%02d (%s:%d)\n", server.Number, server.Ip, server.Port)
	}

	return started, nil
}
func StopServer(server config.Server, force bool) (bool, error) {
	var stopped = false
	if force {
		screen.KillForServer(server)
		utils.DeleteStartScriptForServer(server)
		log.Printf("stopped server #%02d (%s:%d, force=%t)\n", server.Number, server.Ip, server.Port, force)
		stopped = true
	} else if screen.IsNotRunningForServer(server) {
		log.Printf("server #%02d (%s:%d) is not running", server.Number, server.Ip, server.Port)
		stopped = true
	} else {
		err, info := a2s_query.GetServerInfo(server)
		if err != nil {
			return false, errors.New(fmt.Sprintf("error stopping server #%02d (%s:%d): %v", server.Number, server.Ip, server.Port, err))
		}
		if info.Players > 0 {
			return false, errors.New(fmt.Sprintf("error stopping server #%02d (%s:%d): not empty", server.Number, server.Ip, server.Port))
		} else {
			screen.KillForServer(server)
			utils.DeleteStartScriptForServer(server)
			log.Printf("stopped server #%02d (%s:%d)", server.Number, server.Ip, server.Port)
			stopped = true
		}
	}

	return stopped, nil
}
