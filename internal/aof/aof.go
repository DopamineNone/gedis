package aof

import (
	"github.com/DopamineNone/gedis/conf"
	"github.com/DopamineNone/gedis/internal/resp/conn"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/resp/parser"
	"github.com/DopamineNone/gedis/internal/resp/reply"
	"github.com/DopamineNone/gedis/internal/utils"
	"io"
	"log"
	"os"
	"strconv"
)

const aofBufferSize = 1 << 16

type AofHandler struct {
	database    handler.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

type payload struct {
	cmdLine handler.CmdLine
	dbIndex int
}

func NewAofHandler(cfg *conf.Config, db handler.Database) (*AofHandler, error) {

	file, err := os.OpenFile(cfg.AOFFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	h := &AofHandler{
		database:    db,
		aofFilename: cfg.AOFFilename,
		aofFile:     file,
		aofChan:     make(chan *payload, aofBufferSize),
		currentDB:   0,
	}
	go h.handleAOF()
	// load aof

	return h, nil
}

func (h *AofHandler) AddAOF(dbIndex int, cmd handler.CmdLine) {
	h.aofChan <- &payload{
		cmdLine: cmd,
		dbIndex: dbIndex,
	}
}
func (h *AofHandler) handleAOF() {
	for p := range h.aofChan {
		if p.dbIndex != h.currentDB {
			data := reply.NewMultiBulkReply(utils.ToCommandLine("select", []byte(strconv.Itoa(p.dbIndex)))).ToBytes()
			_, err := h.aofFile.Write(data)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			h.currentDB = p.dbIndex
		}
		data := reply.NewMultiBulkReply(p.cmdLine).ToBytes()
		_, err := h.aofFile.Write(data)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (h *AofHandler) LoadAOF() {
	file, err := os.Open(h.aofFilename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()
	ch := parser.ParseStream(file)
	mockConn := conn.New(nil)
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			}
			log.Println(err.Error())
			continue
		}
		if p.Data == nil {
			log.Println("empty payload")
			continue
		}
		data, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			log.Println("invalid payload")
			continue
		}
		h.database.Exec(mockConn, data.Args)
	}
}
