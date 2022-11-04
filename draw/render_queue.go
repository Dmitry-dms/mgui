package draw

import (
	// "fmt"

	"github.com/Dmitry-dms/mgui/fonts"
	//"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/Dmitry-dms/mgui/widgets"
)

type Command struct {
	Type     CmdType
	WidgetId string

	sb *separate_buff

	Shown    bool
	Rect     *Rect_command
	triangle *triangle_command
	RRect    *Rounded_rect
	Window   *Window_command
	WinStart *Window_start_command
	Toolbar  *Toolbar_command
	Text     *Text_command
	Bezier   *BezierQuad_command
	Line     *Line_command
}

type separate_buff struct {
	texid        uint32
	clipRect     [4]float32
	mainClipRect [4]float32
}

type Text_command struct {
	Widget  *widgets.Text
	Text    string
	Font    fonts.Font
	X, Y    float32
	Clr     [4]float32
	Id      string
	Padding int
	Scale   float32
}

type Rect_command struct {
	X, Y, W, H float32
	Clr        [4]float32
	TexId      uint32
	radius     int
	shape      RoundedRectShape
	coords     [4]float32
	// ScaleFactor float32
	// Texture     *gogl.Texture
}
type triangle_command struct {
	x0, y0, x1, y1, x2, y2 float32
	clr                    [4]float32
}

type Window_start_command struct {
	id   string
	x, y float32
}
type Window_command struct {
	Active     bool
	Id         string
	X, Y, W, H float32
	Clr        [4]float32
	Toolbar    Toolbar_command
	Scrollbar  Scrollbar_command
}
type Toolbar_command struct {
	X, Y, W, H float32
	Clr        [4]float32
}

type Scrollbar_command struct {
	X, Y, W, H     float32
	Xb, Yb, Wb, Hb float32
	Radius         int
	ScrollClr      [4]float32
	BtnClr         [4]float32
}
type Rounded_rect struct {
	X, Y, W, H float32
	Clr        [4]float32
	Radius     int
}
type Line_command struct {
	StartX, StartY float32
	EndX, EndY     float32
	Clr            [4]float32
	Points         []utils.Vec2
}
type BezierQuad_command struct {
	StartX, StartY     float32
	SupportX, SupportY float32
	EndX, EndY         float32
	Steps              float32
	Clr                [4]float32
}

type CmdType int

const (
	RectType CmdType = iota
	SeparateBuffer
	BezierQuad
	RectTypeT
	Triangle
	TriangleT
	Line
	LineStrip
	Circle
	RoundedRect
	RoundedRectT
	WindowCmd
	ToolbarCmd
	WindowStartCmd
	ScrollbarCmd
	ScrollButtonCmd
	Text
)
