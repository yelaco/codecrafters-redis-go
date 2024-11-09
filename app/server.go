package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/config"
	core "github.com/codecrafters-io/redis-starter-go/internal/core"
	"github.com/codecrafters-io/redis-starter-go/internal/replication"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/v2/parser"
	"github.com/codecrafters-io/redis-starter-go/pkg/util"
)

func main() {
	cfg := config.GetServerConfig()
	addr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to port %s: %s\n", cfg.Port, err.Error())
		os.Exit(1)
	}

	if cfg.Role == config.ROLE_SLAVE {
		err := replication.Handshake(cfg.Port, cfg.MasterHost, cfg.MasterPort)
		if err != nil {
			log.Fatalf("Failed to replicate from %s:%s: %s", cfg.MasterHost, cfg.MasterPort, err.Error())
		}
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
		}

		go func() {
			defer c.Close()

			core := core.NewCore()
			var respParser resp.RespParser

			reader := bufio.NewReader(c)
			for {
				payload := make([]byte, 512)
				_, err = reader.Read(payload)
				if err != nil {
					if err != io.EOF {
						util.DumpLog(fmt.Sprintf("Error reading from connection: %s\n", err.Error()))
					}
					return
				}
				respParser = parser.NewParser(payload)

				data, err := respParser.Parse()
				if err != nil {
					util.DumpLog(fmt.Sprintf("Error parsing payload: %s\n", err.Error()))
					return
				}

				result, err := core.HandleCommand(data)
				if err != nil {
					util.DumpLog(fmt.Sprintf("Error handling command: %s\n", err.Error()))
				}

				_, err = c.Write([]byte(result.String()))
				if err != nil {
					util.DumpLog(fmt.Sprintf("Error writing to connection: %s\n", err.Error()))
					return
				}
			}
		}()
	}
}
