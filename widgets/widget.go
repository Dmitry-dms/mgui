package widgets

type Widget interface {
	Id() string
	UpdatePosition([4]float32)
	Height() float32
	Width() float32
	BoundingBox() [4]float32
	ToggleUpdate()
	RenderInfo() ([]float32, []int32, int, int)
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

	Vertices []float32
	Indices  []int32
	// LastBufferIndex shows the value of the Buffer.LastInd counter after
	// the widget has been sent to the Buffer. It prevents errors from appearing when constructing the indices Buffer.
	// For example, the first widget was drawn, and in the next frame, a second widget was drawn in front of
	// this widget, which indicates that the first widget needs to be redrawn.
	// But at the stage of drawing the second widget, there is no access to the Updated flag of the first one.
	// Therefore, it is necessary to monitor the status of the indexes Buffer.
	LastBufferIndex int
	VertCount       int
	Updated         bool
}

func (b *baseWidget) ToggleUpdate() {
	b.Updated = true
}

func (b *baseWidget) RenderInfo() ([]float32, []int32, int, int) {
	return b.Vertices, b.Indices, b.VertCount, b.LastBufferIndex
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
