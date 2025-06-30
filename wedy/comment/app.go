package comment

import "GoProj/wedy/pkg/grpcx"

type App struct {
	server    *grpcx.Server
	consumers []events.Consumer
}
