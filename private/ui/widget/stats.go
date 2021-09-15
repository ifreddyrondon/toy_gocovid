package widget

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/ifreddyrondon/gocovid/private/domain"
)

const layoutUS = "January 2, 2006"

type statsWidget struct {
	*widgets.Paragraph
}

func NewStatsWidget(title string, stats domain.Stats) *statsWidget {
	widget := &statsWidget{
		Paragraph: widgets.NewParagraph(),
	}
	printer := message.NewPrinter(language.English)
	lastUpdated := time.Unix(stats.Updated/1000, 0).Format(layoutUS)
	widget.Title = fmt.Sprintf(" %v - %v", title, lastUpdated)
	widget.Text = printer.Sprintf("Infections: %d (%d today)\n", stats.Cases, stats.TodayCases)
	widget.Text += printer.Sprintf("Deaths %d (%d today)\n", stats.Deaths, stats.TodayDeaths)
	widget.Text += printer.Sprintf("Recoveries: %d (%d remaining)\n", stats.Recovered, stats.Active)
	widget.Text += printer.Sprintf("Critical: %d (%.2f%% of cases)\n", stats.Critical, float64(stats.Critical)/float64(stats.Cases)*100)
	widget.Text += printer.Sprintf("Mortality rate: %.2f%%\n", float64(stats.Deaths)/float64(stats.Cases)*100)
	widget.BorderStyle.Fg = ui.ColorCyan
	return widget
}

func (stats *statsWidget) Draw(buf *ui.Buffer) {
	stats.Paragraph.Draw(buf)
}
