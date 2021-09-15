package termui

import (
	"image"
	"strings"

	. "github.com/gizak/termui/v3"
	rw "github.com/mattn/go-runewidth"

	"github.com/ifreddyrondon/gocovid/pkg/runes"
)

const (
	ELLIPSIS = "…"
	CURSOR   = " "
)

type Input struct {
	Block

	Style Style

	Label         string
	Value         string
	ShowWhenEmpty bool
}

func (e *Input) Draw(buf *Buffer) {
	if e.Value == "" && !e.ShowWhenEmpty {
		return
	}

	style := e.Style
	label := e.Label
	label += "["
	style = NewStyle(style.Fg, style.Bg, ModifierBold)
	cursorStyle := NewStyle(style.Bg, style.Fg, ModifierClear)

	p := image.Pt(e.Min.X, e.Min.Y)
	buf.SetString(label, style, p)
	p.X += rw.StringWidth(label)

	tail := "] "
	maxLen := e.Max.X - p.X - rw.StringWidth(tail) - 1
	value := runes.TruncateFront(e.Value, maxLen, ELLIPSIS)
	buf.SetString(value, e.Style, p)
	p.X += rw.StringWidth(value)

	buf.SetString(CURSOR, cursorStyle, p)
	p.X += rw.StringWidth(CURSOR)
	if remaining := maxLen - rw.StringWidth(value); remaining > 0 {
		buf.SetString(strings.Repeat(" ", remaining), e.TitleStyle, p)
		p.X += remaining
	}
	buf.SetString(tail, style, p)
}
