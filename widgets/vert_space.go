package widgets

type VSpace struct {
	baseWidget
	//BoundingBox [4]float32 //x,y,w,h
	//Id          string
}

func NewVertSpace(id string, box [4]float32) *VSpace {
	vs := VSpace{
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     box,
			backgroundColor: [4]float32{255, 255, 255, 1},
		},
	}
	return &vs
}

func (s *VSpace) UpdatePosition(pos [4]float32) {
	s.updatePosition(pos)
}

func (s *VSpace) BoundingBox() [4]float32 {
	return s.boundingBox
}
func (s *VSpace) Visible() bool {
	return true
}
func (s *VSpace) Height() float32 {
	return s.height()
}
func (s *VSpace) Width() float32 {
	return s.width()
}

func (s *VSpace) Color() [4]float32 {
	return s.backgroundColor
}
func (s *VSpace) WidgetId() string {
	return s.id
}
