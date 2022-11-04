package widgets

// TODO: add Base Widget
type Button struct {
	IsActive bool
	//CurrentColor [4]float32
	//Id           string
	//BoundingBox  [4]float32 //x,y,w,h
	base baseWidget
}

func NewButton(id string, x, y, w, h float32, backClr [4]float32) *Button {
	btn := Button{
		base: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			backgroundColor: backClr,
		},
		IsActive: false,
		//CurrentColor: backClr,
		//Id:           id,
		//BoundingBox:  [4]float32{x, y, w, h},
	}
	return &btn
}
func (b *Button) UpdatePosition(pos [4]float32) {
	//b.BoundingBox = pos
	b.base.updatePosition(pos)
}

func (b *Button) ChangeActive() {
	b.IsActive = !b.IsActive
}

func (b *Button) SetWidth(w float32) {
	//b.BoundingBox[2] = w
	b.base.boundingBox[2] = w
}
func (b *Button) SetHeight(h float32) {
	b.base.boundingBox[3] = h
}
func (b *Button) WidgetId() string {
	return b.base.id
}

func (b *Button) Height() float32 {
	return b.base.height()
}
func (b *Button) Width() float32 {
	return b.base.width()
}

func (b *Button) BoundingBox() [4]float32 {
	return b.base.boundingBox
}
func (b *Button) Color() [4]float32 {
	return b.base.backgroundColor
}

func (b *Button) SetColor(clr [4]float32) {
	//b.CurrentColor = clr
	b.base.backgroundColor = clr
}
