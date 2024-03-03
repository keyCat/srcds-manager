package a2s_query

import (
	"fmt"
	"github.com/keyCat/srcds-manager/config"
	"github.com/rumblefrog/go-a2s"
	"time"
)

func GetServerInfo(server config.Server) (error, *a2s.ServerInfo) {
	var client, err = a2s.NewClient(
		fmt.Sprintf("%s:%d", server.Ip, server.Port),
		a2s.TimeoutOption(time.Second*1),
	)
	if err != nil {
		return err, nil
	}
	defer client.Close()

	info, err := client.QueryInfo()
	return err, info
}
