package widgets

import (
	"github.com/Dmitry-dms/mgui/styles"
)

type HybridLayout struct {
	X, Y, W, H            float32
	CursorX, CursorY      float32
	InitY                 float32
	LastWidth, LastHeight float32

	RequireColumn            bool
	CurrentColH, CurrentColW float32
	// CurrentCol     *Column

	// ColumnW, ColumnH float32
	Padding     int
	ItemSpacing int
	style       *styles.Style
	Id          string

	Align RowAlign
}

type Column struct {
	Id   string
	W, H float32
}

type RowAlign uint

const (
	VerticalAlign RowAlign = 1 << iota

	NoAlign
)

func (r *HybridLayout) UpdateColWidth(w float32) {
	if r.CurrentColW < w {
		r.CurrentColW = w
	}
}
func (r *HybridLayout) AddColHeight(h float32) {
	r.CurrentColH += h
}

func NewHLayout(id string, x, y float32, a RowAlign, style *styles.Style) *HybridLayout {
	r := HybridLayout{
		X:     x,
		Y:     y,
		style: style,
		Id:    id,
		Align: a,
	}
	return &r
}

func (r *HybridLayout) Height() float32 {
	return r.LastHeight
}
func (r *HybridLayout) Width() float32 {
	return r.LastWidth
}
func (r *HybridLayout) WidgetId() string {
	return r.Id
}

func (r *HybridLayout) UpdateHeight(h float32) {
	if h > r.H {
		r.H = h
	}
}

func (r *HybridLayout) UpdatePosition(pos [4]float32) {
	r.X = pos[0]
	r.Y = pos[1]

	r.CursorX = pos[0]
	r.CursorY = pos[1]

	r.InitY = pos[1]
}

func (r *HybridLayout) BoundingBox() [4]float32 {
	return [4]float32{}
}
