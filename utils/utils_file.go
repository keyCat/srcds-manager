package utils

import (
	"fmt"
	"github.com/keyCat/srcds-manager/config"
	"log"
	"os"
	"path/filepath"
)

func WriteStartScriptForServer(server config.Server) {
	var cpusched string
	if server.Core != nil {
		// use taskset only if the core was set for the server
		cpusched = fmt.Sprintf("taskset -c %d ", *server.Core)
	}
	if *server.Niceness != 0 {
		cpusched = fmt.Sprintf("%snice -%d ", cpusched, *server.Niceness)
	}
	var contents = fmt.Sprintf(`#!/usr/bin/env bash

###########################################################################################
##                              DO NOT EDIT THIS FILE MANUALLY!
##
## This file was generated by srcds management utility and should not be edited by hand.
## If you want to edit launch params, check out config file at location:
## %s
###########################################################################################

while true
do
	%s%s -ip %s -port %d +map %s +sn_host_num %d +sv_logsdir logs%02d %s %s
	sleep 2
done
`,
		config.Path,
		cpusched,
		filepath.Join(server.Path, "/srcds_run"),
		server.Ip,
		server.Port,
		server.Map,
		server.Number,
		server.Number,
		server.Params,
		server.ParamsAdd,
	)

	var fp = GetStartScriptPathForServer(server)
	err := os.WriteFile(fp, []byte(contents), 0744)
	if err != nil {
		log.Fatalf("error writing server script: %v", err)
	}
}

func DeleteStartScriptForServer(server config.Server) {
	os.Remove(GetStartScriptPathForServer(server))
}

func GetStartScriptPathForServer(server config.Server) string {
	return filepath.Join(server.Path, fmt.Sprintf("/srcds%02d.sh", server.Number))
}
