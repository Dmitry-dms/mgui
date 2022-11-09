package widgets

type Widget interface {
	WidgetId() string
	UpdatePosition([4]float32)
	Height() float32
	Width() float32
	BoundingBox() [4]float32
	ToggleUpdate()
}
type PaddingType uint

const (
	LeftPadding PaddingType = 1 << iota
	TopPadding
	RightPadding
	BotPadding
	AllPadding = LeftPadding | TopPadding | RightPadding | BotPadding
)

type baseWidget struct {
	id              string
	boundingBox     [4]float32
	BackgroundColor [4]float32
	LastVert        []float32
	LastInd         []int32
	Last            int
	LastVertCount   int
	Updated         bool
}

func (b *baseWidget) ToggleUpdate() {
	b.Updated = true
}
func (b *baseWidget) height() float32 {
	return b.boundingBox[3]
}
func (b *baseWidget) updatePosition(p [4]float32) {
	if p[0] != b.boundingBox[0] || p[1] != b.boundingBox[1] ||
		p[2] != b.width() || p[3] != b.height() {
		b.Updated = true
	}
	b.boundingBox = p
}
func (b *baseWidget) width() float32 {
	return b.boundingBox[2]
}
