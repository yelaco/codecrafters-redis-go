package replication

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

const (
	REPL_ID_LEN = 40
	ROLE_MASTER = "master"
	ROLE_SLAVE  = "slave"
)

type Replicator struct {
	replCh  chan []string
	propChs []chan []string
}

func newReplicator() *Replicator {
	repl := &Replicator{
		replCh:  make(chan []string),
		propChs: []chan []string{},
	}
	go func() {
		for {
			cmd := <-repl.replCh
			for _, propCh := range repl.propChs {
				propCh <- cmd
			}
		}
	}()
	return repl
}

var replicator *Replicator

// Master node start the replication process to the slave node
func StartReplication(c net.Conn) {
	if replicator == nil {
		replicator = newReplicator()
	}
	fileContentHex, err := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(c, "$%d\r\n", len(fileContentHex))
	if err != nil {
		return
	}
	_, err = c.Write([]byte(fileContentHex))
	if err != nil {
		return
	}

	propCh := make(chan []string)
	replicator.propChs = append(replicator.propChs, propCh)
	go replicator.PropagateCmd(c, propCh)
}

func (r *Replicator) PropagateCmd(c net.Conn, propCh chan []string) {
	for {
		cmd := <-propCh
		_, err := c.Write([]byte(resp.FormatCommand(cmd)))
		if err != nil {
			return
		}
	}
}

func QueueCmdIfReplicate(cmd []string) {
	if replicator != nil {
		replicator.replCh <- cmd
	}
}
