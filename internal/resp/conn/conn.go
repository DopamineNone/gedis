package conn

import (
	"net"
	"sync"
	"time"
)

type Conn interface {
	GetRemoteAdd() net.Addr
	Close() error
	Write(input []byte) error
	GetDBIndex() int
	SelectDB(n int)
}

type RespConn struct {
	net.Conn
	mu         sync.Mutex
	selectedDB int
	wait       time.Duration
}

func (rc *RespConn) GetRemoteAdd() net.Addr {
	return rc.Conn.RemoteAddr()
}

func (rc *RespConn) Close() error {
	if !rc.mu.TryLock() {
		time.Sleep(rc.wait)
	}
	return rc.Conn.Close()
}

func (rc *RespConn) Write(input []byte) error {
	if len(input) == 0 {
		return nil
	}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	_, err := rc.Conn.Write(input)
	return err
}

func (rc *RespConn) GetDBIndex() int {
	return rc.selectedDB
}

func (rc *RespConn) SelectDB(n int) {
	rc.selectedDB = n
}

func New(c net.Conn) Conn {
	return &RespConn{
		Conn: c,
	}
}
