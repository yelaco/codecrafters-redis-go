package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func Psync(ctx CommandCtx, args []string) (resp.RespData, error) {
	go func() {
		ctx.mu.Lock()
		defer ctx.mu.Unlock()
		fileContentHex, err := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(ctx.conn, "$%d\r\n", len(fileContentHex))
		ctx.conn.Write([]byte(fileContentHex))
	}()
	return resp.RespData{
		Value: fmt.Sprintf("FULLRESYNC %v 0", ctx.serverConfig.ReplicationId),
		Type:  resp.SimpleString,
	}, nil
}
