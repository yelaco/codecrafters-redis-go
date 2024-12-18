package commands

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/config"
	"github.com/codecrafters-io/redis-starter-go/internal/replication"
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
	conn         net.Conn
	serverConfig config.ServerConfig
	mu           *sync.Mutex
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

func NewCommandCtx(conn net.Conn, cfg config.ServerConfig) CommandCtx {
	return CommandCtx{
		conn:         conn,
		serverConfig: cfg,
		mu:           &sync.Mutex{},
	}
}

func (cmd Command) IsWriteCmd() bool {
	switch cmd[0] {
	case SET:
		return true
	default:
		return false
	}
}

func (cmd Command) Execute(ctx CommandCtx) (resp.RespData, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
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
	if cmd.IsWriteCmd() {
		replication.QueueCmdIfReplicate(cmd)
	}
	return result, nil
}
