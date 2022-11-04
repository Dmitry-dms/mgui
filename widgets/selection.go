package widgets

type Selection struct {
	baseWidget
	CurrentItem int
	Opened      bool
}

func NewSelection(id string, x, y, w, h float32) *Selection {
	s := Selection{
		baseWidget: baseWidget{
			id:              id,
			boundingBox:     [4]float32{x, y, w, h},
			backgroundColor: [4]float32{},
		},
		CurrentItem: 0,
	}
	return &s
}

func (s *Selection) WidgetId() string {
	return s.id
}

func (s *Selection) UpdatePosition(pos [4]float32) {
	s.updatePosition(pos)
}

func (s *Selection) Height() float32 {
	return s.height()
}

func (s *Selection) Width() float32 {
	return s.width()
}

func (s *Selection) BoundingBox() [4]float32 {
	return s.boundingBox
}
