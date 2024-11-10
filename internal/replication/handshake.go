package replication

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/v2/parser"
)

func Handshake(port, masterHost, masterPort string) error {
	addr := fmt.Sprintf("%s:%s", masterHost, masterPort)
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	stages := [][]string{
		{"PING"},
		{"REPLCONF", "listening-port", port},
		{"REPLCONF", "capa", "eof", "capa", "psync2"},
		{"PSYNC", "?", "-1"},
	}
	reader := bufio.NewReader(c)
	var respParser resp.RespParser

	for _, stage := range stages {
		_, err = c.Write((resp.FormatCommand(stage)))
		if err != nil {
			return fmt.Errorf("Handshake: Error writing connection: %w", err)
		}

		payload := make([]byte, 512)
		_, err = reader.Read(payload)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("Handshake: Error reading from connection: %w", err)
			}
		}
		respParser = parser.NewParser(payload)

		_, err := respParser.Parse()
		if err != nil {
			return fmt.Errorf("Error parsing payload: %w", err)
		}
	}

	rdbFile := make([]byte, 512)
	_, err = reader.Read(rdbFile)
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("Handshake: Error reading from connection: %w", err)
		}
	}

	return nil
}
