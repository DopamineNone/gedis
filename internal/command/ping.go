package command

import (
	"github.com/DopamineNone/gedis/internal/database"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
)

func Ping(db *database.Core, args handler.CmdLine) reply.Reply {
	return reply.NewPongReply()
}

func init() {
	database.RegisterCommand("ping", Ping, 1)
}
