package ui

type Toolbar struct {
	x, y float32 // top-left corner
	w, h float32
	clr  [4]float32
}

func NewToolbar(x, y, w, h float32) Toolbar {
	tb := Toolbar{
		x:   x,
		y:   y,
		w:   w,
		h:   h,
		clr: [4]float32{35, 53, 79, 1},
	}
	return tb
}
