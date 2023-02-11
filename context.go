package ui

import (
	"github.com/Dmitry-dms/mgui/cache"
	"github.com/Dmitry-dms/mgui/draw"
	"github.com/Dmitry-dms/mgui/fonts"
	"image"

	"github.com/Dmitry-dms/mgui/styles"
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/Dmitry-dms/mgui/widgets"
)

var uiCtx *UiContext

func init() {
	uiCtx = NewContext(nil)
}
func ctx() *UiContext {
	return uiCtx
}

type UiContext struct {
	// rq *RenderQueue

	renderer UiRenderer
	io       *Io

	//Widgets
	Windows          []*Window
	sortedWindows    []*Window
	ActiveWidget     string
	LastActiveWidget string

	ActiveWindow *Window

	//widg space
	ActiveWidgetSpaceId              string
	WantScrollFocusWidgetSpaceId     string
	WantScrollFocusWidgetSpaceLastId string
	ActiveWidgetSpace                *WidgetSpace
	FocusedWidgetSpace               *WidgetSpace

	// For now, it uses as mark which tells there is something that must be handled inside global buffer.
	// For example GlobalWidgetSpace.
	CurrentGlobalWidgetSpace *WidgetSpace

	//SelectableText *widgets.Text
	SelectedText string

	SelectedTexts     []*widgets.Text
	SelectedTextStart *widgets.Text
	SelectedTextEnd   *widgets.Text

	FocusedTextInput *widgets.Text

	PriorWindow       *Window
	HoveredWindow     *Window
	LastHoveredWindow *Window
	WindowCounter     int // Shows how much active windows on the screen

	displaySize [2]float32

	//cache
	windowCache *cache.RamCache[*Window]
	windowStack Stack
	//widgetsCache   *cache.RamCache[widgets.Widget]
	widgetsCache   *widgets.WidgetsCache
	widgSpaceCache *cache.RamCache[*WidgetSpace]

	//refactor
	Time float32

	//style
	CurrentStyle       *styles.Style
	prevStyle          *styles.Style
	styleChangeCounter int
	StyleChanged       bool

	//fonts
	font *fonts.Font

	globalBuffer   *draw.CmdBuffer
	delayedWidgets []func()
}

const (
	Selectable      = widgets.Selectable
	SplitWords      = widgets.SplitWords
	SplitChars      = widgets.SplitChars
	FitContent      = widgets.FitContent
	Editable        = widgets.Editable
	DefaultTextFlag = widgets.Default
)

func NewContext(frontRenderer UiRenderer) *UiContext {
	c := UiContext{
		// rq:            NewRenderQueue(),
		renderer:      frontRenderer,
		io:            NewIo(),
		Windows:       make([]*Window, 0),
		sortedWindows: make([]*Window, 0),
		windowCache:   cache.NewRamCache[*Window](),
		//widgetsCache:   cache.NewRamCache[widgets.Widget](),
		widgetsCache:   widgets.New(),
		windowStack:    NewStack(),
		widgSpaceCache: cache.NewRamCache[*WidgetSpace](),
		CurrentStyle:   &styles.DefaultStyle,
		SelectedTexts:  []*widgets.Text{},
	}
	c.globalBuffer = draw.NewBuffer(c.io.DisplaySize)

	return &c
}

func UploadFont(path string, size int, dpi float32, from int, to int) (*fonts.Font, *image.RGBA) {
	c := ctx()
	f, data := fonts.NewFont(path, size, dpi, from, to)
	c.font = f
	return f, data
}

func AddRenderer(frontRenderer UiRenderer) {
	c := ctx()
	c.renderer = frontRenderer
}

func SetChangeCursorFunc(f func(c CursorType)) {
	c := ctx()
	c.io.SetCursor = f
}

func (c *UiContext) getPeekWindow() *Window {
	if c.windowStack.Length() != 0 {
		return c.windowStack.Peek()
	}
	return nil
}

