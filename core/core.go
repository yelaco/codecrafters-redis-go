package core

import (
	"bytes"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/util"
)

type Core struct {
	db   map[string]string
	info ServerInfo
}

func NewCore(config util.Config) *Core {
	return &Core{
		db:   make(map[string]string),
		info: NewServerInfo(config),
	}
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

	switch cmd[0] {
	case "ping":
		return c.Ping()
	case "echo":
		if len(cmd) >= 2 {
			return c.Echo(cmd[1])
		}
	case "set":
		// TODO: parse the command
		if len(cmd) >= 3 {
			options := SetOptions{}
			if len(cmd) >= 5 {
				t, _ := strconv.Atoi(cmd[4])
				if cmd[3] == "px" {
					options.expireType = PX
					options.expireTime = int64(t)
				}
			}
			return c.Set(cmd[1], cmd[2], options)
		}
	case "get":
		if len(cmd) >= 2 {
			return c.Get(cmd[1])
		}
	case "info":
		if len(cmd) >= 2 {
			return c.Info(cmd[1])
		}
	default:
		return resp.RespData{
			Value: "ERROR: command not found",
			Type:  resp.SimpleError,
		}, nil
	}

	return resp.RespData{
		Value: "ERROR: invalid command",
		Type:  resp.SimpleError,
	}, nil
}
