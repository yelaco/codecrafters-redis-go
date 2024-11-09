package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func Ping(ctx CommandCtx, args []string) (resp.RespData, error) {
	return resp.RespData{
		Value: "PONG",
		Type:  resp.SimpleString,
	}, nil
}
