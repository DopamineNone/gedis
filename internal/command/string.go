package command

import (
	"github.com/DopamineNone/gedis/internal/database"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
)

func execGet(db *database.Core, args handler.CmdLine) reply.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.NewNullBulkReply()
	}
	return reply.NewBulkReply(entity.Data.([]byte))
}

func execSet(db *database.Core, args handler.CmdLine) reply.Reply {
	key := string(args[0])
	value := args[1]
	db.PutEntity(key, &database.DataEntity{Data: value})
	return reply.NewOkReply()
}

func execSetnx(db *database.Core, args handler.CmdLine) reply.Reply {
	key := string(args[0])
	value := args[1]
	return reply.NewIntReply(db.PutIfAbsent(key, &database.DataEntity{Data: value}))
}

func execGetSet(db *database.Core, args handler.CmdLine) reply.Reply {
	key := string(args[0])
	value := string(args[1])
	oldValue, exist := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: value})
	if !exist {
		return reply.NewNullBulkReply()
	}
	return reply.NewBulkReply(oldValue.Data.([]byte))
}

func execStrLen(db *database.Core, args handler.CmdLine) reply.Reply {
	key := string(args[0])
	value, exist := db.GetEntity(key)
	if !exist {
		return reply.NewNullBulkReply()
	}
	return reply.NewIntReply(len(value.Data.([]byte)))
}

func init() {
	database.RegisterCommand("get", execGet, 2)
	database.RegisterCommand("set", execSet, 3)
	database.RegisterCommand("setnx", execSetnx, 3)
	database.RegisterCommand("getset", execGetSet, 3)
	database.RegisterCommand("strlen", execStrLen, 2)
}
