package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/core"
	"github.com/codecrafters-io/redis-starter-go/core/replication"
	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/resp/v2/parser"
	"github.com/codecrafters-io/redis-starter-go/util"
)

func main() {
	config, err := util.NewConfig()
	if err != nil {
		panic(err)
	}
	addr := fmt.Sprintf("%v:%v", config.Host, config.Port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind to port %s: %s\n", config.Port, err.Error())
		os.Exit(1)
	}

	serverInfo := core.NewServerInfo(config)
	if serverInfo.Replication.IsReplica() {
		masterHost := *serverInfo.Replication.MasterHost
		masterPort := *serverInfo.Replication.MasterPort
		err := replication.Handshake(config.Port, masterHost, masterPort)
		if err != nil {
			log.Fatalf("Failed to replicate from %s:%s: %s", masterHost, masterPort, err.Error())
		}
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
		}

		go func() {
			defer c.Close()

			core := core.NewCore(config, serverInfo)
			var respParser resp.RespParser

			reader := bufio.NewReader(c)
			for {
				payload := make([]byte, 512)
				_, err = reader.Read(payload)
				if err != nil {
					if err != io.EOF {
						fmt.Println("Error reading from connection: ", err.Error())
					}
					return
				}
				respParser = parser.NewParser(payload)

				data, err := respParser.Parse()
				if err != nil {
					fmt.Println("Error parsing payload: ", err.Error())
					return
				}

				result, err := core.HandleCommand(data)
				if err != nil {
					fmt.Println("Error handling command: ", err.Error())
				}

				_, err = c.Write([]byte(result.String()))
				if err != nil {
					fmt.Println("Error writing to connection: ", err.Error())
					return
				}
			}
		}()
	}
}
