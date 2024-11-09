package commands

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

type ReplConfigOptions struct {
	listeningPort string
	capabilities  []string
}

func ReplConfig(ctx CommandCtx, args []string) (resp.RespData, error) {
	return resp.RespData{Value: "OK", Type: resp.SimpleString}, nil
}
