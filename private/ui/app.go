package ui

import (
	"context"
	"fmt"
	"log"
	"net/http"

	ui "github.com/gizak/termui/v3"

	"github.com/ifreddyrondon/gocovid/private/domain"
	"github.com/ifreddyrondon/gocovid/private/fetching"
	"github.com/ifreddyrondon/gocovid/private/ui/country"
	"github.com/ifreddyrondon/gocovid/private/ui/index"
	"github.com/ifreddyrondon/gocovid/private/ui/loading"
)

type (
	Component interface {
		ui.Drawable
		// EventHandler handler
		EventHandler(ui.Event) bool
	}
	Requester interface {
		//Do executes the received *http.Request and returns the generated response, with err in case
		//of transport errors.
		Do(request *http.Request) (response *http.Response, err error)
	}
)

type App struct {
	ctx    context.Context
	cancel func()
	crrCmp Component
	client Requester
}

func NewApp(parentCtx context.Context, client Requester) (*App, context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	return &App{
		cancel: cancel,
		ctx:    parentCtx,
		client: client,
	}, ctx
}

func (a *App) Run() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	setDefaultTermuiColors()
	w, h := ui.TerminalDimensions()
	a.crrCmp = loading.Component(w, h)

	fetcher := fetching.NewIndex(a.client)
	dataChan := fetcher.FetchChan(a.ctx)
	events := ui.PollEvents()
	for {
		select {
		case data := <-dataChan:
			a.crrCmp = index.Component(w, h, data.All, data.Countries, func(data domain.Country) {
				a.crrCmp = country.Component(w, h, data)
			})
		case <-a.ctx.Done():
			if 1 == 1 {

			}
			return
		case e := <-events:
			if a.crrCmp.EventHandler(e) {
				continue
			}
			switch e.ID {
			case "<C-c>":
				fmt.Println(`Type "q" to exit. If you are in filter mode type "Esc" and "q"`)
			case "q":
				a.cancel()
				return
			}
		}
	}
}

func (a *App) Draw() {
	ui.Clear()
	ui.Render(a.crrCmp)
}

func setDefaultTermuiColors() {
	ui.Theme.Default = ui.NewStyle(ui.ColorWhite, ui.ColorClear)
	ui.Theme.Block.Title = ui.NewStyle(ui.ColorWhite, ui.ColorClear)
	ui.Theme.Block.Border = ui.NewStyle(ui.ColorCyan, ui.ColorClear)
}
