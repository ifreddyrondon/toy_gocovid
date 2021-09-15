package termui

import (
	"fmt"
	"image"
	"log"
	"strings"

	. "github.com/gizak/termui/v3"
)

type Table struct {
	*Block

	Header []string
	Rows   [][]string

	ColWidths []int
	ColGap    int
	PadLeft   int

	ShowCursor  bool
	CursorColor Color

	ShowLocation bool

	UniqueCol    int    // the column used to uniquely identify each table row
	SelectedItem string // used to keep the cursor on the correct item if the data changes
	SelectedRow  int
	TopRow       int // used to indicate where in the table we are scrolled at

	ColResizer     func()
	ScrollCallback func(int)
}

// NewTable returns a new Table3 instance
func NewTable() *Table {
	return &Table{
		Block:       NewBlock(),
		SelectedRow: 0,
		TopRow:      0,
		UniqueCol:   0,
		ColResizer:  func() {},
	}
}

func (t *Table) Draw(buf *Buffer) {
	t.Block.Draw(buf)

	if t.ShowLocation {
		t.drawPagination(buf)
	}

	t.ColResizer()

	// finds exact column starting position
	var colXPos []int
	cur := 1 + t.PadLeft
	for _, w := range t.ColWidths {
		colXPos = append(colXPos, cur)
		cur += w
		cur += t.ColGap
	}

	// prints header
	for i, h := range t.Header {
		width := t.ColWidths[i]
		if width == 0 {
			continue
		}
		// don't render column if it doesn't fit in widget
		if width > (t.Inner.Dx()-colXPos[i])+1 {
			continue
		}
		buf.SetString(
			h,
			NewStyle(Theme.Default.Fg, ColorClear, ModifierBold),
			image.Pt(t.Inner.Min.X+colXPos[i]-1, t.Inner.Min.Y),
		)
	}

	if t.TopRow < 0 {
		log.Printf("table widget TopRow value less than 0. TopRow: %v", t.TopRow)
		return
	}

	// prints each row
	for rowNum := t.TopRow; rowNum < t.TopRow+t.Inner.Dy()-1 && rowNum < len(t.Rows); rowNum++ {
		row := t.Rows[rowNum]
		y := (rowNum + 2) - t.TopRow

		// prints cursor
		style := NewStyle(Theme.Default.Fg)
		if t.ShowCursor {
			if (t.SelectedItem == "" && rowNum == t.SelectedRow) || (t.SelectedItem != "" && t.SelectedItem == row[t.UniqueCol]) {
				style.Fg = t.CursorColor
				style.Modifier = ModifierReverse
				for _, width := range t.ColWidths {
					if width == 0 {
						continue
					}
					buf.SetString(
						strings.Repeat(" ", t.Inner.Dx()),
						style,
						image.Pt(t.Inner.Min.X, t.Inner.Min.Y+y-1),
					)
				}
				t.SelectedItem = row[t.UniqueCol]
				t.SelectedRow = rowNum
			}
		}

		// prints each col of the row
		for i, width := range t.ColWidths {
			if width == 0 {
				continue
			}
			// don't render column if width is greater than distance to end of widget
			if width > (t.Inner.Dx()-colXPos[i])+1 {
				continue
			}
			r := TrimString(row[i], width)
			buf.SetString(
				r,
				style,
				image.Pt(t.Inner.Min.X+colXPos[i]-1, t.Inner.Min.Y+y-1),
			)
		}
	}
}

func (t *Table) drawPagination(buf *Buffer) {
	total := len(t.Rows)
	topRow := t.TopRow + 1
	if topRow > total {
		topRow = total
	}
	bottomRow := t.TopRow + t.Inner.Dy() - 1
	if bottomRow > total {
		bottomRow = total
	}

	loc := fmt.Sprintf(" %d - %d of %d ", topRow, bottomRow, total)

	width := len(loc)
	buf.SetString(loc, t.TitleStyle, image.Pt(t.Max.X-width-2, t.Min.Y))
}

func (t *Table) calcPos() {
	t.SelectedItem = ""

	if t.SelectedRow < 0 {
		t.SelectedRow = 0
	}
	if t.SelectedRow < t.TopRow {
		t.TopRow = t.SelectedRow
	}

	if t.SelectedRow > len(t.Rows)-1 {
		t.SelectedRow = len(t.Rows) - 1
	}
	if t.SelectedRow > t.TopRow+(t.Inner.Dy()-2) {
		t.TopRow = t.SelectedRow - (t.Inner.Dy() - 2)
	}
}

func (t *Table) ScrollUp() {
	t.SelectedRow--
	t.calcPos()
	t.ScrollCallback(t.SelectedRow)
}

func (t *Table) ScrollDown() {
	t.SelectedRow++
	t.calcPos()
	t.ScrollCallback(t.SelectedRow)
}

func (t *Table) ScrollTop() {
	t.SelectedRow = 0
	t.calcPos()
	t.ScrollCallback(t.SelectedRow)
}

func (t *Table) ScrollBottom() {
	t.SelectedRow = len(t.Rows) - 1
	t.calcPos()
	t.ScrollCallback(t.SelectedRow)
}
