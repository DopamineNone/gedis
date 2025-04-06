package tcp

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/DopamineNone/gedis/internal/app"
)

type EchoHandler struct {
	connSet   sync.Map
	isClosing atomic.Bool
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.isClosing.Load() {
		_ = conn.Close()
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Println(err.Error())
				return
			}
			// return content
			conn.Write([]byte(msg))
		}
	}
}

func (h *EchoHandler) Close() error {
	h.isClosing.Store(true)
	return nil
}

func NewHandler() app.Handler {
	return &EchoHandler{
		connSet: sync.Map{},
	}
}
