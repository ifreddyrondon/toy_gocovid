package index

import (
	"fmt"
	"unicode/utf8"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/ifreddyrondon/gocovid/pkg/termui"
	"github.com/ifreddyrondon/gocovid/private/domain"
)

var tableHeader = []string{
	"#",
	"Country",
	"Total Cases",
	"Cases (today)",
	"Total Deaths",
	"Deaths (today)",
	"Recoveries",
	"Active",
	"Critical",
	"Mortality",
}

type CountriesWidget struct {
	*termui.Table
	printer *message.Printer
}

func NewCountriesWidget(scrollCallback func(int)) *CountriesWidget {
	table := termui.NewTable()
	widget := &CountriesWidget{
		Table:   table,
		printer: message.NewPrinter(language.English),
	}
	widget.Title = " Countries "
	widget.ShowCursor = true
	widget.ShowLocation = true
	widget.ColGap = 3
	widget.PadLeft = 2
	widget.ColWidths = []int{3, 20, 12, 12, 12, 12, 12, 12, 12, 12}
	widget.CursorColor = ui.ColorBlue
	widget.ScrollCallback = scrollCallback

	widget.Header = tableHeader
	return widget
}

func (c *CountriesWidget) Update(data []domain.Country) {
	strings := make([][]string, len(data))
	for i, v := range data {
		strings[i] = make([]string, len(tableHeader))
		strings[i][0] = fmt.Sprintf("%d", i+1)
		strings[i][1] = v.Country
		strings[i][2] = c.printer.Sprintf("%d", v.Cases)
		strings[i][3] = c.printer.Sprintf("%d", v.TodayCases)
		strings[i][4] = c.printer.Sprintf("%d", v.Deaths)
		strings[i][5] = c.printer.Sprintf("%d", v.TodayDeaths)
		strings[i][6] = c.printer.Sprintf("%d", v.Recovered)
		strings[i][7] = c.printer.Sprintf("%d", v.Active)
		strings[i][8] = c.printer.Sprintf("%d", v.Critical)
		strings[i][9] = c.printer.Sprintf("%.2f%s", float64(v.Deaths)/float64(v.Cases)*100, "%")
	}
	c.Rows = strings
}

func (c *CountriesWidget) Draw(buf *ui.Buffer) {
	c.Table.Draw(buf)
}

func (c *CountriesWidget) EventHandler(e ui.Event) bool {
	switch e.ID {
	case "k", "<Up>", "<MouseWheelUp>":
		c.ScrollUp()
	case "j", "<Down>", "<MouseWheelDown>":
		c.ScrollDown()
	default:
		return false
	}
	return true
}

type InputWidget struct {
	*termui.Input
	callback func(string)
}

func NewInputWidget(callback func(string)) *InputWidget {
	return &InputWidget{
		Input: &termui.Input{
			Style:         ui.NewStyle(ui.ColorWhite, ui.ColorClear),
			Label:         " Filter: ",
			Value:         "",
			ShowWhenEmpty: true,
		},
		callback: callback,
	}
}

func (i *InputWidget) EventHandler(e ui.Event) bool {
	if utf8.RuneCountInString(e.ID) == 1 {
		i.Value += e.ID
		go i.callback(i.Value)
		return true
	}
	switch e.ID {
	case "<Backspace>":
		if i.Value != "" {
			r := []rune(i.Value)
			i.Value = string(r[:len(r)-1])
			go i.callback(i.Value)
		}
	case "<Space>":
		i.Value += " "
	default:
		return false
	}
	return true
}

func (i *InputWidget) Draw(buf *ui.Buffer) {
	i.Input.Draw(buf)
}

type SortWidget struct {
	*widgets.Paragraph
}

func NewSortWidget() *SortWidget {
	widget := &SortWidget{
		Paragraph: widgets.NewParagraph(),
	}
	widget.Title = " Sorting options "
	widget.Text = fmt.Sprintf("1: Total Cases | 2: Cases Today | 3: Total Deaths | 4: Deaths Today | 5: Recoveries | 6: Active | 7: Mortality")
	return widget
}

func (w *SortWidget) Draw(buf *ui.Buffer) {
	w.Paragraph.Draw(buf)
}
