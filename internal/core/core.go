package core

import (
	"bytes"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/config"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type Core struct {
	Dict map[string]string
}

func NewCore() *Core {
	c := &Core{
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

	return commands.Command(cmd).Execute(commands.NewCommandCtx(c.Dict, config.GetServerConfig()))
}