func (c *UiContext) dragBehaviorInWindow(rect utils.Rect, captured *bool) {
	if !*captured {
		*captured = utils.PointInRect(c.io.MousePos, rect) && c.io.DragStarted(rect) && c.io.IsDragging
		if c.ActiveWindow != nil {
			// FIXME(@Dmitry-dms): The next frame window moves a bit.
			c.ActiveWindow.capturedInsideWin = *captured
		}
	} else {
	}
	if c.io.MouseReleased[0] {
		*captured = false
		if c.ActiveWindow != nil {
			c.ActiveWindow.capturedInsideWin = false
		}
	}
}

func (c *UiContext) dragBehavior(rect utils.Rect, captured *bool) {
	if !*captured {
		*captured = utils.PointInRect(c.io.MousePos, rect) && c.io.DragStarted(rect) && c.io.IsDragging
	}
	if c.io.MouseReleased[0] {
		*captured = false
	}
}

func (c *UiContext) SetScrollY(scrollY float32) {
	wnd := c.windowStack.Peek()
	wnd.currentWidgetSpace.setScrollY(scrollY)
}

func (c *UiContext) AddWidget(id string, w widgets.Widget) bool {
	return c.widgetsCache.Add(id, w)
}

func (c *UiContext) GetWidget(id string) (widgets.Widget, bool) {
	return c.widgetsCache.Get(id)
}

func GetIo() *Io {
	c := ctx()
	return c.io
}

func NewFrame(displaySize [2]float32) {
	c := ctx()
	c.UpdateMouseInputs()

	c.renderer.NewFrame()

	c.io.SetDisplaySize(displaySize[0], displaySize[1])

}
func (c *UiContext) pushWindowFront(w *Window) {
	for i := len(c.sortedWindows) - 1; i >= 0; i-- {
		if c.sortedWindows[i] == w {
			if i == len(c.sortedWindows)-1 {
				return
			}
			c.sortedWindows[i] = c.sortedWindows[len(c.sortedWindows)-1]
			c.sortedWindows[len(c.sortedWindows)-1] = w
			return
		}
	}
}

func (c *UiContext) findHoveredWindow() {
	var hovered *Window
	if len(c.sortedWindows) == 0 {
		return
	}

	for i := 0; i <= len(c.sortedWindows)-1; i++ {
		window := c.sortedWindows[i]
		bb := window.outerRect

		if !bb.Contains(c.io.MousePos) {
			continue
		}
		if c.io.MouseClicked[0] && c.ActiveWindow != window {
			if !utils.PointInRect(c.io.MousePos, c.ActiveWindow.outerRect) {
				c.ActiveWindow = window
				c.pushWindowFront(window)
			}
		}

		if c.ActiveWindow == window {
			hovered = window
		} else {
			hovered = c.LastHoveredWindow
		}
		c.LastHoveredWindow = window

	}
	if c.ActiveWindow == nil {
		c.ActiveWindow = c.sortedWindows[len(c.sortedWindows)-1]
	}
	c.HoveredWindow = hovered

}

