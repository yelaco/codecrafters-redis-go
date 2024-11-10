package core

import (
	"bytes"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/config"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type Core struct {
	Conn net.Conn
	Dict map[string]string
}

func NewCore(conn net.Conn) *Core {
	c := &Core{
		Conn: conn,
		Dict: make(map[string]string),
	}
	return c
}

func (c *Core) HandleCommand(command resp.RespData) (resp.RespData, error) {
	if command.Type != resp.Array {
		return resp.RespData{}, ErrInvalidCommandRespDataType
	}
	dataArr := command.Value.(resp.RespDataArray)

	cmd := make([]string, 0, len(dataArr))
	for _, data := range dataArr {
		v := bytes.ToLower(data.Value.([]byte))
		cmd = append(cmd, string(v))
	}

	return commands.Command(cmd).Execute(commands.NewCommandCtx(c.Conn, c.Dict, config.GetServerConfig()))
}
