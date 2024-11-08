package replication

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/resp/v2/parser"
)

func Handshake(host, port string) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	stages := [][]string{
		{"PING"},
		// {"REPLCONF"},
		// {"PSYNC"},
	}
	reader := bufio.NewReader(c)
	var respParser resp.RespParser

	for _, stage := range stages {
		_, err = c.Write((resp.FormatCommand(stage)))
		if err != nil {
			return err
		}

		payload := make([]byte, 512)
		_, err = reader.Read(payload)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("Error reading from connection: %w", err)
			}
		}
		respParser = parser.NewParser(payload)

		_, err := respParser.Parse()
		if err != nil {
			return fmt.Errorf("Error parsing payload: %w", err)
		}
	}

	return nil
}
