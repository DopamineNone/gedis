package database

import (
	"github.com/DopamineNone/gedis/conf"
	"github.com/DopamineNone/gedis/internal/datastruct/dict"
	"github.com/DopamineNone/gedis/internal/resp/conn"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
	"github.com/google/wire"
	"log"
	"strconv"
	"strings"
)

var ProvideSet = wire.NewSet(New)

type DB struct {
	dbSet []*Core
}

func New(cfg *conf.Config) handler.Database {
	if cfg.Count == 0 {
		cfg.Count = 16
	}
	set := make([]*Core, cfg.Count)
	for i := range set {
		set[i] = NewCore(i, dict.NewSyncDict())
	}
	return &DB{dbSet: set}
}

func (db *DB) Exec(client conn.Conn, args handler.CmdLine) reply.Reply {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.NewArgNumErrReply(cmdName)
		}
		return execSelect(client, db, args)
	}
	dbIndex := client.GetDBIndex()
	dbCore := db.dbSet[dbIndex]
	return dbCore.Exec(client, args)
}

func (db *DB) Close() {
	//TODO implement me
	panic("implement me")
}

func (db *DB) AfterClientClose(client conn.Conn) {
	//TODO implement me
	panic("implement me")
}

func execSelect(c conn.Conn, database *DB, args [][]byte) reply.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewErrReply("ERR invalid db index")
	}
	if dbIndex >= len(database.dbSet) {
		return reply.NewErrReply("ERR index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.NewOkReply()
}
