// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/DopamineNone/gedis/conf"
	"github.com/DopamineNone/gedis/internal/app"
	"github.com/DopamineNone/gedis/internal/database"
	"github.com/DopamineNone/gedis/internal/resp/handler"
	"github.com/DopamineNone/gedis/internal/tcp"
)

// Injectors from wire.go:

func wireApp() *app.App {
	config := conf.New()
	handlerDatabase := database.New(config)
	appHandler := handler.NewHandler(handlerDatabase)
	listener := tcp.MustListener(config)
	appApp := app.New(appHandler, listener)
	return appApp
}