func (c *UiContext) UpdateMouseInputs() {

	io := c.io

	if io.IsMousePosValid(&io.MousePos) && io.IsMousePosValid(&io.MousePosPrev) {
		io.MouseDelta = io.MousePos.Sub(io.MousePosPrev)
	} else {
		io.MouseDelta = utils.Vec2{0, 0}
	}

	io.MousePosPrev = io.MousePos
	for i := 0; i < len(io.MouseDown); i++ {
		io.MouseClicked[i] = io.MouseDown[i] && io.MouseDownDuration[i] < 0
		io.MouseClickedCount[i] = 0
		io.MouseReleased[i] = !io.MouseDown[i] && io.MouseDownDuration[i] >= 0
		io.MouseDownDurationPrev[i] = io.MouseDownDuration[i]
		if io.MouseDown[i] {
			if io.MouseDownDuration[i] < 0 {
				io.MouseDownDuration[i] = 0
			} else {
				io.MouseDownDuration[i] += io.DeltaTime
			}
		} else {
			io.MouseDownDuration[i] = -1
		}

		if io.MouseClicked[i] {
			isRepeatedClick := false
			if c.Time-float32(io.MouseClickedTime[i]) < io.MouseDoubleClickTime {
				var delta utils.Vec2
				if io.IsMousePosValid(&io.MousePos) {
					delta = io.MousePos.Sub(io.MouseClickedPos[i])
				} else {
					delta = utils.Vec2{0, 0}
				}

				if delta.LengthSqr() < io.MouseDoubleClickMaxDist*io.MouseDoubleClickMaxDist {
					isRepeatedClick = true
				}
			}

			if isRepeatedClick {
				io.MouseClickedLastCount[i]++
			} else {
				io.MouseClickedLastCount[i] = 1
			}

			io.MouseClickedTime[i] = c.Time
			io.MouseClickedPos[i] = io.MousePos
			io.MouseClickedCount[i] = io.MouseClickedLastCount[i]
			io.MouseDragMaxDistanceSqr[i] = 0
		} else if io.MouseDown[i] {
			// Maintain the maximum distance we reaching from the initial click position, which is used with dragging threshold
			var deltaSqrPos float32
			if io.IsMousePosValid(&io.MousePos) {
				deltaSqrPos = (io.MousePos.Sub(io.MouseClickedPos[i])).LengthSqr()
			} else {
				deltaSqrPos = 0
			}
			io.MouseDragMaxDistanceSqr[i] = utils.Max(io.MouseDragMaxDistanceSqr[i], deltaSqrPos)
		}
		// We provide io.MouseDoubleClicked[] as a legacy service
		io.MouseDoubleClicked[i] = (io.MouseClickedCount[i] == 2)

	}

}

func copyWindows(w []*Window) []*Window {
	r := make([]*Window, len(w))
	for i, v := range w {
		r[i] = v
	}
	return r
}

var lastWinL = 0

func EndFrame(size [2]float32) {
	c := ctx()
	//if c.WindowCounter == 0 {
	//	return
	//}
	// Если количество окон не изменилось, в копировании нет нужды
	if lastWinL != len(c.Windows) {
		c.sortedWindows = copyWindows(c.Windows)
		lastWinL = len(c.Windows)
	}

	c.findHoveredWindow()
	//if len(c.sortedWindows) == 0 {
	//	return
	//}

	for _, widget := range c.delayedWidgets {
		widget()
	}

	for _, v := range c.sortedWindows {
		c.renderer.Draw(size, *v.buffer)
		v.buffer.Clear()
	}
	//fmt.Println(len(c.globalBuffer.DrawCalls))
	if len(c.globalBuffer.DrawCalls) != 0 {
		c.renderer.Draw(size, *c.globalBuffer)
		c.globalBuffer.Clear()
	}

	// c.renderer.End()

	//if !c.io.IsDragging && c.wantResizeH == true {
	//	c.wantResizeH = false
	//
	//} else if !c.io.IsDragging && c.wantResizeV == true {
	//	c.wantResizeV = false
	//}

	c.io.ScrollX = 0
	c.io.ScrollY = 0

	c.LastActiveWidget = c.ActiveWidget
	c.ActiveWidget = ""

	if c.ActiveWindow != nil {
	} else {
		c.ActiveWidgetSpaceId = ""
	}

	// Используется для удаления фокуса, если ЛКМ была кликнута снаружи WidgetSpace
	if c.FocusedWidgetSpace != nil {
		if utils.PointOutsideRect(c.io.MouseClickedPos[0], utils.NewRectS(c.FocusedWidgetSpace.ClipRect)) {
			c.FocusedWidgetSpace = nil
		}
	}

	//if c.FocusedTextInput != nil {
	//	if utils.PointOutsideRect(c.io.MouseClickedPos[0], utils.NewRectS(c.FocusedTextInput.BoundingBox())) {
	//		ToggleAllWidgets()
	//		fmt.Println("outside - ", c.FocusedTextInput.WidgetId())
	//		c.FocusedTextInput = nil
	//	}
	//}

	c.WantScrollFocusWidgetSpaceLastId = c.WantScrollFocusWidgetSpaceId
	c.WantScrollFocusWidgetSpaceId = ""

	c.io.MouseClickedPos[0] = utils.Vec2{}
	//c.ActiveWindow.widgSpaces = []*WidgetSpace{}

	c.io.PressedKey = GuiKey_None
	c.io.modPressed = [8]bool{}
	c.io.KeyPressedThisFrame = false

	c.WindowCounter = 0
}

