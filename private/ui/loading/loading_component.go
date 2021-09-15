package loading

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type component struct {
	*ui.Grid
}

func Component(width, height int) *component {
	grid := ui.NewGrid()
	grid.SetRect(0, 0, width, height)
	return &component{grid}
}

func (l *component) Draw(buf *ui.Buffer) {
	widget := widgets.NewParagraph()
	widget.Text = `
                                           d8b      888
                                           Y8P      888
                                                    888
.d88b.   .d88b.   .d8888b .d88b.  888  888 888  .d88888
d88P"88b d88""88b d88P"   d88""88b 888  888 888 d88" 888
888  888 888  888 888     888  888 Y88  88P 888 888  888
Y88b 888 Y88..88P Y88b.   Y88..88P  Y8bd8P  888 Y88b 888
"Y88888  "Y88P"   "Y8888P "Y88P"    Y88P   888  "Y88888
    888
Y8b d88P
"Y88P"


[gocovid](fg:black,bg:yellow) by Freddy Rondon

Worldwide Coronavirus (COVID-19) Statistics for your terminal
	`
	widget.Border = false
	l.Grid.Set(ui.NewRow(1, widget))
	l.Grid.Draw(buf)
}

func (l *component) EventHandler(ui.Event) bool { return false }
