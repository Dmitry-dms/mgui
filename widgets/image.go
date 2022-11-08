package widgets

type Image struct {
	//Texture *gogl.Texture
	baseWidget
	TexId     uint32
	TexCoords [4]float32
	//Id      string
	//BoundingBox  [4]float32 //x,y,w,h
	//CurrentColor [4]float32
}

func NewImage2(id string, x, y, w, h float32, texid uint32, texcoords, clr [4]float32) *Image {
	i := Image{
		//Texture: tex,
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			BackgroundColor: clr,
		},
		TexCoords: texcoords,
		TexId:     texid,
	}
	return &i
}
func NewImage(id string, x, y, w, h float32, clr [4]float32) *Image {
	i := Image{
		//Texture: tex,
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			BackgroundColor: clr,
		},
	}
	return &i
}

//func (i *Image) WasUpdated() bool {
//	return i.WasUpdated()
//}

func (i *Image) SetColor(clr [4]float32) {
	if clr == i.BackgroundColor {
		return
	}
	i.BackgroundColor = clr
	i.Updated = true
}

func (i *Image) BoundingBox() [4]float32 {
	return i.boundingBox
}
func (i *Image) UpdatePosition(pos [4]float32) {
	//i.BoundingBox = pos
	i.updatePosition(pos)
	//i.Updated = false
}
func (i *Image) Color() [4]float32 {
	return i.BackgroundColor
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
