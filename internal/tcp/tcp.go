package tcp

import "github.com/google/wire"

var ProvideSet = wire.NewSet(MustListener)
