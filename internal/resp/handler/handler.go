package handler

import (
	"context"
	"github.com/google/wire"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/DopamineNone/gedis/internal/app"
	"github.com/DopamineNone/gedis/internal/resp/conn"
	"github.com/DopamineNone/gedis/internal/resp/parser"
	"github.com/DopamineNone/gedis/internal/resp/reply"
)

var ProvideSet = wire.NewSet(NewHandler)

type CmdLine = [][]byte

type Database interface {
	Exec(client conn.Conn, args CmdLine) reply.Reply
	Close()
	AfterClientClose(client conn.Conn)
}

type RespHandler struct {
	connSet   sync.Map
	isClosing atomic.Bool
	db        Database
}

func (rh *RespHandler) closeClient(client conn.Conn) {
	_ = client.Close()
	rh.db.AfterClientClose(client)
	rh.connSet.Delete(client)
}

func (rh *RespHandler) Handle(ctx context.Context, c net.Conn) {
	if rh.isClosing.Load() {
		_ = c.Close()
	}

	client := conn.New(c)
	rh.connSet.Store(client, struct{}{})

	for payload := range parser.ParseStream(c) {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// close
				rh.closeClient(client)
				log.Println("connection closed: " + client.GetRemoteAdd().String())
				return
			}
			result := reply.NewErrReply(payload.Err.Error())
			err := client.Write(result.ToBytes())
			if err != nil {
				rh.closeClient(client)
				return
			}
			continue
		}
		if payload.Data != nil {
			data, ok := payload.Data.(*reply.MultiBulkReply)
			if !ok {
				//log.Println("required multi bulk reply")
				client.Write(payload.Data.ToBytes())
				continue
			}
			result := rh.db.Exec(client, data.Args)
			if result != nil {
				client.Write(result.ToBytes())
			} else {
				client.Write(reply.NewUnknownErrReply().ToBytes())
			}
		}
	}
}

func (rh *RespHandler) Close() error {
	log.Println("handler shutting down")
	rh.isClosing.Store(true)
	rh.connSet.Range(func(client, _ any) bool {
		client.(conn.Conn).Close()
		return true
	})
	rh.db.Close()
	return nil
}

func NewHandler(db Database) app.Handler {
	return &RespHandler{
		connSet: sync.Map{},
		db:      db,
	}
}
