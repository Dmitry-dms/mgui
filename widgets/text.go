package widgets

import (
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

	Editor *Editor
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

func NewText(id, text string, x, y float32, font *fonts.Font, style *styles.Style, flag TextFlag) *Text {
	t := Text{
		Message: text,
		baseWidget: baseWidget{
			id:              id,
			BackgroundColor: style.TransparentColor,
		},
		CurrentColor: style.TextColor,
		Size:         style.TextSize,
		Padding:      style.TextPadding * int(style.FontScale),
		Scale:        style.FontScale,
		Flag:         flag,
		Editor:       NewEditor(font, style.FontScale),
	}
	t.Editor.InsertText(text)
	t.SetWH(t.Editor.TextWidth, t.Editor.TextHeight)
	t.boundingBox = [4]float32{x, y, t.Editor.TextWidth, t.Editor.TextHeight + float32(style.TextPadding)}
	return &t
}

func (t *Text) UpdatePosition(pos [4]float32) {
	t.updatePosition(pos)
}

func (t *Text) SetWH(width, height float32) {
	t.boundingBox[2] = width
	t.boundingBox[3] = height //+ float32(t.AllPadding)
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
func (t *Text) Id() string {
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
