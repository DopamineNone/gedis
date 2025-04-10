package database

import (
	"github.com/DopamineNone/gedis/conf"
	"github.com/DopamineNone/gedis/internal/aof"
	"github.com/DopamineNone/gedis/internal/datastruct/dict"
	"github.com/DopamineNone/gedis/internal/resp/conn"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/reply"
	"github.com/google/wire"
	"log"
	"strconv"
	"strings"
)

var (
	ProvideSet = wire.NewSet(New)
)

type DB struct {
	dbSet []*Core
	*aof.AofHandler
}

func New(cfg *conf.Config) handler.Database {
	if cfg.Count == 0 {
		cfg.Count = 16
	}
	set := make([]*Core, cfg.Count)
	db := &DB{dbSet: set}
	if cfg.AppendOnly {
		var err error
		db.AofHandler, err = aof.NewAofHandler(cfg, db)
		if err != nil {
			panic(err)
		}

		for i := range set {
			set[i] = NewCore(i, dict.NewSyncDict(), func(line handler.CmdLine) {
				db.AofHandler.AddAOF(i, line)
			})
		}
		db.LoadAOF()
	} else {
		for i := range set {
			set[i] = NewCore(i, dict.NewSyncDict(), nil)
		}
	}
	return db
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
	return
}

func (db *DB) AfterClientClose(client conn.Conn) {
	//TODO implement me
	return
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
