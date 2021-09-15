package country

import (
	"net/http"

	ui "github.com/gizak/termui/v3"

	"github.com/ifreddyrondon/gocovid/private/domain"
	"github.com/ifreddyrondon/gocovid/private/fetching"
	"github.com/ifreddyrondon/gocovid/private/ui/widget"
)

type w struct {
	stats ui.Drawable
	plot  ui.Drawable
}

type state struct {
}

type component struct {
	*ui.Grid
	widgets w
	state   state
}

func Component(width, height int, data domain.Country) *component {
	c := &component{
		Grid: ui.NewGrid(),
		state: state{
		},
	}
	c.widgets = w{
		stats: widget.NewStatsWidget(data.Country, data.Stats),
	}
	fetcher := fetching.Historical{Doer: &http.Client{}}
	dataChan := fetcher.FetchCountryChan(data.CountryInfo.Iso2)
	select {
	case data := <-dataChan:
		c.widgets.plot = NewPlotWidget(data.HistoricalCountry)
	}
	c.SetRect(0, 0, width, height)
	c.Set(
		ui.NewRow(0.2, ui.NewCol(1, c.widgets.stats)),
		ui.NewRow(0.8, ui.NewCol(1, c.widgets.plot)),
	)
	return c
}

func (c *component) Draw(buf *ui.Buffer) {
	c.Grid.Draw(buf)
}

func (c *component) EventHandler(ui.Event) bool { return false }
