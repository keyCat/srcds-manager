package a2s_query

import (
	"fmt"
	"github.com/keyCat/srcds-manager/config"
	"github.com/rumblefrog/go-a2s"
	"time"
)

func GetServerInfo(server config.Server) (*a2s.ServerInfo, error) {
	var client, err = a2s.NewClient(
		fmt.Sprintf("%s:%d", server.Ip, server.Port),
		a2s.TimeoutOption(time.Second*1),
	)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.QueryInfo()
}
