package widgets

type Rectangle struct {
	baseWidget
}

func NewRectangle(id string, x, y, w, h float32, clr [4]float32) *Rectangle {
	r := Rectangle{
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			BackgroundColor: clr,
			LastVert:        nil,
			LastInd:         nil,
			Last:            0,
			LastVertCount:   0,
			Updated:         true,
		},
	}
	return &r
}

func (r *Rectangle) WidgetId() string {
	return r.baseWidget.id
}

func (r *Rectangle) UpdatePosition(pos [4]float32) {
	r.baseWidget.updatePosition(pos)
}

func (r *Rectangle) Height() float32 {
	return r.baseWidget.height()
}

func (r *Rectangle) Width() float32 {
	return r.baseWidget.width()
}

func (r *Rectangle) BoundingBox() [4]float32 {
	return r.baseWidget.boundingBox
}

func (r *Rectangle) ToggleUpdate() {
	r.baseWidget.ToggleUpdate()
}
