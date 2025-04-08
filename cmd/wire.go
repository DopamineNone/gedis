//go:build wireinject
// +build wireinject

package main

import (
	"github.com/DopamineNone/gedis/conf"
	"github.com/DopamineNone/gedis/internal/app"
	"github.com/DopamineNone/gedis/internal/database"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/tcp"
	"github.com/google/wire"
)

func wireApp() *app.App {
	panic(wire.Build(handler.ProvideSet, app.ProvideSet, conf.ProvideSet, database.ProvideSet, tcp.ProvideSet))
}
