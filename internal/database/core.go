package database

import (
	"github.com/DopamineNone/gedis/internal/datastruct/dict"
	"github.com/DopamineNone/gedis/internal/resp/conn"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
	"strings"
)

type DataEntity struct {
	Data any
}

type Core struct {
	index  int
	data   dict.Dict
	AddAOF func(line handler.CmdLine)
}

func NewCore(index int, data dict.Dict, addAOF func(line handler.CmdLine)) *Core {
	return &Core{index: index, data: data, AddAOF: addAOF}
}

func (db *Core) Exec(client conn.Conn, args handler.CmdLine) reply.Reply {
	cmdName := strings.ToLower(string(args[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.NewErrReply("ERR unknown command: " + cmdName)
	}
	if !validateArgsLen(cmd.arity, args) {
		return reply.NewArgNumErrReply(cmdName)
	}
	return cmd.executor(db, args[1:])
}

func validateArgsLen(arity int, args handler.CmdLine) bool {
	if arity >= 0 {
		return arity == len(args)
	}
	return -arity <= len(args)
}

func (db *Core) Close() {
	//TODO implement me
	panic("implement me")
}

func (db *Core) AfterClientClose(client conn.Conn) {
	//TODO implement me
	panic("implement me")
}

func (db *Core) GetEntity(key string) (*DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	data, ok := raw.(*DataEntity)
	if !ok {
		return nil, false
	}
	return data, true
}

func (db *Core) PutEntity(key string, entity *DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *Core) PutIfExists(key string, entity *DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *Core) PutIfAbsent(key string, entity *DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *Core) Remove(key string) int {
	return db.data.Remove(key)
}

func (db *Core) RemoveMultiKeys(keys ...string) int {
	deleted := 0
	for _, key := range keys {
		deleted += db.data.Remove(key)
	}
	return deleted
}

func (db *Core) Flush() {
	db.data.Clear()
}

func (db *Core) ForEach(consumer dict.Consumer) {
	db.data.ForEach(consumer)
}
