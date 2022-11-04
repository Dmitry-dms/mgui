package widgets

type Widget interface {
	WidgetId() string
	UpdatePosition([4]float32)
	Height() float32
	Width() float32
	BoundingBox() [4]float32
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
	backgroundColor [4]float32
}

func (b *baseWidget) height() float32 {
	return b.boundingBox[3]
}
func (b *baseWidget) updatePosition(p [4]float32) {
	b.boundingBox = p
}
func (b *baseWidget) width() float32 {
	return b.boundingBox[2]
}
