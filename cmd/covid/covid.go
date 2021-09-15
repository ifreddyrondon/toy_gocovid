package main

import (
	"context"
	"net/http"
	"time"

	"github.com/ifreddyrondon/gocovid/private/ui"
)

const defaultFPS = 100

func main() {
	var c http.Client
	app, ctx := ui.NewApp(context.Background(), &c)
	go app.Run()
	ticks := time.Tick(time.Second / defaultFPS)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			app.Draw()
		}
	}
}
