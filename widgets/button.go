package widgets

// TODO: add Base Widget
type Button struct {
	IsActive bool
	//CurrentColor [4]float32
	//Id           string
	//BoundingBox  [4]float32 //x,y,w,h
	baseWidget
}

func NewButton(id string, x, y, w, h float32, backClr [4]float32) *Button {
	btn := Button{
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			BackgroundColor: backClr,
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
	b.updatePosition(pos)
}

func (b *Button) ChangeActive() {
	b.IsActive = !b.IsActive
}

func (b *Button) SetWidth(w float32) {
	if w == b.width() {
		return
	}
	b.boundingBox[2] = w
	b.ToggleUpdate()
}
func (b *Button) SetHeight(h float32) {
	b.boundingBox[3] = h
}
func (b *Button) Id() string {
	return b.id
}

func (b *Button) Height() float32 {
	return b.height()
}
func (b *Button) Width() float32 {
	return b.width()
}

func (b *Button) BoundingBox() [4]float32 {
	return b.boundingBox
}
func (b *Button) Color() [4]float32 {
	return b.BackgroundColor
}

func (b *Button) SetColor(clr [4]float32) {
	if clr == b.BackgroundColor {
		return
	}
	b.ToggleUpdate()
	b.BackgroundColor = clr
}
