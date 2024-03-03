package screen

import (
	"bytes"
	"fmt"
	"github.com/keyCat/srcds-manager/config"
	"github.com/keyCat/srcds-manager/utils"
	"log"
	"os/exec"
	"time"
)

func GetNameForServer(server config.Server) string {
	return fmt.Sprintf("%s%02d", config.Value.Project, server.Number)
}

func StartForServer(server config.Server) error {
	cmd := exec.Command("screen", "-U", "-m", "-d", "-S", GetNameForServer(server), utils.GetStartScriptPathForServer(server))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return err
}

// IsRunningForServer checks if screen is running for server config
func IsRunningForServer(server config.Server) bool {
	return IsRunning(GetNameForServer(server))
}

// IsNotRunningForServer checks if screen is not running for server config
func IsNotRunningForServer(server config.Server) bool {
	return !IsRunningForServer(server)
}

// IsRunning checks if screen with given name is running
// screen -list | grep %name%
func IsRunning(name string) bool {
	cmd := exec.Command("screen", "-list")
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()
	grep := exec.Command("grep", fmt.Sprintf("%s", name))
	grep.Stdin = stdout
	cmd.Start()
	out, _ := grep.Output()
	cmd.Wait()

	return string(out) != ""
}

// KillForServer kills screen for server config
func KillForServer(server config.Server) {
	Kill(GetNameForServer(server))
}

// Kill sends a quite command to a screen with the given name
// screen -r %name% -X quit
func Kill(name string) {
	if IsRunning(name) {
		cmd := exec.Command("screen", "-S", fmt.Sprintf("%s", name), "-X", "quit")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("could not kill the screen \"%s\": %v — %v", name, err, stderr.String())
		}

		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			if !IsRunning(name) {
				break
			}
		}
	}
}
