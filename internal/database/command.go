package database

import (
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
	"strings"
)

var cmdTable = make(map[string]*command)

type ExecFunc func(db *Core, args handler.CmdLine) reply.Reply

type command struct {
	executor ExecFunc
	arity    int
}

func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}
