package ui

import "github.com/Dmitry-dms/mgui/utils"

type Scrollbar struct {
	x, y, w, h     float32
	bX, bY, bW, bH float32
	clr            [4]float32
	isActive       bool
}

func NewScrolBar(bound, btn utils.Rect, clr [4]float32) *Scrollbar {
	sb := Scrollbar{
		x:        bound.Min.X,
		y:        bound.Min.Y,
		w:        bound.Width(),
		h:        bound.Height(),
		bX:       btn.Min.X,
		bY:       btn.Min.Y,
		bW:       btn.Width(),
		bH:       btn.Height(),
		clr:      clr,
		isActive: false,
	}
	return &sb
}
