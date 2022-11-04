package widgets

type Image struct {
	//Texture *gogl.Texture
	baseWidget
	//Id      string
	//BoundingBox  [4]float32 //x,y,w,h
	//CurrentColor [4]float32
}

func NewImage(id string, x, y, w, h float32, clr [4]float32) *Image {
	i := Image{
		//Texture: tex,
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			backgroundColor: clr,
		},
	}
	return &i
}

func (i *Image) BoundingBox() [4]float32 {
	return i.boundingBox
}
func (i *Image) UpdatePosition(pos [4]float32) {
	//i.BoundingBox = pos
	i.updatePosition(pos)
}
func (i *Image) Color() [4]float32 {
	return i.backgroundColor
}
func (i *Image) WidgetId() string {
	return i.id
}

func (i *Image) Visible() bool {
	return true
}

func (i *Image) Height() float32 {
	return i.height()
}
func (i *Image) Width() float32 {
	return i.width()
}
