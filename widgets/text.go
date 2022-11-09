package widgets

import (
	"fmt"
	"github.com/Dmitry-dms/mgui/fonts"
	"github.com/Dmitry-dms/mgui/styles"
)

type Text struct {
	baseWidget
	Message      string
	CurrentColor [4]float32
	Chars        []fonts.CombinedCharInfo // Slice of all chars. It calculates only when Message changes. Is used in draw package.
	Flag         TextFlag
	Lines        []fonts.TextLine // Lines of text, each Line contains its own Chars. Used when user selects texts.

	Size    int
	Padding int
	Scale   float32

	StartInd, StartLine int
	EndInd, EndLine     int

	//MultiLineTextInput
	LastWidth             float32
	CursorInd, CursorLine int
}

type TextFlag uint

const (
	Selectable TextFlag = 1 << iota
	SplitWords
	SplitChars
	FitContent
	Editable
	Default

	// MultiLineTextInput flags
	MultiLine
)

func NewText(id, text string, x, y, w, h float32, chars []fonts.CombinedCharInfo, l []fonts.TextLine, style *styles.Style, flag TextFlag) *Text {
	t := Text{
		Message: text,
		Chars:   chars,
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h + float32(style.TextPadding)},
			BackgroundColor: style.TransparentColor,
		},
		CurrentColor: style.TextColor,
		Size:         style.TextSize,
		Padding:      style.TextPadding * int(style.FontScale),
		Scale:        style.FontScale,
		Flag:         flag,
		Lines:        l,
	}
	return &t
}

func (t *Text) CursorHelper(dx int) {
	if dx > 0 {
		if t.CursorInd+dx <= len(t.Lines[t.CursorLine].Text) {
			t.CursorInd += dx
		} else {
			if t.CursorLine < len(t.Lines)-1 {
				t.CursorLine++
				t.CursorInd = 0
			}
		}
	} else {
		if t.CursorLine != 0 {
			if t.CursorInd == 0 {
				t.CursorInd = len(t.Lines[t.CursorLine-1].Text)
				t.CursorLine--
			} else {
				t.CursorInd--
			}
		} else {
			if t.CursorInd != 0 {
				t.CursorInd--
			}
		}
	}
	fmt.Println(t.CursorLine, t.CursorInd)
}

func (t *Text) CalculateCursorPos() (x, y, w, h float32) {
	line := t.Lines[t.CursorLine]
	x = line.StartX
	y = line.StartY
	if t.CursorInd >= len(line.Text) {
		x += line.Width
	} else {
		char := line.Text[t.CursorInd]
		x += char.Pos.X - float32(char.Char.LeftBearing)
	}

	h = line.Height
	w = 5
	return
}

func (t *Text) UpdatePosition(pos [4]float32) {
	t.updatePosition(pos)
	//t.Base.boundingBox = pos
}

func (t *Text) SetWH(width, height float32) {
	t.boundingBox[2] = width
	t.boundingBox[3] = height //+ float32(t.Padding)
}
func (t *Text) SetTextColor(clr [4]float32) {
	if clr == t.CurrentColor {
		return
	}
	t.Updated = true
	t.CurrentColor = clr
}
func (t *Text) SetBackGroundColor(clr [4]float32) {
	if clr == t.BackgroundColor {
		return
	}
	t.BackgroundColor = clr
}

func (t *Text) BoundingBox() [4]float32 {
	return t.boundingBox
}
func (t *Text) GetBackgroundColor() [4]float32 {
	return t.BackgroundColor
}
func (t *Text) Color() [4]float32 {
	return t.CurrentColor
}
func (t *Text) WidgetId() string {
	return t.id
}

func (t *Text) Height() float32 {
	return t.height()
}
func (t *Text) Visible() bool {
	return true
}
func (t *Text) Width() float32 {
	return t.width()
}
