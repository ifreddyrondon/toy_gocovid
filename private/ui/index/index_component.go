package index

import (
	ui "github.com/gizak/termui/v3"

	"github.com/ifreddyrondon/gocovid/private/domain"
	"github.com/ifreddyrondon/gocovid/private/ui/widget"
)

type mode string

const (
	insert mode = "insert"
	view   mode = "view"
)

type state struct {
	crrMode      mode
	countries    []domain.Country
	countriesCpy []domain.Country
	crrTableLine int
}

type w struct {
	all       ui.Drawable
	sort      *SortWidget
	input     *InputWidget
	countries *CountriesWidget
}

type component struct {
	*ui.Grid
	widgets  w
	state    state
	callback func(country domain.Country)
}

func Component(width, height int, all domain.All, countries []domain.Country, cb func(domain.Country)) *component {
	c := &component{
		Grid: ui.NewGrid(),
		state: state{
			crrMode:   view,
			countries: countries,
		},
		callback: cb,
	}
	c.widgets = w{
		all:  widget.NewStatsWidget("Global statistics", all.Stats),
		sort: NewSortWidget(),
		input: NewInputWidget(func(term string) {
			c.filterCountries(domain.FilterByName(term))
		}),
		countries: NewCountriesWidget(func(i int) {
			c.state.crrTableLine = i
		}),
	}
	c.SetRect(0, 0, width, height)
	c.sortCountries(domain.SortByCases)
	c.showSort()
	return c
}

func (c *component) filterCountries(filterFunc domain.FilterFunc) {
	c.state.countriesCpy = filterFunc(c.state.countries)
	c.widgets.countries.Update(c.state.countriesCpy)
}

func (c *component) sortCountries(sort domain.SortFunc) {
	sort(c.state.countries)
	c.state.countriesCpy = c.state.countries
	c.widgets.countries.Update(c.state.countries)
}

func (c *component) Draw(buf *ui.Buffer) {
	c.Grid.Draw(buf)
}

func (c *component) EventHandler(e ui.Event) bool {
	if c.widgets.countries.EventHandler(e) {
		return true
	}
	if c.state.crrMode == insert && c.widgets.input.EventHandler(e) {
		return true
	}
	if c.state.crrMode == view && c.eventHandlerView(e) {
		return true
	}
	if c.eventHandlerInput(e) {
		return true
	}
	switch e.ID {
	case "<Enter>":
		c.callback(c.state.countriesCpy[c.state.crrTableLine])
	default:
		return false
	}
	return true
}

func (c *component) eventHandlerView(e ui.Event) bool {
	switch e.ID {
	case "1":
		c.sortCountries(domain.SortByCases)
	case "2":
		c.sortCountries(domain.SortByCasesToday)
	case "3":
		c.sortCountries(domain.SortByDeaths)
	case "4":
		c.sortCountries(domain.SortByDeathsToday)
	case "5":
		c.sortCountries(domain.SortByRecoveries)
	case "6":
		c.sortCountries(domain.SortByActive)
	case "7":
		c.sortCountries(domain.SortByMortality)
	case "/":
		c.state.crrMode = insert
		c.showInput()
	default:
		return false
	}
	return true
}

func (c *component) eventHandlerInput(e ui.Event) bool {
	switch e.ID {
	case "<Escape>":
		c.state.crrMode = view
		c.sortCountries(domain.SortByCases)
		c.showSort()
	default:
		return false
	}
	return true
}

func (c *component) showSort() {
	c.Set(
		ui.NewRow(0.2, ui.NewCol(1, c.widgets.all)),
		ui.NewRow(0.7, ui.NewCol(1, c.widgets.countries)),
		ui.NewRow(0.1, ui.NewCol(1, c.widgets.sort)),
	)
}

func (c *component) showInput() {
	c.Set(
		ui.NewRow(0.2, ui.NewCol(1, c.widgets.all)),
		ui.NewRow(0.7, ui.NewCol(1, c.widgets.countries)),
		ui.NewRow(0.1, ui.NewCol(1, c.widgets.input)),
	)
}
