package core

import "github.com/codecrafters-io/redis-starter-go/core/replication"

type ServerInfo struct {
	Replication replication.InfoSection
}

func NewServerInfo() ServerInfo {
	return ServerInfo{
		Replication: replication.NewInfoSection(),
	}
}