func SetDisplaySize(w, h float32) {
	c := ctx()
	c.io.SetDisplaySize(w, h)
	// Need to redraw all widgets
	ToggleAllWidgets()
}

// ToggleAllWidgets is responsible for setting the Update=true flag for all widgets.
// TODO(@Dmitry-dms): Should it be only visible widgets?
func ToggleAllWidgets() {
	c := ctx()
	for _, i := range c.widgetsCache.Map() {
		i.ToggleUpdate()
	}
}

func GET_CONTEXT() *UiContext {
	return ctx()
}

type StyleVar4f uint
type StyleVar1f uint

const (
	ButtonActiveColor StyleVar4f = iota
	ButtonHoveredColor
	Padding
)

const (
	Margin StyleVar1f = iota
	FontScale
	AllPadding
	LeftPadding
	TopPadding
	RightPadding
	BottomPadding
)

func PushStyleVar1f(v StyleVar1f, val float32) {
	c := ctx()
	if c.styleChangeCounter == 0 {
		prev := *c.CurrentStyle
		c.prevStyle = &prev
	}
	switch v {
	case Margin:
		c.CurrentStyle.Margin = val
	case FontScale:
		c.CurrentStyle.FontScale = val
	case LeftPadding:
		c.CurrentStyle.Padding.Left = val
	case RightPadding:
		c.CurrentStyle.Padding.Right = val
	case TopPadding:
		c.CurrentStyle.Padding.Top = val
	case BottomPadding:
		c.CurrentStyle.Padding.Bottom = val
	case AllPadding:
		c.CurrentStyle.Padding = styles.Padding{
			Left:   val,
			Top:    val,
			Right:  val,
			Bottom: val,
		}
	}
	c.StyleChanged = true
	c.styleChangeCounter++
}

func (c *UiContext) PushStyleVar4f(v StyleVar4f, m [4]float32) {
	if c.styleChangeCounter == 0 {
		prev := *c.CurrentStyle
		c.prevStyle = &prev
	}
	switch v {
	case ButtonHoveredColor:
		c.CurrentStyle.BtnHoveredColor = m
	case ButtonActiveColor:
		c.CurrentStyle.BtnActiveColor = m
	}
	c.StyleChanged = true
	c.styleChangeCounter++
}
func PopStyleVar() {
	c := ctx()
	c.CurrentStyle = c.prevStyle
	c.styleChangeCounter = 0
}
func (c *UiContext) LineArc() {
	//wnd := c.windowStack.Peek()
	//x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	//wnd.buffer.CreateBezierQuad(300, 500, 20, 150, 300, 300, [4]float32{255, 0, 0, 1}, wnd.DefaultClip())
	//wnd.buffer.RoundedBorderRectangle(x, y, 200, 100, 30, 15, red, wnd.DefaultClip())
}

func (c *UiContext) Bezier() {
	wnd := c.windowStack.Peek()
	wnd.buffer.CreateBezierQuad(300, 500, 20, 150, 300, 300, 20, [4]float32{255, 0, 0, 1})
}
func (c *UiContext) Line(end float32) {
	wnd := c.windowStack.Peek()
	wnd.buffer.CreateLine(0, 0, end, end, [4]float32{255, 0, 0, 1})

}

func (c *UiContext) SeparateBuffer(wnd *Window, texid float32, clip draw.ClipRectCompose) {

}

type UiRenderer interface {
	NewFrame()
	Draw(displaySize [2]float32, buffer draw.CmdBuffer)
}
