//go:build wireinject
// +build wireinject

package main

import (
	"github.com/DopamineNone/gedis/app"
	"github.com/DopamineNone/gedis/conf"
	"github.com/DopamineNone/gedis/tcp"
	"github.com/google/wire"
)

func wireApp() *app.App {
	panic(wire.Build(tcp.ProvideSet, app.ProvideSet, conf.ProvideSet))
}
