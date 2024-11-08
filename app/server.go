package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/core"
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
		fmt.Println("Failed to bind to port 6379: ", err.Error())
		os.Exit(1)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func() {
			defer c.Close()

			core := core.NewCore(config)
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
