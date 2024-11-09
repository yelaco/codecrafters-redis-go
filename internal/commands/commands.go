package commands

import (
	"context"
	"errors"

	"github.com/codecrafters-io/redis-starter-go/internal/config"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

var (
	ErrUnknownCommand = errors.New("unknown command")
	ErrInvalidCommand = errors.New("invalid command")
)

const (
	PING     = "ping"
	ECHO     = "echo"
	SET      = "set"
	GET      = "get"
	INFO     = "info"
	REPLCONF = "replconf"
	PSYNC    = "psync"
)

type Command []string

type CommandCtx struct {
	context.Context
	dict         map[string]string
	serverConfig config.ServerConfig
}

var cmds = map[string]func(CommandCtx, []string) (resp.RespData, error){
	PING:     Ping,
	ECHO:     Echo,
	SET:      Set,
	GET:      Get,
	INFO:     Info,
	REPLCONF: ReplConfig,
	PSYNC:    Psync,
}

func NewCommandCtx(dict map[string]string, cfg config.ServerConfig) CommandCtx {
	return CommandCtx{
		dict:         dict,
		serverConfig: cfg,
	}
}

func (cmd Command) Execute(ctx CommandCtx) (resp.RespData, error) {
	f, ok := cmds[cmd[0]]
	if !ok {
		return resp.RespData{
			Value: "ERROR: Unknown command",
			Type:  resp.SimpleError,
		}, ErrUnknownCommand
	}
	result, err := f(ctx, cmd[1:])
	if err != nil {
		return resp.RespData{
			Value: "ERROR: Invalid command",
			Type:  resp.SimpleError,
		}, ErrInvalidCommand
	}
	return result, nil
}
