package replication

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/resp/v2/parser"
)

// Slave node handshakes with the master node
func Handshake(port, masterHost, masterPort string) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%s", masterHost, masterPort)
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	stages := [][]string{
		{"PING"},
		{"REPLCONF", "listening-port", port},
		{"REPLCONF", "capa", "psync2"},
		{"PSYNC", "?", "-1"},
	}
	reader := bufio.NewReader(c)
	var respParser resp.RespParser

	// PING stage
	_, err = c.Write((resp.FormatCommand(stages[0])))
	if err != nil {
		return nil, fmt.Errorf("error writing connection: %w", err)
	}
	payload := make([]byte, 512)
	_, err = reader.Read(payload)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error reading from connection: %w", err)
		}
	}
	respParser = parser.NewParser(payload)
	_, _ = respParser.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %w", err)
	}

	// REPLCONF stage
	_, err = c.Write((resp.FormatCommand(stages[1])))
	if err != nil {
		return nil, fmt.Errorf("error writing connection: %w", err)
	}
	payload = make([]byte, 512)
	_, err = reader.Read(payload)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error reading from connection: %w", err)
		}
	}
	respParser = parser.NewParser(payload)
	_, _ = respParser.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %w", err)
	}
	_, err = c.Write((resp.FormatCommand(stages[2])))
	if err != nil {
		return nil, fmt.Errorf("error writing connection: %w", err)
	}
	payload = make([]byte, 512)
	_, err = reader.Read(payload)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error reading from connection: %w", err)
		}
	}
	respParser = parser.NewParser(payload)
	_, _ = respParser.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %w", err)
	}

	// PSYNC stage
	_, err = c.Write((resp.FormatCommand(stages[3])))
	if err != nil {
		return nil, fmt.Errorf("error writing connection: %w", err)
	}

	payload = make([]byte, 4096)
	_, err = reader.Read(payload)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error reading from connection: %w", err)
		}
	}
	respParser = parser.NewParser(payload)

	_, _ = respParser.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %w", err)
	}

	rdbFile := make([]byte, 512)
	_, err = reader.Read(rdbFile)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error reading from connection: %w", err)
		}
	}
	fmt.Println(string(rdbFile))
	return c, nil
}

func ReplicateFrom(c net.Conn, handler func(resp.RespData) (resp.RespData, error)) {
	var respParser resp.RespParser
	reader := bufio.NewReader(c)
	for {
		payload := make([]byte, 512)
		_, err := reader.Read(payload)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from connection: %s\n", err.Error())
			}
			// return
		}
		if payload[0] == '\n' {
			continue
		} else if payload[0] == 0 {
			return
		}
		respParser = parser.NewParser(payload)

		fmt.Println(string(payload))

		data, err := respParser.Parse()
		if err != nil {
			fmt.Printf("Error parsing payload: %s\n", err.Error())
			// return
		}

		fmt.Println(data)

		_, err = handler(data)
		if err != nil {
			fmt.Printf("Error handling command: %s\n", err.Error())
		}
	}
}
