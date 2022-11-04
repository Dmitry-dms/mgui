package widgets

type TabBar struct {
	baseWidget
	Bars       []*TabItem
	CurrentTab int
	BarHeight  float32

	CursorX float32
}

func NewTabBar(id string, x, y, w, h float32) *TabBar {
	t := TabBar{baseWidget: baseWidget{
		id:              id,
		boundingBox:     [4]float32{x, y, w, h},
		backgroundColor: [4]float32{255, 255, 255, 1},
	},
		Bars:    []*TabItem{},
		CursorX: x}
	return &t
}
func (t *TabBar) FindTabItem(name, widgSpaceId string) (item *TabItem, index int) {
	for i, it := range t.Bars {
		if it.Name == name {
			item = it
			index = i
			return
		}
	}
	item = &TabItem{
		WidgetSpaceId: widgSpaceId,
		Name:          name,
	}
	t.Bars = append(t.Bars, item)
	index = len(t.Bars) - 1
	return
}

func (t *TabBar) WidgetId() string {
	return t.id
}

func (t *TabBar) UpdatePosition(pos [4]float32) {
	t.updatePosition(pos)
	t.CursorX = pos[0]
}
func (t *TabBar) ChangeActive(b *TabItem) {
	for _, bar := range t.Bars {
		if bar != b {
			bar.Active = false
		} else {
			bar.Active = true
		}
	}
}
func (t *TabBar) SetHeight(h float32) {
	//if h > t.height() {
	t.boundingBox[3] = h
	//}
}

func (t *TabBar) SetWidth(w float32) {
	if w > t.width() {
		t.boundingBox[2] = w
	}
}
func (t *TabBar) Height() float32 {
	return t.height()
}

func (t *TabBar) Width() float32 {
	return t.width()
}

func (t *TabBar) BoundingBox() [4]float32 {
	return t.boundingBox
}
