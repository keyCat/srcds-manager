package main

import (
	"github.com/keyCat/srcds-manager/cmd"
	"github.com/keyCat/srcds-manager/utils"
)

func main() {
	utils.FatalIfCommandIsNotAvailable("taskset")
	utils.FatalIfCommandIsNotAvailable("nice")
	utils.FatalIfCommandIsNotAvailable("screen")
	cmd.Execute()
}
