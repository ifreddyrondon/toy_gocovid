package country

import (
	"math"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/ifreddyrondon/gocovid/private/domain"
)

type plotWidget struct {
	*widgets.Plot
}

func NewPlotWidget(data domain.HistoricalCountry) *plotWidget {
	plotData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = 1 + math.Sin(float64(i)/5)
			data[1][i] = 1 + math.Cos(float64(i)/5)
		}
		return data
	}()

	w := plotWidget{Plot: widgets.NewPlot()}
	w.Title = "braille-mode Line Chart"
	w.Data = plotData
	w.AxesColor = ui.ColorWhite
	w.LineColors[0] = ui.ColorGreen
	return &w
}

func (w *plotWidget) Draw(buf *ui.Buffer) {
	w.Plot.Draw(buf)
}
