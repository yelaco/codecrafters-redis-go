package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func Info(ctx CommandCtx, args []string) (resp.RespData, error) {
	switch args[0] {
	case "replication":
		return resp.RespData{
			Value: fmt.Sprintf("# Replication\n%v", replication.GetInfoSection(ctx.serverConfig)),
			Type:  resp.BulkString,
		}, nil
	default:
		return resp.RespData{
			Value: "ERROR: unknown section",
			Type:  resp.SimpleError,
		}, nil
	}
}
