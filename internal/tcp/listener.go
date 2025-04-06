package tcp

import (
	"net"

	"github.com/DopamineNone/gedis/conf"
)

func MustListener(cfg *conf.Config) net.Listener {
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		panic(err)
	}
	return listener
}