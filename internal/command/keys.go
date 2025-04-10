package command

import (
	"github.com/DopamineNone/gedis/internal/database"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
	"github.com/DopamineNone/gedis/internal/utils"
	"path/filepath"
)

func execDel(db *database.Core, args handler.CmdLine) reply.Reply {
	keys := make([]string, len(args))
	for i, key := range args {
		keys[i] = string(key)
	}
	deleted := db.RemoveMultiKeys(keys...)
	if deleted > 0 {
		db.AddAOF(utils.ToCommandLine("del", args...))
	}
	return reply.NewIntReply(deleted)
}

func execExists(db *database.Core, args handler.CmdLine) reply.Reply {
	result := 0
	for _, key := range args {
		_, exist := db.GetEntity(string(key))
		if exist {
			result++
		}
	}
	return reply.NewIntReply(result)
}

func execKeys(db *database.Core, args handler.CmdLine) reply.Reply {
	pattern := string(args[0])

	result := make([][]byte, 0)
	db.ForEach(func(key string, val any) {
		if match, err := filepath.Match(pattern, key); match && err != nil {
			result = append(result, []byte(key))
		}
	})
	return reply.NewMultiBulkReply(result)
}

func execFlushDB(db *database.Core, args handler.CmdLine) reply.Reply {
	db.Flush()
	db.AddAOF(utils.ToCommandLine("flush", args...))
	return reply.NewOkReply()
}

func execType(db *database.Core, args handler.CmdLine) reply.Reply {
	entity, exist := db.GetEntity(string(args[0]))
	if !exist {
		return reply.NewStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.NewStatusReply("string")
		//TODO:
	}
	return reply.NewUnknownErrReply()
}
func execRename(db *database.Core, args handler.CmdLine) reply.Reply {
	key := string(args[0])
	newKey := string(args[1])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(newKey, entity)
	db.Remove(key)
	db.AddAOF(utils.ToCommandLine("rename", args...))
	return reply.NewOkReply()
}

func execRenameNx(db *database.Core, args handler.CmdLine) reply.Reply {
	key := string(args[0])
	newKey := string(args[1])
	_, exist := db.GetEntity(newKey)
	if exist {
		return reply.NewIntReply(0)
	}
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(newKey, entity)
	db.Remove(key)
	db.AddAOF(utils.ToCommandLine("renamenx", args...))
	return reply.NewIntReply(1)
}

func init() {
	database.RegisterCommand("del", execDel, -2)
	database.RegisterCommand("exists", execExists, -2)
	database.RegisterCommand("keys", execKeys, 2)
	database.RegisterCommand("flushdb", execFlushDB, 1)
	database.RegisterCommand("type", execType, 2)
	database.RegisterCommand("rename", execRename, 3)
	database.RegisterCommand("renamenx", execRenameNx, 3)
}
