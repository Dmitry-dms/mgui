package widgets

import (
	"github.com/Dmitry-dms/mgui/fonts"
	"github.com/Dmitry-dms/mgui/styles"
)

type TextButton struct {
	Id    string
	align TextAlign
	pad   PaddingType
	style *styles.Style
	*Text
	*Button
}

type TextAlign uint32

const (
	Center TextAlign = iota
	Left
	Right
)

// NewTextButton TODO(@Dmitry-dms): Should TextButton exist or just use Text and Button together when needed?
func NewTextButton(id string, x, y, w, h float32,
	text string, tPos []fonts.CombinedCharInfo, align TextAlign, pad PaddingType, style *styles.Style) *TextButton {
	tb := TextButton{
		Id:    id,
		align: align,
		style: style,
		pad:   pad,
	}
	txt := NewText("", text, x, y, w, h, tPos, nil, style, Default)

	btn := NewButton("", x, y, 2*style.Padding+w, 2*style.Padding+h, style.BtnColor)

	tb.Button = btn
	tb.Text = txt

	tb.UpdateTextPos(tb.Button.BoundingBox()[0], tb.Button.BoundingBox()[1])

	return &tb
}
func (tb *TextButton) ToggleUpdate() {
	tb.Text.ToggleUpdate()
	tb.Button.ToggleUpdate()
}
func (tb *TextButton) RenderInfo() ([]float32, []int32, int, int) {
	return tb.Button.RenderInfo()
}

// UpdateTextPos TODO: Need to improve correct calculation of text position
func (tb *TextButton) UpdateTextPos(x, y float32) {
	var xNew, yNew float32
	switch tb.align {
	case Center:
		xNew = ((tb.Width() - tb.Text.Width() - tb.style.Padding*2) / 2) + tb.style.Padding
	case Left:
		xNew = tb.style.Padding
	case Right:
		xNew = tb.Width() - tb.Text.Width() - tb.style.Padding*3
	}
	yNew = tb.style.Padding
	if tb.Width() <= tb.Text.Width()+tb.style.Padding*3 {
		xNew = tb.style.Padding
	}
	tb.Text.UpdatePosition([4]float32{x + xNew, y + yNew, tb.Text.Width(), tb.Text.Height()})
}

func (tb *TextButton) Active() bool {
	return tb.Button.IsActive
}

func (tb *TextButton) Color() [4]float32 {
	return tb.Button.Color()
}
func (tb *TextButton) SetBackgroundColor(clr [4]float32) {
	tb.Button.SetColor(clr)
}

func (tb *TextButton) ChangeActive() {
	tb.Button.ChangeActive()
}

func (tb *TextButton) UpdatePosition(pos [4]float32) {
	tb.Button.UpdatePosition(pos)
	tb.UpdateTextPos(tb.Button.BoundingBox()[0], tb.Button.BoundingBox()[1])
}
func (tb *TextButton) WidgetId() string {
	return tb.Id
}
func (tb *TextButton) Height() float32 {
	return tb.Button.Height()
}

func (tb *TextButton) BoundingBox() [4]float32 {
	return tb.Button.BoundingBox()
}
func (tb *TextButton) Width() float32 {
	return tb.Button.Width()
}
