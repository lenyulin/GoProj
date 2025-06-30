package main

import (
	"GoProj/wedy/ioc"
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func main() {
	app := InitWebServer()
	app.cron.Start()
	defer func() {
		//wait for cron job finished
		ctx := app.cron.Stop()
		<-ctx.Done()
	}()
	initPrometheus()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	tpCancel := ioc.InitOTEL()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		tpCancel(ctx)
	}()
	server := app.server
	store := cookie.NewStore([]byte("your-secret-key"))
	server.Use(sessions.Sessions("userId", store))
	server.Run(":8080")
}
func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
