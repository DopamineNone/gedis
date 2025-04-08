package database

import (
	"github.com/DopamineNone/gedis/internal/resp/conn"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(New)

type EchoDatabase struct {
}

func New() handler.Database {
	return &EchoDatabase{}
}

func (ed *EchoDatabase) Exec(client conn.Conn, args handler.CmdLine) reply.Reply {
	return reply.NewMultiBulkReply(args)
}

func (ed *EchoDatabase) Close() {
}

func (ed *EchoDatabase) AfterClientClose(client conn.Conn) {
}
