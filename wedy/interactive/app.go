package main

import (
	"GoProj/wedy/interactive/events"
	"GoProj/wedy/pkg/grpcx"
)

type App struct {
	server    *grpcx.Server
	consumers []events.Consumer
}
