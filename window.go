package ui

import (
	// "fmt"
	"fmt"
	"github.com/Dmitry-dms/mgui/fonts"
	"math"
	"strings"
	"time"

	"math/rand"

	// "math/rand"

	//"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/mgui/draw"
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/Dmitry-dms/mgui/widgets"
)

type Window struct {
	toolbar Toolbar
	x, y    float32 // top-left corner
	w, h    float32
	active  bool
	Id      string

	outerRect  utils.Rect
	minW, minH float32

	//render
	buffer *draw.CmdBuffer

	mainWidgetSpace    *WidgetSpace
	currentWidgetSpace *WidgetSpace
	focusedWidgetSpace *WidgetSpace

	//widgSpaceWantFocus bool

	widgSpaces []*WidgetSpace

	capturedV, capturedH  bool
	capturedWin           bool
	capturedInsideWin     bool
	capturedTextSelection bool

	delayedWidgets []func()

	VisibleTexts []*widgets.Text
	textRegions  []utils.Rect

	//intent
	wantResizeH, wantResizeV bool
}

func genWindowId() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprint(rand.Intn(100000))
}

func NewCustomWindow(x, y, w, h float32) *Window {
	id := genWindowId()
	wnd := Window{
		Id:              id,
		toolbar:         Toolbar{},
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		outerRect:       utils.Rect{Min: utils.Vec2{X: x, Y: y}, Max: utils.Vec2{X: x + w, Y: y + h}},
		minW:            200,
		minH:            50,
		mainWidgetSpace: newWidgetSpace(fmt.Sprintf("main-widg-space-%s", id), x, y, w, h, Default),
		buffer:          draw.NewBuffer(uiCtx.io.DisplaySize),
		widgSpaces:      make([]*WidgetSpace, 0),
		delayedWidgets:  []func(){},
		VisibleTexts:    []*widgets.Text{},
	}

	wnd.currentWidgetSpace = wnd.mainWidgetSpace
	return &wnd
}

func NewWindow(x, y, w, h float32) *Window {
	tb := NewToolbar(x, y, w, 30)
	id := genWindowId()

	wnd := Window{
		Id:              id,
		toolbar:         tb,
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		outerRect:       utils.Rect{Min: utils.Vec2{X: x, Y: y}, Max: utils.Vec2{X: x + w, Y: y + h}},
		minW:            200,
		minH:            50,
		mainWidgetSpace: newWidgetSpace(fmt.Sprintf("main-widg-space-%s", id), x, y+tb.h, w, h-tb.h, Default),
		buffer:          draw.NewBuffer(uiCtx.io.DisplaySize),
		widgSpaces:      make([]*WidgetSpace, 0),
		delayedWidgets:  []func(){},
		VisibleTexts:    []*widgets.Text{},
	}

	wnd.currentWidgetSpace = wnd.mainWidgetSpace
	return &wnd
}

const (
	defx, defy, defw, defh = 300, 100, 900, 500
	scrollChange           = 2
)

func BeginCustomWindow(id string, x, y, w, h, wsX, wsY, wsW, wsH float32, texId uint32, textCoords [4]float32, widgFunc func()) {
	c := ctx()
	var wnd *Window
	wnd, ok := c.windowCache.Get(id)
	if !ok {
		wnd = NewCustomWindow(x, y, w, h)
		c.Windows = append(c.Windows, wnd)
		wnd.Id = id
		c.windowCache.Add(id, wnd)
	}
	newX := wnd.x
	newY := wnd.y
	newH := wnd.h
	newW := wnd.w

	// logic
	{
		//Прямоугольник справа
		//vResizeRect := utils.NewRect(wnd.x+wnd.w-scrollChange, wnd.y, scrollChange+5, wnd.h)
		//hResizeRect := utils.NewRect(wnd.x, wnd.y+wnd.h-scrollChange, wnd.w, scrollChange+5)
		//if utils.PointInRect(c.io.MousePos, hResizeRect) && c.ActiveWindow == wnd {
		//	c.io.SetCursor(VResizeCursor)
		//	c.wantResizeH = true
		//} else if utils.PointInRect(c.io.MousePos, vResizeRect) && c.ActiveWindow == wnd {
		//	c.io.SetCursor(HResizeCursor)
		//	c.wantResizeV = true
		//} else {
		//	c.io.SetCursor(ArrowCursor)
		//}
		//c.dragBehavior(vResizeRect, &wnd.capturedV)
		//c.dragBehavior(hResizeRect, &wnd.capturedH)
		//// Изменение размеров окна
		//if c.wantResizeH && c.ActiveWindow == wnd && wnd.capturedH {
		//	n := newH
		//	n += c.io.MouseDelta.Y
		//	if n > wnd.minH {
		//		newH = n
		//		if wnd.mainWidgetSpace.scrlY != 0 {
		//			wnd.mainWidgetSpace.scrlY -= c.io.MouseDelta.Y
		//		}
		//	}
		//} else if c.wantResizeV && c.ActiveWindow == wnd && wnd.capturedV {
		//	n := newW
		//	n += c.io.MouseDelta.X
		//	if n > wnd.minW {
		//		newW = n
		//	}
		//}

		c.dragBehavior(wnd.outerRect, &wnd.capturedWin)
		// Изменение положения окна
		if c.ActiveWindow == wnd && wnd.capturedWin && !wnd.wantResizeV &&
			!wnd.wantResizeH && !wnd.capturedInsideWin && !wnd.capturedTextSelection {
			newX += c.io.MouseDelta.X
			newY += c.io.MouseDelta.Y
		}
	}

	wnd.x = newX
	wnd.y = newY
	wnd.h = newH
	wnd.w = newW

	wnd.outerRect = utils.NewRect(wnd.x, wnd.y, wnd.w-wnd.mainWidgetSpace.verticalScrollbar.w, wnd.h)

	wnd.mainWidgetSpace.X = wnd.x + wsX
	wnd.mainWidgetSpace.Y = wnd.y + wsY
	wnd.mainWidgetSpace.W = wsW
	wnd.mainWidgetSpace.H = wsH

	wnd.mainWidgetSpace.cursorX = wnd.mainWidgetSpace.X
	wnd.mainWidgetSpace.cursorY = wnd.mainWidgetSpace.Y

	wnd.mainWidgetSpace.ClipRect = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y,
		wnd.mainWidgetSpace.W - wnd.mainWidgetSpace.verticalScrollbar.w, wnd.mainWidgetSpace.H}

	//wnd.buffer.CreateWindow(cmdw, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
	//wnd.widgetSpaceLogic(wnd.mainWidgetSpace, func() draw.ClipRectCompose {
	//	cl := [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}
	//	return draw.NewClip(draw.EmptyClip, cl)
	//})

	//// Draw selected text regions. We doo it here because we don't want to draw it in front of text.
	//// Maybe in future I will change text selection algorithm and rework this.
	//for _, region := range wnd.textRegions {
	//	b := region.Min
	//	wnd.buffer.CreateRect(b.X, b.Y, region.Width(), region.Height(), 0, draw.StraightCorners, 0,
	//		softGreen, wnd.DefaultClip())
	//}

	DrawImage(wnd.buffer, "custom-wnd-back-"+id, wnd.x, wnd.y, wnd.w, wnd.h, texId,
		textCoords, whiteColor, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
	c.windowStack.Push(wnd)

	widgFunc()

	for _, f := range wnd.delayedWidgets {
		f()
	}
	wnd.delayedWidgets = []func(){}

	wnd.VisibleTexts = []*widgets.Text{}
	wnd = c.windowStack.Pop()

	wnd.mainWidgetSpace.checkVerScroll()
	var clip = draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect)
	wnd.buffer.SeparateBuffer(0, clip) // Make sure that we didn't miss anything
	wnd.mainWidgetSpace.AddVirtualHeight(c.CurrentStyle.BotMargin)

	wnd.mainWidgetSpace.lastVirtualHeight = wnd.mainWidgetSpace.virtualHeight
	wnd.mainWidgetSpace.virtualHeight = 0
	wnd.mainWidgetSpace.lastVirtualWidth = wnd.mainWidgetSpace.virtualWidth
	wnd.mainWidgetSpace.virtualWidth = 0
}

func BeginWindow(windowName string, opened *bool) {
	if *opened == false {
		return
	}
	c := ctx()

	wnd, ok := c.windowCache.Get(windowName)
	if !ok {
		r := rand.Intn(500)
		g := rand.Intn(300)
		wnd = NewWindow(defx+float32(r), defy+float32(g), defw, defh)
		c.Windows = append(c.Windows, wnd)
		wnd.Id = windowName
		c.windowCache.Add(windowName, wnd)
	}
	c.WindowCounter++

	newX := wnd.x
	newY := wnd.y
	newH := wnd.h
	newW := wnd.w

	// logic
	{
		//Прямоугольник справа
		vResizeRect := utils.NewRect(wnd.x+wnd.w-scrollChange, wnd.y, scrollChange+5, wnd.h)
		hResizeRect := utils.NewRect(wnd.x, wnd.y+wnd.h-scrollChange, wnd.w, scrollChange+5)
		if utils.PointInRect(c.io.MousePos, hResizeRect) && c.ActiveWindow == wnd {
			c.io.SetCursor(VResizeCursor)
			wnd.wantResizeH = true
		} else if utils.PointInRect(c.io.MousePos, vResizeRect) && c.ActiveWindow == wnd {
			c.io.SetCursor(HResizeCursor)
			wnd.wantResizeV = true
		} else {
			c.io.SetCursor(ArrowCursor)
		}
		c.dragBehavior(vResizeRect, &wnd.capturedV)
		c.dragBehavior(hResizeRect, &wnd.capturedH)
		// Изменение размеров окна
		if wnd.wantResizeH && c.ActiveWindow == wnd && wnd.capturedH {
			n := newH
			n += c.io.MouseDelta.Y
			if n > wnd.minH {
				newH = n
				if wnd.mainWidgetSpace.scrlY != 0 {
					wnd.mainWidgetSpace.scrlY -= c.io.MouseDelta.Y
				}
			}
		} else if wnd.wantResizeV && c.ActiveWindow == wnd && wnd.capturedV {
			n := newW
			n += c.io.MouseDelta.X
			if n > wnd.minW {
				newW = n
			}
		} else {
			wnd.wantResizeH = false
			wnd.wantResizeV = false
		}

		c.dragBehavior(wnd.outerRect, &wnd.capturedWin)
		// Изменение положения окна
		if c.ActiveWindow == wnd && wnd.capturedWin && !wnd.wantResizeV &&
			!wnd.wantResizeH && !wnd.capturedInsideWin && !wnd.capturedTextSelection {
			newX += c.io.MouseDelta.X
			newY += c.io.MouseDelta.Y
		}
	}

	wnd.x = newX
	wnd.y = newY
	wnd.h = newH
	wnd.w = newW

	wnd.outerRect = utils.NewRect(wnd.x, wnd.y, wnd.w-wnd.mainWidgetSpace.verticalScrollbar.w, wnd.h)

	wnd.mainWidgetSpace.X = newX
	wnd.mainWidgetSpace.Y = newY + wnd.toolbar.h
	wnd.mainWidgetSpace.W = newW
	wnd.mainWidgetSpace.H = newH - wnd.toolbar.h

	wnd.mainWidgetSpace.cursorX = wnd.mainWidgetSpace.X + uiCtx.CurrentStyle.LeftMargin
	wnd.mainWidgetSpace.cursorY = wnd.mainWidgetSpace.Y + uiCtx.CurrentStyle.TopMargin

	wnd.mainWidgetSpace.ClipRect = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y,
		wnd.mainWidgetSpace.W - wnd.mainWidgetSpace.verticalScrollbar.w, wnd.mainWidgetSpace.H}

	hoveredWs := c.hoverBehavior(wnd, utils.NewRectS(wnd.mainWidgetSpace.ClipRect))
	actWind := c.ActiveWindow == wnd && c.HoveredWindow == wnd
	wdig := wnd.mainWidgetSpace.widgetSpaceLogic(hoveredWs, actWind, c.io.ScrollY != 0, wnd.buffer, func() draw.ClipRectCompose {
		cl := [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}
		return draw.NewClip(draw.EmptyClip, cl)
	})
	wnd.addDelayedWidget(wdig)

	DrawRoundedRect(wnd.buffer, windowName+"-main", wnd.x, wnd.y, wnd.w, wnd.h, mainWindowClr2, draw.AllRounded, 10)
	//DrawRect(wnd.buffer, windowName+"-main", wnd.x, wnd.y, wnd.w, wnd.h, mainWindowClr2)
	toolbar := wnd.toolbar
	DrawRoundedRect(wnd.buffer, windowName+"-toolbar", wnd.x, wnd.y, wnd.w, toolbar.h, toolbar.clr, draw.TopRect, 10)
	wnd.buffer.SeparateBuffer(0, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))

	if len(c.SelectedTexts) != 0 {
	}

	c.windowStack.Push(wnd)

	if wnd.toolbar.h != 0 {
		//b, hovered, clicked := c.ButtonEX(windowName+"-btn", wnd, wnd.x+wnd.w-25, wnd.y, 25, 25)
		//if hovered {
		//	b.SetColor(softGreen)
		//} else {
		//	b.SetColor(red)
		//}
		//if clicked {
		//	*opened = !*opened
		//}

		//txt, _, _, _, out := c.TextEX(wnd.mainWidgetSpace, windowName+"-txt", windowName, 0, DefaultTextFlag)
		//if out { // Prevent cursor changing because we draw it outside main widget space
		//	wnd.mainWidgetSpace.addCursor(-txt.Width(), -txt.Height())
		//	wnd.mainWidgetSpace.AddVirtualWH(-txt.Width(), -txt.Height())
		//}
		//txt.UpdatePosition([4]float32{wnd.x, wnd.y, txt.Width(), txt.Height()})
		//c.DrawText(wnd.x, wnd.y, txt, c.font.TextureId, wnd.buffer, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
	}

	// Draw selected text regions. We do it here because we don't want to draw it in front of text.
	// Maybe in future I will change text selection algorithm and rework this.
	for i, region := range wnd.textRegions {
		coords := region.Min
		bound := region.Max
		DrawRect(wnd.buffer, fmt.Sprint(i)+"tregfg", coords.X, coords.Y, bound.X, bound.Y, softGreen)
	}
	wnd.buffer.SeparateBuffer(0, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
}

var step float32 = 100

func (ws *WidgetSpace) widgetSpaceLogic(hovering, actWind, isScrollable bool, buffer *draw.CmdBuffer, scrollClip func() draw.ClipRectCompose) (delayedWidgets func()) {
	c := uiCtx

	if hovering {
		c.ActiveWidgetSpaceId = ws.id
		if ws.flags&Scrollable != 0 {
			c.WantScrollFocusWidgetSpaceId = ws.id
		}
	}
	// Scrollbar behavior
	if ws.flags&Scrollable != 0 {
		ws.vertScrollBar()
		//if c.ActiveWindow == wnd && c.HoveredWindow == wnd && c.io.ScrollY != 0 && c.WantScrollFocusWidgetSpaceLastId == ws.id && c.FocusedWidgetSpace == nil {
		if actWind && isScrollable && c.WantScrollFocusWidgetSpaceLastId == ws.id && c.FocusedWidgetSpace == nil {
			ws.handleMouseScroll(float32(c.io.ScrollY))
		}
		if ws.flags&ShowScrollbar != 0 && ws.isVertScrollShown {
			delayedWidgets = func() {
				scrlClip := scrollClip()
				scrl := ws.verticalScrollbar
				DrawRoundedRect(buffer, ws.id+"-vertScroll", scrl.x, scrl.y, scrl.w, scrl.h, scrl.clr, draw.AllRounded, 5)
				DrawRoundedRect(buffer, ws.id+"-vertScroll-btn", scrl.bX, scrl.bY, scrl.bW, scrl.bH, [4]float32{255, 0, 0, 1}, draw.AllRounded, 5)
				buffer.SeparateBuffer(0, scrlClip)
			}
		}
	}
	return
}

var (
	whiteColor     = [4]float32{255, 255, 255, 1}
	softGreen      = [4]float32{231, 240, 162, 0.8}
	black          = [4]float32{0, 0, 0, 1}
	red            = [4]float32{255, 0, 0, 1}
	transparent    = [4]float32{0, 0, 0, 0}
	mainWindowClr  = [4]float32{231, 158, 162, 0.8}
	mainWindowClr2 = [4]float32{29, 29, 29, 1}
)

func (c *UiContext) getWidget(id string, f func() widgets.Widget) widgets.Widget {
	var widg widgets.Widget
	widg, ok := c.GetWidget(id)
	if !ok {
		widg = f()
		c.AddWidget(id, widg)
	}
	return widg
}

//func ResolvePadding(x, y float32) (padLeft, padTop float32, box [4]float32) {
//	c := ctx()
//	padLeft = c.CurrentStyle.Padding.Left
//	padTop = c.CurrentStyle.Padding.Top
//
//}

func TextButton(id string, msg string) bool {
	c := ctx()
	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return false
	}
	txt, _, _, isRow, _ := c.TextEX(ws, id+"-text", msg, 0, DefaultTextFlag)
	padding := c.CurrentStyle.Padding
	var clicked, hovered bool
	btn, x, y, isRow, out := c.ButtonEX(ws, id+"-btn", txt.Width()+padding.WidthSum(), txt.Height()+padding.HeightSum())
	if out {
		return false
	}
	hovered = IsHovered(wnd, ws, btn)
	clicked = hovered && c.io.MouseClicked[0]
	if hovered {
		btn.SetColor(c.CurrentStyle.BtnHoveredColor)
		if clicked {
			btn.ChangeActive()
		}
	} else if btn.IsActive {
		btn.SetColor(c.CurrentStyle.BtnActiveColor)
	} else {
		btn.SetColor(c.CurrentStyle.BtnColor)
	}
	clip, buffer := ws.UpdateWidgetPosition(x, y, isRow, wnd, btn)
	c.DrawButton(x, y, btn, buffer, clip)

	txt.UpdatePosition([4]float32{x + padding.Left, y + padding.Top, txt.Width(), txt.Height()})

	c.DrawText(x+padding.Left, y+padding.Top, txt, c.font.TextureId, buffer, clip)
	return clicked
}

func Slider(id string, i *float32, min, max float32) {
	c := ctx()
	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return
	}
	slider := c.getWidget(id, func() widgets.Widget {
		return widgets.NewSlider(id, 0, 0, 200, 50, min, max, c.CurrentStyle)
	}).(*widgets.Slider)
	x, y, isRow, out := c.WidgetEX(ws, slider.Width(), slider.Height())
	if out {
		return
	}

	// logic
	{
		slider.HandleMouseDrag(c.io.MouseDelta.X, i, c.dragBehaviorInWindow)
		slider.CalculateNumber(i)
		// In the first launch if number more or less than borders values, we have to make it equal one of them
		if *i > max {
			*i = max
		} else if *i < min {
			*i = min
		}
	}

	clip, buffer := ws.UpdateWidgetPosition(x, y, isRow, wnd, slider)

	sl := slider.MainSliderPos()
	DrawRect(buffer, id+"slider", sl[0], sl[1], sl[2], sl[3], softGreen)
	btn := slider.BtnSliderPos()
	DrawRect(buffer, id+"sliderBtn", btn[0], btn[1], btn[2], btn[3], black)
	wnd.buffer.SeparateBuffer(0, clip)
}
func DrawRect(buff *draw.CmdBuffer, id string, x, y, w, h float32, clr [4]float32) *widgets.Rectangle {
	return DrawRectEX(buff, id, x, y, w, h, clr, draw.StraightCorners, 0)
}
func DrawRoundedRect(buff *draw.CmdBuffer, id string, x, y, w, h float32, clr [4]float32, shape draw.RoundedRectShape, radius int) *widgets.Rectangle {
	return DrawRectEX(buff, id, x, y, w, h, clr, shape, radius)
}

func DrawImage(buffer *draw.CmdBuffer, id string, x, y, w, h float32, texId uint32, texCoords, clr [4]float32, clip draw.ClipRectCompose) {
	c := ctx()
	img := c.getWidget(id, func() widgets.Widget {
		img := widgets.NewImage2(id, x, y, w, h, texId, texCoords, clr)
		img.ToggleUpdate()
		return img
	}).(*widgets.Image)
	img.UpdatePosition([4]float32{x, y, w, h})

	if img.Updated || buffer.CheckIndicesChange(img) {
		img.Vertices, img.Indices, img.VertCount, img.LastBufferIndex = buffer.CreateTexturedRect(x, y,
			img.Width(), img.Height(), img.TexId, img.TexCoords, img.Color())
		buffer.SendToBuffer(img.Vertices, img.Indices, img.VertCount, 0)
		img.Updated = false
	} else {
		buffer.SendToBuffer(img.Vertices, img.Indices, img.VertCount, img.LastBufferIndex)
	}
	buffer.SeparateBuffer(img.TexId, clip)
}

func DrawRectEX(buff *draw.CmdBuffer, id string, x, y, w, h float32, clr [4]float32, shape draw.RoundedRectShape, radius int) *widgets.Rectangle {
	c := ctx()
	r := c.getWidget(id, func() widgets.Widget {
		r := widgets.NewRectangle(id, x, y, w, h, clr)
		r.ToggleUpdate()
		return r
	}).(*widgets.Rectangle)

	r.UpdatePosition([4]float32{x, y, w, h})

	if r.Updated || buff.CheckIndicesChange(r) {
		switch shape {
		case draw.StraightCorners:
			r.Vertices, r.Indices, r.VertCount, r.LastBufferIndex = buff.CreateRect(x, y, r.Width(), r.Height(), r.BackgroundColor)
		default:
			r.Vertices, r.Indices, r.VertCount, r.LastBufferIndex = buff.CreateRoundedRect(x, y, r.Width(), r.Height(), radius, shape, r.BackgroundColor)
		}
		buff.SendToBuffer(r.Vertices, r.Indices, r.VertCount, 0)
		r.Updated = false
	} else {
		buff.SendToBuffer(r.Vertices, r.Indices, r.VertCount, r.LastBufferIndex)
	}
	return r
}

func RuneIndex(s string, c rune, fromIndex int) int {
	r := []rune(s)
	ind := 0
	for i := fromIndex; i <= len(r)-1; i++ {
		if r[i] != c {
			ind++
		} else {
			break
		}
	}
	return ind + fromIndex
}
func LastRuneIndex(s string, c rune) int {
	r := []rune(s)
	ind := 0
	for i := len(r) - 1; i >= 0; i-- {
		if r[i] != c {
			ind++
		} else {
			ind++
			break
		}
	}
	return len(r) - ind
}

// wrap Taken from https://commons.apache.org/proper/commons-lang/apidocs/org/apache/commons/lang3/text/WordUtils.html
// FIXME: found a bug when width of the first word less than wrapLength. The first char replaces with '\n'.
func wrap(msg string, wrapLength int) string {
	str := []rune(msg)
	inputLineLength := len(str)
	offset := 0
	sb := strings.Builder{}
	sb.Grow(len(str))
	for inputLineLength-offset > wrapLength {
		if str[offset] == ' ' {
			offset++
			continue
		}
		spaceToWrapAt := LastRuneIndex(string(str[:wrapLength+offset+1]), ' ')
		if spaceToWrapAt >= offset {
			sb.Write([]byte(string(str[offset:spaceToWrapAt])))
			sb.Write([]byte(string(' '))) // Add space in case of words selection. Without it, words can stick together.
			sb.Write([]byte(string('\n')))
			offset = spaceToWrapAt + 1
		} else {
			spaceToWrapAt = RuneIndex(string(str), ' ', wrapLength+offset)
			if spaceToWrapAt >= 0 {
				sb.Write([]byte(string(str[offset:spaceToWrapAt])))
				sb.Write([]byte(string('\n')))

				offset = spaceToWrapAt + 1
			} else {
				sb.Write([]byte(string(str[offset:])))
				offset = inputLineLength
			}
		}
	}
	if offset > len(str) {
	} else {
		sb.Write([]byte(string(str[offset:])))
	}
	return sb.String()
}

// FitTextToWidth is used to split long string to lines with defined width in px.
// If the length of the word more than border, its chars will be moved to the next line.
// |the quick brown fox jum|
// |ps over the lazy dog   |
func (c *UiContext) FitTextToWidth(x, w float32, msg string) string {
	r := []rune(msg)
	sb := strings.Builder{}
	sb.Grow(len(r))
	var dw float32 = 0
	for _, l := range r {
		char := c.font.GetCharacter(l)
		if x+dw+float32(char.Advance) < x+w {
			dw += float32(char.Advance)
			sb.Write([]byte(string(l)))
		} else {
			dw = 0
			sb.Write([]byte(string('\n')))
			// Do not let to create space at the new line
			if l != ' ' {
				sb.Write([]byte(string(l)))
			}
		}
	}
	return sb.String()
}

func insertIntoString2(s string, lineIndx, index int, val string, lines []fonts.TextLine) string {
	sb := strings.Builder{}
	sb.Grow(len(s))
	tmp := []rune(s)

	realIndx := 0
	{
		if lineIndx == 0 {
			realIndx = index
		} else {
			for i := 0; i < lineIndx; i++ {
				realIndx += len(lines[i].Text) + 1
			}
			realIndx += index
		}
	}
	index = realIndx

	for i, _ := range tmp[:realIndx] {
		sb.WriteString(string(tmp[i]))
	}
	sb.WriteString(val)
	for _, sr := range tmp[realIndx:] {
		sb.WriteString(string(sr))
	}
	return sb.String()
}

func removeFromString(s string, lineIndx, index int, lines []fonts.TextLine) string {
	tmp := []rune(s)
	if index < 0 {
		return s
	}
	realIndx := 0
	{
		if lineIndx == 0 {
			realIndx = index
		} else {
			for i := 0; i < lineIndx; i++ {
				realIndx += len(lines[i].Text) + 1
			}
			realIndx += index
		}
	}
	index = realIndx
	if index == 0 && lineIndx == 0 {
		return string(tmp[1:])
	} else if index == len(tmp) {
		fmt.Println("here")
	}
	tmp = append(tmp[:index], tmp[index+1:]...)
	return string(tmp)
}

func (c *UiContext) ClearTextSelection(wnd *Window) {
	c.SelectedTextStart = nil
	c.SelectedTextEnd = nil
	c.SelectedText = ""
	wnd.textRegions = []utils.Rect{}
}

func (c *UiContext) InputTextEX(ws *WidgetSpace, wnd *Window, id string, inputMsg string, inpmsg *string, key GuiKey, flag widgets.TextFlag) (txt *widgets.Text, x, y float32, isRow, outOfWs bool) {
	txt = c.getWidget(id, func() widgets.Widget {
		txt = widgets.NewTextNew(id, *inpmsg, ws.cursorX, ws.cursorY, c.font, c.CurrentStyle, flag)
		txt.Updated = true
		return txt
	}).(*widgets.Text)

	x, y, isRow, outOfWs = c.WidgetEX(ws, txt.Width(), txt.Height())

	if key != GuiKey_None && flag&widgets.Editable != 0 && c.FocusedTextInput == txt {
		if key == GuiKey_Backspace {
			txt.Editor.Backspace()
		} else if key == GuiKey_RightArrow {
			txt.Editor.MoveCharRight()
		} else if key == GuiKey_UpArrow {
			txt.Editor.MoveLineUp()
		} else if key == GuiKey_DownArrow {
			txt.Editor.MoveLineDown()
		} else if key == GuiKey_Delete {
			txt.Editor.Delete()
		} else if key == GuiKey_LeftArrow {
			txt.Editor.MoveCharLeft()
		} else if key == GuiKey_Enter {
			txt.Editor.InsertText("\n")
		} else if IsCommandKey(key) {

		} else if key != GuiKey_None {
			txt.Editor.InsertText(inputMsg)
		}
		//width, h, l, chars := c.font.CalculateTextBounds(*inpmsg, c.CurrentStyle.FontScale)
		//txt.Lines = l
		//txt.Chars = chars
		//txt.SetWH(width, h)
		//txt.Message = *inpmsg

		txt.ToggleUpdate()
		//ToggleAllWidgets()
	}
	return
}

func (c *UiContext) getTextInput() (string, GuiKey) {
	k := ""
	key := GuiKey_None
	if c.io.KeyPressedThisFrame && c.FocusedTextInput != nil {
		key = c.io.PressedKey
		k = c.io.keyToString(key)
	}
	return k, key
}

func TextInput(id string, w, h float32, message *string) {
	c := ctx()
	wnd, rootWs := c.IsWidgetSpaceAvailable()
	if rootWs == nil {
		fmt.Println("Can't find any widget spaces")
		return
	}
	var txt *widgets.Text
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws := c.subWidgetSpaceHelperWithBackground(wnd, wnd.buffer, id, x, y, w, h, 0, 0, softGreen, draw.StraightCorners, Scrollable|FitWidth, func() {
		currWs := wnd.currentWidgetSpace
		msg, key := c.getTextInput()
		txtTmp, x, y, isRow, out := c.InputTextEX(currWs, wnd, id, msg, message, key, widgets.Editable|Selectable)
		if out {
			return
		} else {
			wnd.VisibleTexts = append(wnd.VisibleTexts, txtTmp)
		}
		txt = txtTmp

		clip, buffer := currWs.UpdateWidgetPosition(x, y, isRow, wnd, txt)
		//if c.FocusedTextInput == txt {
		//	xc, yc, wc, hc := txt.CalculateCursorPos()
		//	DrawRect(buffer, id+"cursor", x+xc, y+yc, wc, hc, red)
		//}
		buffer.SeparateBuffer(0, clip)
		c.DrawText(x, y, txt, c.font.TextureId, buffer, clip)
	})
	// It is necessary to always monitor the state of the focused text.
	if c.FocusedTextInput == txt && c.io.MouseClicked[0] {
		if utils.PointOutsideRect(c.io.MouseClickedPos[0], utils.NewRectS(ws.ClipRect)) {
			c.FocusedTextInput = nil
			txt.ToggleUpdate()
		}
	}

	if c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect)) && c.io.MouseClicked[0] {
		//txt.ToggleUpdate()
		//txt.CursorInd = len(txt.Chars)
		c.FocusedTextInput = txt
		//pos := c.io.MouseClickedPos[0]
		//startFounded := false
		//for _, line := range txt.Lines {
		//	if pos.Y > line.StartY+y && pos.Y <= line.StartY+y+line.Height && !startFounded {
		//		startFounded = true
		//		txt.CursorInd = len(line.Text)
		//	}
		//}
	}
	rootWs.AddVirtualHeight(ws.H)
	rootWs.addCursor(ws.W, ws.H)
	//wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	//wnd.addCursor(ws.W, ws.H)
}

// MultiLineTextInput TODO: Add want input flag
func MultiLineTextInput(id string, message *string) {
	c := ctx()
	wnd, rootWs := c.IsWidgetSpaceAvailable()
	if rootWs == nil {
		fmt.Println("Can't find any widget spaces")
		return
	}
	var txt *widgets.Text
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws := c.subWidgetSpaceHelperWithBackground(wnd, wnd.buffer, id, x, y, wnd.mainWidgetSpace.W-(x-wnd.x)-wnd.mainWidgetSpace.verticalScrollbar.w, 200, 0, 0, softGreen, draw.StraightCorners, Scrollable|ShowScrollbar|FitWidth, func() {
		currWs := wnd.currentWidgetSpace
		msg, key := c.getTextInput()
		txtTmp, x, y, isRow, out := c.InputTextEX(currWs, wnd, id, msg, message, key, widgets.Editable|Selectable|widgets.MultiLine)
		if out {
			return
		} else {
			wnd.VisibleTexts = append(wnd.VisibleTexts, txtTmp)
		}
		txt = txtTmp

		clip, buffer := currWs.UpdateWidgetPosition(x, y, isRow, wnd, txt)
		if c.FocusedTextInput == txt {
			xc, yc, wc, hc := txt.CalculateCursorPos()
			DrawRect(buffer, id+"cursor", x+xc, y+yc, wc, hc, red)
		}
		//DrawRect(buffer, id+"-background", ws.X, ws.Y, ws.W, ws.H, whiteColor)
		buffer.SeparateBuffer(0, clip)
		c.DrawText(x, y, txt, c.font.TextureId, buffer, clip)
	})
	if c.FocusedTextInput == txt && c.io.MouseClicked[0] {
		if utils.PointOutsideRect(c.io.MouseClickedPos[0], utils.NewRectS(ws.ClipRect)) {
			c.FocusedTextInput = nil
			//txt.ToggleUpdate()
		}
	}

	if c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect)) && c.io.MouseClicked[0] {
		//txt.ToggleUpdate()
		txt.CursorInd = len(txt.Chars)
		c.FocusedTextInput = txt
		pos := c.io.MouseClickedPos[0]
		startFounded := false
		for _, line := range txt.Lines {
			if pos.Y > line.StartY+y && pos.Y <= line.StartY+y+line.Height && !startFounded {
				startFounded = true
				txt.CursorInd = len(line.Text)
			}
		}
	}
	rootWs.AddVirtualHeight(ws.H)
	rootWs.addCursor(ws.W, ws.H)
	//wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	//wnd.addCursor(ws.W, ws.H)
}

func TextFitted(id string, w float32, msg string) {
	c := ctx()
	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return
	}
	txt, x, y, isRow, out := c.TextEX(ws, id, msg, w, Selectable|widgets.SplitWords)
	if out {
		return
	}
	hovered := c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox())) && c.ActiveWidgetSpaceId == ws.id
	if hovered {
		txt.SetTextColor(softGreen)
	} else {
		txt.SetTextColor(c.CurrentStyle.TextColor)
	}
	clip, buffer := ws.UpdateWidgetPosition(x, y, isRow, wnd, txt)

	c.DrawText(x, y, txt, c.font.TextureId, buffer, clip)
}

func (wnd *Window) DefaultClip() draw.ClipRectCompose {
	return draw.NewClip(wnd.currentWidgetSpace.ClipRect, wnd.mainWidgetSpace.ClipRect)
}

func CalculateWidgetInfo(w, h float32, ws *WidgetSpace) (x, y float32, isRow, outWindow bool) {
	x, y, isRow = ws.getCursorPosition()
	if y+h > ws.Y && y <= ws.Y+ws.H {
		y += ws.resolveRowAlign(h)
		return
	} else {
		outWindow = true
		return
	}
}

// Widget
// Может рисоваться внутри окна (w,h)->clip.default или на произвольном месте (x,y,w,h)->clip.ignore
// Если внутри, нужна проверка на сущ. окно и возвращает => (x,y)
// ImageEX(id,x,y,w,h,tuid,tcoords)
// Image(id,w,h,tuid,tcoords) -> WidgetInfo() -> ImageEX(x,y)
// Global.Image(id,x,y,w,h,tuid,tcoords) -> ImageEX(x,y)

// ButtonEx(id,x,y,w,h,tuid,tcoords,clr)
// Button(id,w,h) -> WidgetInfo() -> ButtonEx(x,y,w,h,0,nil)
// Global.Button -> ButtonEx
// ImageButton() -> WidgetInfo() -> ButtonEX(x,y,w,h,1,[4])

// TextEX(id,x,y,w,h,msg,flag)

func (c *UiContext) WidgetEX(ws *WidgetSpace, w, h float32) (x, y float32, isRow, outOfWs bool) {
	x, y, isRow, outOfWs = CalculateWidgetInfo(w, h, ws)
	if outOfWs {
		ws.addCursor(w, h)
		if !isRow {
			ws.AddVirtualWH(w, h)
		}
	}
	return
}

func (c *UiContext) TextEX(ws *WidgetSpace, id string, msg string, newWidth float32, flag widgets.TextFlag) (txt *widgets.Text, x, y float32, isRow, outOfWs bool) {
	txt = c.getWidget(id, func() widgets.Widget {
		txt := widgets.NewTextNew(id, msg, ws.cursorX, ws.cursorY, c.font, c.CurrentStyle, flag)
		txt.Updated = true
		return txt
	}).(*widgets.Text)
	x, y, isRow, outOfWs = c.WidgetEX(ws, txt.Width(), txt.Height())

	if newWidth != 0 && txt.LastWidth != newWidth || txt.Message != msg {
		txt.ToggleUpdate()
		if txt.Flag&widgets.SplitChars != 0 {
			msg = c.FitTextToWidth(txt.BoundingBox()[0], newWidth, msg)
		} else if txt.Flag&widgets.SplitWords != 0 {
			numChars := int(math.Floor(float64(newWidth / c.font.XCharAdvance())))
			msg = wrap(msg, numChars)
		}
		width, height := txt.Editor.ReplaceBuffer(msg)
		txt.Message = msg
		txt.SetWH(width, height)
		txt.LastWidth = newWidth
	}
	return
}

func (c *UiContext) ImageEX(ws *WidgetSpace, id string, w, h float32, texId uint32, texCoords, clr [4]float32) (img *widgets.Image, x, y float32, isRow, outOfWs bool) {
	img = c.getWidget(id, func() widgets.Widget {
		img2 := widgets.NewImage2(id, ws.cursorX, ws.cursorY, w, h, texId, texCoords, clr)
		img2.Updated = true
		return img2
	}).(*widgets.Image)

	x, y, isRow, outOfWs = CalculateWidgetInfo(w, h, ws)
	if outOfWs {
		ws.addCursor(w, h)
		if !isRow {
			ws.AddVirtualWH(w, h)
		}
		outOfWs = true
		return
	}
	return
}

func ClipRect(wnd *Window) draw.ClipRectCompose {
	var clip draw.ClipRectCompose
	if wnd.currentWidgetSpace.flags&IgnoreClipping != 0 {
		clip = draw.NewClip(draw.EmptyClip, wnd.currentWidgetSpace.ClipRect)
	} else {
		clip = wnd.DefaultClip()
	}
	return clip
}
func (c *UiContext) DrawButton(x, y float32, btn *widgets.Button, buffer *draw.CmdBuffer, clip draw.ClipRectCompose) {
	if btn.Updated || buffer.CheckIndicesChange(btn) {
		btn.Vertices, btn.Indices, btn.VertCount, btn.LastBufferIndex = buffer.CreateRect(x, y,
			btn.Width(), btn.Height(), btn.Color())
		buffer.SendToBuffer(btn.Vertices, btn.Indices, btn.VertCount, 0)
		btn.Updated = false
	} else {
		buffer.SendToBuffer(btn.Vertices, btn.Indices, btn.VertCount, btn.LastBufferIndex)
	}
	buffer.SeparateBuffer(0, clip)
}
func (c *UiContext) DrawImage(x, y float32, img *widgets.Image, buffer *draw.CmdBuffer, clip draw.ClipRectCompose) {
	if img.Updated || buffer.CheckIndicesChange(img) {
		img.Vertices, img.Indices, img.VertCount, img.LastBufferIndex = buffer.CreateTexturedRect(x, y,
			img.Width(), img.Height(), img.TexId, img.TexCoords, img.Color())
		buffer.SendToBuffer(img.Vertices, img.Indices, img.VertCount, 0)
		img.Updated = false
	} else {
		buffer.SendToBuffer(img.Vertices, img.Indices, img.VertCount, img.LastBufferIndex)
	}
	buffer.SeparateBuffer(img.TexId, clip)
}

func GlobalWidgetSpace(id string, x, y, w, h float32, flag WidgetSpaceFlag, widgFunc func()) {
	c := ctx()
	c.subWidgetSpaceHelper(nil, c.globalBuffer, id+"global-ws", x, y, w, h, flag, func() {
		ws := c.getWidgetSpace(id+"global-ws", w, h, flag)
		c.CurrentGlobalWidgetSpace = ws
		wnd := c.windowStack.Peek()
		// If global ws was called inside ui.BeginWindow/End we should handle it. For example Tooltip widget
		// can be called inside window, but it contents should be rendered after window contents.
		if wnd == nil {
			c.CurrentGlobalWidgetSpace = ws
			widgFunc()
			c.CurrentGlobalWidgetSpace = nil
		} else {
			prevWs := wnd.currentWidgetSpace
			wnd.currentWidgetSpace = ws
			widgFunc()
			wnd.currentWidgetSpace = prevWs
		}
		c.CurrentGlobalWidgetSpace = nil
	})
}

func GlobalImage(id string, x, y, w, h float32, texId uint32, texCoords [4]float32) bool {
	c := ctx()
	var ws *WidgetSpace
	var clicked bool
	ws = c.subWidgetSpaceHelper(nil, c.globalBuffer, "wsp-"+id, x, y, w, h, IgnoreClipping, func() {
		ws = c.getWidgetSpace("wsp-"+id, w, h, IgnoreClipping)
		img, x, y, isRow, out := c.ImageEX(ws, id, w, h, texId, texCoords, whiteColor)
		if out {
			return
		}
		hovered := utils.PointInRect(c.io.MousePos, utils.NewRectS(img.BoundingBox()))
		clicked = hovered && c.io.MouseClicked[0]
		if hovered {
			img.SetColor(red)
		} else {
			img.SetColor(whiteColor)
		}
		clip, buff := ws.UpdateWidgetPosition(x, y, isRow, nil, img)

		c.DrawImage(x, y, img, buff, clip)
	})

	return clicked
}

func (c *UiContext) IsWidgetSpaceAvailable() (wnd *Window, ws *WidgetSpace) {
	wnd = c.getPeekWindow()
	if wnd != nil {
		ws = wnd.currentWidgetSpace
	} else if c.CurrentGlobalWidgetSpace != nil {
		ws = c.CurrentGlobalWidgetSpace
	}
	return
}

func IsHovered(wnd *Window, ws *WidgetSpace, w widgets.Widget) bool {
	c := ctx()
	if wnd == nil {
		return utils.PointInRect(c.io.MousePos, utils.NewRectS(w.BoundingBox())) && c.ActiveWidgetSpaceId == ws.id
	} else {
		hovered := c.hoverBehavior(wnd, utils.NewRectS(w.BoundingBox()))
		if hovered {
			//fmt.Println(w.WidgetId(), " ", w.BoundingBox())
		}
		//box := w.BoundingBox()
		//DrawRect(c.globalBuffer, w.WidgetId()+"debugdsd", box[0]-15, box[1]-15, box[2]+10, box[3]+10, whiteColor)
		return hovered && c.ActiveWidgetSpaceId == ws.id
	}
}

func Image(id string, w, h float32, texId uint32, texCoords [4]float32) bool {
	c := ctx()

	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return false
	}

	img, x, y, isRow, out := c.ImageEX(ws, id, w, h, texId, texCoords, whiteColor)
	if out {
		return false
	}

	hovered := IsHovered(wnd, ws, img)
	clicked := hovered && c.io.MouseClicked[0]
	if hovered {
		img.SetColor(red)
		c.setActiveWidget(img.WidgetId())
	} else {
		img.SetColor(whiteColor)
	}
	clip, buffer := ws.UpdateWidgetPosition(x, y, isRow, wnd, img)

	c.DrawImage(x, y, img, buffer, clip)
	return clicked
}

func Text(id string, msg string, flag widgets.TextFlag) {
	c := ctx()
	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return
	}

	txt, x, y, isRow, out := c.TextEX(ws, id, msg, 0, flag)

	if out {
		return
	} else if flag&Selectable != 0 {
		wnd.VisibleTexts = append(wnd.VisibleTexts, txt)
	}
	hovered := IsHovered(wnd, ws, txt)
	if hovered {
		txt.SetTextColor(softGreen)
	} else {
		txt.SetTextColor(c.CurrentStyle.TextColor)
	}
	clip, buffer := ws.UpdateWidgetPosition(x, y, isRow, wnd, txt)

	c.DrawText(x, y, txt, c.font.TextureId, buffer, clip)
	//wnd.debugDrawS(txt.BoundingBox())
	//wnd.buffer.SeparateBuffer(0, clip)
}

func (c *UiContext) DrawText(x, y float32, txt *widgets.Text, texid uint32, buffer *draw.CmdBuffer, clip draw.ClipRectCompose) {
	if txt.Updated || buffer.CheckIndicesChange(txt) {
		txt.Vertices, txt.Indices, txt.VertCount, txt.LastBufferIndex = buffer.CreateText(x, y, txt, txt.Scale, *c.font)
		buffer.SendToBuffer(txt.Vertices, txt.Indices, txt.VertCount, 0)
		txt.Updated = false
	} else {
		buffer.SendToBuffer(txt.Vertices, txt.Indices, txt.VertCount, txt.LastBufferIndex)
	}
	buffer.SeparateBuffer(texid, clip)
}
func (ws *WidgetSpace) UpdateWidgetPosition(xPos, yPos float32, isRow bool, wnd *Window, w widgets.Widget) (clip draw.ClipRectCompose, buffer *draw.CmdBuffer) {
	c := ctx()
	w.UpdatePosition([4]float32{xPos, yPos, w.Width(), w.Height()})
	ws.addCursor(w.Width(), w.Height())
	if !isRow {
		ws.AddVirtualWH(w.Width(), w.Height())
	}
	if wnd == nil || c.CurrentGlobalWidgetSpace != nil {
		clip = ws.Clip()
		buffer = c.globalBuffer
	} else {
		clip = ClipRect(wnd)
		buffer = wnd.buffer
	}
	return
}

// TODO: measure performance
func (wnd *Window) endWidget(xPos, yPos float32, isRow bool, w widgets.Widget) draw.ClipRectCompose {
	w.UpdatePosition([4]float32{xPos, yPos, w.Width(), w.Height()})
	wnd.addCursor(w.Width(), w.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualWH(w.Width(), w.Height())
	}

	//wnd.debugDraw(xPos, yPos, w.Width(), w.Height())

	var clip draw.ClipRectCompose
	if wnd.currentWidgetSpace.flags&IgnoreClipping != 0 {
		clip = draw.NewClip(draw.EmptyClip, wnd.currentWidgetSpace.ClipRect)
	} else {
		clip = wnd.DefaultClip()
	}
	return clip
}

func (w *Window) addDelayedWidget(f func()) {
	if f != nil {
		w.delayedWidgets = append(w.delayedWidgets, f)
	}
}

func VSpace(id string) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var s *widgets.VSpace
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	s = c.getWidget(id, func() widgets.Widget {
		s := widgets.NewVertSpace(id, [4]float32{x, y, 100, 20})
		return s
	}).(*widgets.VSpace)

	wnd.endWidget(x, y, isRow, s)
}

func (c *UiContext) hoverBehavior(wnd *Window, rect utils.Rect) bool {
	inRect := utils.PointInRect(c.io.MousePos, utils.NewRect(rect.Min.X, rect.Min.Y, rect.Width(), rect.Height()))
	if wnd == nil {
		// FIXME: incorrect calculation
		focusedWidgSpace := false
		// Accept widget actions only from focused widget space
		if c.FocusedWidgetSpace != nil {
			if c.CurrentGlobalWidgetSpace.id != c.ActiveWidgetSpaceId {
				focusedWidgSpace = true
			}
		}
		return inRect && !focusedWidgSpace
	}
	inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y, wnd.w, wnd.h)

	focusedWidgSpace := false
	// Accept widget actions only from focused widget space
	if c.FocusedWidgetSpace != nil {
		if wnd.currentWidgetSpace != c.FocusedWidgetSpace {
			focusedWidgSpace = true
		}
	}
	return c.ActiveWindow == wnd && inRect && inWindow && !focusedWidgSpace
}

func TreeNode(id string, msg string, widgFunc func()) bool {
	c := ctx()

	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return false
	}
	PushStyleVar1f(FontScale, 1)
	txt, _, _, isRow, out := c.TextEX(ws, id+"-header", msg, 0, DefaultTextFlag)
	PopStyleVar()
	if out {
		// Prevent cursor changing, because by default if widget is outside widget space,
		// cursor automatically increments. But in this case TreeNode is a composite
		// widget, so we don't want to lose correct positions inside TreeNode.
		// TODO(@Dmitry-dms): Illustrate this problem.
		ws.addCursor(-txt.Width(), -txt.Height())
		ws.AddVirtualWH(-txt.Width(), -txt.Height())
	}

	var clicked, hovered bool
	btn, x, y, isRow, out2 := c.ButtonEX(ws, id+"-btn", wnd.w, txt.Height())
	if out2 && !btn.IsActive {
		return false
	} else if out2 && btn.IsActive {
		ws.addCursor(-btn.Width(), -btn.Height())
		ws.AddVirtualWH(-btn.Width(), -btn.Height())
	}
	{
		hovered = IsHovered(wnd, ws, btn)
		clicked = hovered && c.io.MouseClicked[0]
		if hovered {
			btn.SetColor(c.CurrentStyle.BtnHoveredColor)
			if clicked {
				btn.ChangeActive()
			}
		} else if btn.IsActive {
			btn.SetColor(c.CurrentStyle.BtnActiveColor)
		} else {
			btn.SetColor(c.CurrentStyle.BtnColor)
		}
		btn.SetWidth(wnd.w)
	}
	clip, buffer := ws.UpdateWidgetPosition(x, y, isRow, wnd, btn)
	c.DrawButton(x, y, btn, buffer, clip)

	txt.UpdatePosition([4]float32{x, y, txt.Width(), txt.Height()})

	c.DrawText(x, y, txt, c.font.TextureId, buffer, clip)
	if btn.IsActive {
		x += 50
		childWs := c.subWidgetSpaceHelper(wnd, wnd.buffer, id, x, y+btn.Height(), 0, 0, NotScrollable|Resizable, widgFunc)
		ws.addCursor(childWs.W, childWs.H)
		ws.AddVirtualWH(childWs.W, childWs.H)
	}

	return clicked
}

func (c *UiContext) textButton(id string, wnd *Window, msg string, x, y float32, align widgets.TextAlign) (tBtn *widgets.TextButton, hovered, clicked bool) {
	tBtn = c.getWidget(id, func() widgets.Widget {
		w, h, _, p := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewTextButton(id, x, y, w, h, msg, p, align, widgets.AllPadding, c.CurrentStyle)
	}).(*widgets.TextButton)
	if msg != tBtn.Message {
		tBtn.Message = msg
		w, h, _, p := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		tBtn.Text.Chars = p
		tBtn.SetWH(w, h)
	}
	hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.BoundingBox()))
	if hovered {
		c.setActiveWidget(tBtn.Id)
	}
	clicked = c.io.MouseClicked[0] && hovered
	if clicked {
		tBtn.ChangeActive()
	}
	return
}

func (c *UiContext) ButtonEX(ws *WidgetSpace, id string, w, h float32) (btn *widgets.Button, x, y float32, isRow, outOfWs bool) {
	btn = c.getWidget(id, func() widgets.Widget {
		return widgets.NewButton(id, ws.cursorX, ws.cursorY, w, h, c.CurrentStyle.BtnColor)
	}).(*widgets.Button)
	x, y, isRow, outOfWs = CalculateWidgetInfo(w, h, ws)
	if outOfWs {
		ws.addCursor(w, h)
		if !isRow {
			ws.AddVirtualWH(w, h)
		}
		outOfWs = true
		return
	}
	return
}

func Button(id string) bool {
	c := ctx()

	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return false
	}

	var clicked, hovered bool
	btn, x, y, isRow, out := c.ButtonEX(ws, id, 100, 100)
	if out {
		return false
	}
	hovered = IsHovered(wnd, ws, btn)
	clicked = hovered && c.io.MouseClicked[0]
	if hovered {
		btn.SetColor(c.CurrentStyle.BtnHoveredColor)
	} else if btn.IsActive {
		btn.SetColor(c.CurrentStyle.BtnActiveColor)
	} else {
		btn.SetColor(c.CurrentStyle.BtnColor)
	}
	clip, buffer := ws.UpdateWidgetPosition(x, y, isRow, wnd, btn)

	c.DrawButton(x, y, btn, buffer, clip)
	return clicked
}

func (wnd *Window) addCursor(width, height float32) {
	row, ok := wnd.currentWidgetSpace.getCurrentRow()
	if !ok {
		wnd.currentWidgetSpace.cursorY += height
	} else {
		if row.RequireColumn {
			row.CursorY += height
			row.UpdateColWidth(width)
			row.AddColHeight(height)
		} else {
			row.CursorX += width
			row.W += width
			row.UpdateHeight(height)
		}
	}
}

func (c *UiContext) setActiveWidget(id string) {
	c.ActiveWidget = id
}

func Selection(id string, index *int, data []string, texId uint32, texCoords [4]float32) {
	c := ctx()
	wnd, ws := c.IsWidgetSpaceAvailable()
	if ws == nil {
		fmt.Println("Can't find any widget spaces")
		return
	}
	var s *widgets.Selection
	// Need to use WS because text may not fit into ButtonEX, so it should be clipped
	SubWidgetSpace(id+"---", 0, 0, Resizable|NotScrollable, func() {
		Row(id+"row--", widgets.VerticalAlign, func() {
			originX, originY, isRow := wnd.currentWidgetSpace.getCursorPosition()
			s = c.getWidget(id, func() widgets.Widget {
				return widgets.NewSelection(id, originX, originY, 300, 40)
			}).(*widgets.Selection)

			//x2, y2, _ := wnd.currentWidgetSpace.getCursorPosition()
			img, _, imgY, _, out := c.ImageEX(ws, id+"sel-arrow", s.Height(), s.Height(), texId, texCoords, whiteColor)
			if out {
				return
			}
			//clip := wnd.endWidget(originX+s.Width()-imgX, imgY, isRow, img)
			img.UpdatePosition([4]float32{originX + s.Width() - img.Width(), imgY, img.Width(), img.Height()})
			hovered := IsHovered(wnd, ws, img)
			clicked := hovered && c.io.MouseClicked[0]
			DrawRect(wnd.buffer, id+"sel-rect", originX, originY, s.Width(), s.Height(), whiteColor)
			wnd.buffer.SeparateBuffer(0, wnd.DefaultClip())
			wnd.endWidget(originX, originY, isRow, s)
			//img, _, clicked := c.imageHelper(id+"arrow", x2, y2, s.Height(), s.Height(), func() *widgets.Image {
			//	return nil
			//})
			//
			c.DrawImage(originX+s.Width()-img.Width(), imgY, img, wnd.buffer, wnd.DefaultClip())

			txt, _, _, _, _ := c.TextEX(ws, data[*index]+"--"+id, data[*index], 0, DefaultTextFlag)
			txt.UpdatePosition([4]float32{originX, originY, txt.Width(), txt.Height()})
			c.DrawText(originX, originY, txt, c.font.TextureId, wnd.buffer, wnd.DefaultClip())
			//txt, _ := c.textHelper(data[*index]+"--"+id, x, y, 0, data[*index], widgets.Default)
			//wnd.buffer.CreateText(x+c.CurrentStyle.AllPadding, y+(s.Height()-txt.Height())/2, txt,
			//	*c.font, draw.NewClip(draw.EmptyClip, [4]float32{x, y, s.Width(), s.Height()}))
			//wnd.buffer.CreateTexturedRect(x2-s.Height(), y2, img.Width(), img.Height(), texId, texCoords, img.Color(), wnd.DefaultClip())
			if clicked {
				fmt.Println("clcd")
				s.Opened = true
				c.setActiveWidget(id)
			}
			if c.ActiveWidget != id {
				s.Opened = false
			}
		})
	})

	//ContextMenu(id, IgnoreClipping, func() {
	//	for i, datum := range data {
	//		x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	//		tbt, _, clicked := c.textButton(datum+"_btnT_"+id, wnd, datum, x, y, widgets.Left)
	//		if clicked {
	//			*index = i
	//			c.FocusedWidgetSpace = nil
	//		}
	//		tbt.SetWidth(s.Width())
	//		clip := wnd.endWidget(x, y, false, tbt)
	//		wnd.buffer.CreateButtonT(x, y, tbt, *c.font, clip)
	//	}
	//})
}

//func ContextMenu(ownerWidgetId string, flag WidgetSpaceFlag, widgFunc func()) {
//	c := ctx()
//	wnd := c.windowStack.Peek()
//	var bb [4]float32
//	widg, ok := c.GetWidget(ownerWidgetId)
//	if !ok {
//		return
//	}
//	bb = widg.BoundingBox()
//	id := ownerWidgetId + "-ws-context"
//	ws := c.getWidgetSpace(id, 0, 0, Resizable|FitWidth|flag)
//	if c.LastActiveWidget == widg.WidgetId() {
//		c.FocusedWidgetSpace = ws
//	}
//	if c.FocusedWidgetSpace == ws {
//		f := func() {
//			ws.ClipRect = [4]float32{ws.X, ws.Y, ws.W, ws.H}
//			clip := draw.NewClip(draw.EmptyClip, ws.ClipRect)
//			wnd.buffer.CreateRect(bb[0], bb[1]+widg.Height(), ws.W, ws.H, 0, draw.StraightCorners, 0, black, clip)
//			c.subWidgetSpaceHelper(wnd, wnd.buffer, id, bb[0], bb[1]+widg.Height(), widg.Width(), 0, Resizable|FitWidth|flag, widgFunc)
//		}
//		wnd.addDelayedWidget(f)
//	}
//}
func Tooltip(id string, widgFunc func()) {
	c := ctx()
	x, y := c.io.MousePos.X+10, c.io.MousePos.Y+5
	//wnd := c.windowStack.Peek()

	//ws := c.getWidgetSpace(id, 0, 0, Resizable|IgnoreClipping)
	GlobalWidgetSpace(id+"global-ws", x, y, 0, 0, Resizable|IgnoreClipping, widgFunc)

	//wnd.addDelayedWidget(func() {
	//GlobalWidgetSpace(id+"glWs", x, y, 0, 0, Resizable|IgnoreClipping, widgFunc)
	//c.subWidgetSpaceHelperWithBackground(wnd, wnd.buffer, id, x, y, 0, 0, 0, 0, black, draw.StraightCorners, Resizable|IgnoreClipping, widgFunc)
	//wnd.buffer.CreateRect(x, y, ws.W, ws.H, 0, draw.StraightCorners, 0, black, draw.NewClip(draw.EmptyClip, ws.ClipRect))
	//c.subWidgetSpaceHelper(wnd, wnd.buffer, id, x, y, 0, 0, Resizable|IgnoreClipping, widgFunc)
	//})
}

func Column(id string, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var hl *widgets.HybridLayout
	hl, ok := wnd.currentWidgetSpace.getCurrentRow()
	if !ok {
		return
	}
	hl.RequireColumn = true
	hl.CurrentColH, hl.CurrentColW = 0, 0

	widgFunc()

	hl.RequireColumn = false
	hl.CursorY = hl.InitY

	hl.W += hl.CurrentColW
	hl.CursorX += hl.CurrentColW
	hl.UpdateHeight(hl.CurrentColH)
}

func (wnd *Window) debugDrawS(x [4]float32) {
	wnd.buffer.CreateBorderBox(x[0], x[1], x[2], x[3], 2, red)
}
func (wnd *Window) debugDraw(x, y, w, h float32, clr [4]float32) {
	wnd.buffer.CreateBorderBox(x, y, w, h, 2, clr)
}

var wsWantRebuffer bool

func (c *UiContext) getWidgetSpace(id string, width, height float32, flags WidgetSpaceFlag) *WidgetSpace {
	ws, ok := c.widgSpaceCache.Get(id)
	if !ok {
		ws = newWidgetSpace(id, 0, 0, width, height, flags)
		c.widgSpaceCache.Add(id, ws)
		wsWantRebuffer = true
	}
	//wnd.widgSpaces = append(wnd.widgSpaces, ws)
	return ws
}

func (c *UiContext) subWidgetSpaceHelperEx(wnd *Window, buff *draw.CmdBuffer, id string, x, y, width, height float32, texId uint32, radius int, clr [4]float32, shape draw.RoundedRectShape, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	//wnd := c.windowStack.Peek()

	ws := c.getWidgetSpace(id, width, height, flags)
	var prevWS *WidgetSpace
	if wnd != nil {
		prevWS = wnd.currentWidgetSpace
		wnd.currentWidgetSpace = ws
	}

	ws.X = x
	ws.Y = y
	ws.cursorY = y
	ws.cursorX = x

	//wnd.debugDraw(ws.X, ws.Y, ws.W, ws.H, red)

	outOfWindow := false
	if wnd != nil {
		if y < wnd.mainWidgetSpace.Y {
			outOfWindow = true
			// vs-clip-1.png
			if ws.isVertScrollShown {
				ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W - ws.verticalScrollbar.w, ws.H - (wnd.mainWidgetSpace.Y - y)}
			} else {
				//ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W, ws.H - (wnd.mainWidgetSpace.Y - y)}
				ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W, ws.H}
			}
			//wnd.debugDrawS(ws.ClipRect)
		} else if y+ws.H > wnd.mainWidgetSpace.Y+wnd.mainWidgetSpace.H {
			ws.ClipRect = [4]float32{x, y, ws.W, ws.H - (wnd.mainWidgetSpace.Y - y)}
		} else {
			if ws.isVertScrollShown {

				ws.ClipRect = [4]float32{x, y, ws.W - ws.verticalScrollbar.w, ws.H}
			} else {
				ws.ClipRect = [4]float32{x, y, ws.W, ws.H}
			}
		}
	} else {
		if ws.isVertScrollShown {
			ws.ClipRect = [4]float32{x, y, ws.W - ws.verticalScrollbar.w, ws.H}
		} else {
			ws.ClipRect = [4]float32{x, y, ws.W, ws.H}
		}
	}

	if flags&FitWidth != 0 {
		if ws.isVertScrollShown {
			ws.ClipRect[2] = width - ws.verticalScrollbar.w
		} else {
			ws.ClipRect[2] = width
		}
		ws.W = width
	}

	var hoveredWs, actWind, scrollable bool
	if wnd != nil {
		//hoveredWs = c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect))
		hoveredWs = c.hoverBehavior(wnd, utils.NewRect(ws.X, ws.Y, ws.W, ws.H))
		scrollable = c.io.ScrollY != 0
		actWind = c.ActiveWindow == wnd && c.HoveredWindow == wnd
	} else {
		hoveredWs = utils.PointInRect(c.io.MousePos, utils.NewRectS(ws.ClipRect))
		scrollable = c.io.ScrollY != 0
		actWind = true
	}

	// TODO(@Dmitry-dms): Should each widget space have its own LastVertices,LastIndices etc... for drawing a scrollbar.
	// Now, global buffer handle it, so it redraws every frame if needed. If it would be a window buffer
	// we should handle its vertices carefully, because it breaks all widgets inside the window.
	ws.widgetSpaceLogic(hoveredWs, actWind, scrollable, buff, func() draw.ClipRectCompose {
		cl := [4]float32{ws.X, ws.Y, ws.W, ws.H}
		if outOfWindow {
			cl[1] = wnd.mainWidgetSpace.Y
		}
		return draw.NewClip(cl, wnd.mainWidgetSpace.ClipRect)
	})

	if flags&FillBackground != 0 {
		if texId == 0 {
			DrawRect(buff, id+"back", ws.X, ws.Y, ws.W, ws.H, clr)
			//buff.CreateRect(ws.X, ws.Y, ws.W, ws.H, radius, shape, 0, clr, wnd.DefaultClip())
		} else {
			//TODO: Add textured rect method
			//wnd.buffer.CreateTexturedRect()
		}
		buff.SeparateBuffer(0, ws.Clip())
	}

	widgFunc()
	ws.checkVerScroll()

	ws.lastVirtualHeight = ws.virtualHeight
	ws.virtualHeight = 0
	ws.lastVirtualWidth = ws.virtualWidth
	ws.virtualWidth = 0

	if ws.flags&Resizable != 0 {
		ws.H = ws.lastVirtualHeight
		ws.W = ws.lastVirtualWidth
	}

	//wnd.buffer.CreateRect(wnd.mainWidgetSpace.X, ws.H+y, wnd.w, 2,
	//	0, draw.StraightCorners, 0, c.CurrentStyle.WidgSpaceDividerColor, wnd.DefaultClip())
	//wnd.debugDraw(x, y, ws.W, ws.H)
	if wnd != nil {
		wnd.currentWidgetSpace = prevWS
	}
	return ws
}

func (c *UiContext) subWidgetSpaceHelperWithBackground(wnd *Window, buff *draw.CmdBuffer, id string, x, y, width, height float32, texId uint32, radius int, clr [4]float32, shape draw.RoundedRectShape, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	return c.subWidgetSpaceHelperEx(wnd, buff, id, x, y, width, height, texId, radius, clr, shape, flags|FillBackground, widgFunc)
}

func (c *UiContext) subWidgetSpaceHelper(wnd *Window, buff *draw.CmdBuffer, id string, x, y, width, height float32, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	return c.subWidgetSpaceHelperEx(wnd, buff, id, x, y, width, height, 0, 0, black, draw.StraightCorners, flags, widgFunc)
}

func SubWidgetSpace(id string, width, height float32, flags WidgetSpaceFlag, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var ws *WidgetSpace

	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws = c.subWidgetSpaceHelper(wnd, wnd.buffer, id, x, y, width, height, flags, widgFunc)

	wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	wnd.addCursor(ws.W, ws.H)
}

func TabItem(name string, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var tb *widgets.TabBar
	tb, ok := wnd.currentWidgetSpace.getCurrentTabBar()
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	if !ok {
		return
	}
	wspId := name + "-wsp-" + tb.WidgetId()
	_, index := tb.FindTabItem(name, wspId)

	var ws *WidgetSpace
	if index == tb.CurrentTab {
		ws = c.subWidgetSpaceHelper(wnd, wnd.buffer, wspId, x, y, 0, 0, Resizable|HideScrollbar, widgFunc)
	}
	if ws != nil {
		tb.SetHeight(ws.H)
		tb.SetWidth(ws.W)
	}
}
func TabBar(id string, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var tab *widgets.TabBar
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()

	tab = c.getWidget(id, func() widgets.Widget {
		return widgets.NewTabBar(id, x, y, 0, 0)
	}).(*widgets.TabBar)

	var rowHeight, rowWidth float32
	ws := c.subWidgetSpaceHelper(wnd, wnd.buffer, id, x, y, 0, 0, Resizable|NotScrollable, func() {
		//cr := wnd.currentWidgetSpace
		//wnd.buffer.CreateRect(cr.X, cr.Y, cr.W, cr.H, 10, draw.AllRounded, 0, softGreen, wnd.DefaultClip())
		//Row("rowds", func() {
		//	row, _ := wnd.currentWidgetSpace.getCurrentRow()
		//	//wnd.buffer.CreateRect(row.X, row.Y, row.Width(), row.Height(), 10, draw.AllRounded, 0, transparent, wnd.DefaultClip())
		//	for i, item := range tab.Bars {
		//		x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
		//		tbtn, hovered, clicked := c.textButton(fmt.Sprint(id, "-", i), wnd, item.Name, x, y, widgets.Center)
		//		if clicked {
		//			tab.CurrentTab = i
		//		}
		//		if hovered {
		//			tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnActiveColor)
		//			if clicked {
		//				tab.ChangeActive(item)
		//			}
		//		} else if item.Active {
		//			tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnActiveColor)
		//		} else {
		//			tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnColor)
		//		}
		//
		//		tbtn.Text.CurrentColor = whiteColor
		//		tbtn.SetHeight(tbtn.Height() - (tbtn.Height() - tbtn.Text.Height()) + c.CurrentStyle.AllPadding)
		//		clip := wnd.endWidget(x, y, false, tbtn)
		//		//wnd.buffer.CreateRect(x, y, tbtn.Width(), tbtn.Height(), 10, draw.TopRect, 0, tbtn.Color(), clip)
		//		//wnd.buffer.CreateRect(x, y, tbtn.Width(), tbtn.Height(), 10, draw.TopRect, 0, tbtn.Color(), clip)
		//		wnd.buffer.SeparateBuffer(0, clip)
		//		wnd.buffer.CreateText(tbtn.Text.BoundingBox()[0], tbtn.Text.BoundingBox()[1], tbtn.Text, *c.font, clip)
		//
		//		row, _ := wnd.currentWidgetSpace.getCurrentRow()
		//		if i != len(tab.Bars)-1 {
		//			row.CursorX += 10
		//			row.W += 10
		//		}
		//		if row.Height() > rowHeight {
		//			rowHeight = row.Height()
		//		}
		//		rowWidth = row.Width()
		//	}
		//})

		//wnd.buffer.CreateRect(wnd.x, y+rowHeight, wnd.w, 2, 0, draw.StraightCorners, 0, c.CurrentStyle.TabBtnActiveColor, draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect))
		wnd.buffer.SeparateBuffer(0, draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect))
		wnd.currentWidgetSpace.cursorY += 5

		tab.BarHeight = rowHeight
		tab.SetWidth(rowWidth)
		wnd.currentWidgetSpace.tabStack.Push(tab)
		widgFunc()
		wnd.currentWidgetSpace.tabStack.Pop()
	})
	ws.W = tab.Width()
	ws.H = tab.Height() + tab.BarHeight + 5

	//wnd.buffer.CreateRect(wnd.mainWidgetSpace.X, ws.H+y, wnd.w, 2,
	//	0, draw.StraightCorners, 0, c.CurrentStyle.WidgSpaceDividerColor, wnd.DefaultClip())

	wnd.addCursor(ws.W, ws.H)
	wnd.currentWidgetSpace.AddVirtualWH(ws.W, ws.H)
}

func Row(id string, align widgets.RowAlign, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var row *widgets.HybridLayout
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	// Need to return cursor back, because internal row cursor shouldn't know anything about outer
	y += wnd.currentWidgetSpace.scrlY

	row = c.getWidget(id, func() widgets.Widget {
		return widgets.NewHLayout(id, x, y, align, c.CurrentStyle)
	}).(*widgets.HybridLayout)
	row.UpdatePosition([4]float32{x, y, row.Width(), row.Height()})
	wnd.currentWidgetSpace.rowStack.Push(row)

	widgFunc()

	hl := wnd.currentWidgetSpace.rowStack.Pop()
	wnd.addCursor(0, hl.H)
	//wnd.endWidget(x, y, false, row)

	wnd.currentWidgetSpace.AddVirtualWH(hl.W, hl.H)
	hl.LastWidth = hl.W
	hl.LastHeight = hl.H
	hl.H = 0
	hl.W = 0
}

// solveTextSelection is responsible for finding selected texts on window
// TODO: improve this algorithm
// TODO: add upstring selection with several Text widgets. Today, works only downstring selection
func (c *UiContext) solveTextSelection(wnd *Window) {
	startFounded := false
	// Iterate through all visible texts on the screen
	for _, t := range wnd.VisibleTexts {
		if t.Flag&Selectable == 0 {
			continue
		}
		// If there are focused input text, don't let to select others text widgets
		// FIXME: Is this logic should be separated?
		if c.FocusedTextInput != nil {
			if t != c.FocusedTextInput {
				continue
			}
		}
		if c.hoverBehavior(wnd, utils.NewRectS(t.BoundingBox())) {
			c.dragBehavior(wnd.outerRect, &wnd.capturedTextSelection)
			// If we don't have a selected text yet, place a cursor to it
			if c.io.MouseClicked[0] {
				b := t.BoundingBox()
				x := c.io.dragStarted.X - b[0]
				y := c.io.dragStarted.Y - b[1]
				for _, line := range t.Editor.Lines {
					for i := line.Begin; i < line.End; i++ {
						char := t.Editor.CharsInfo[i]
						if x >= char.Xpos-5 && x <= char.Xpos+float32(char.Info.Advance) && !startFounded &&
							y <= line.Ypos+line.Height {

							t.StartInd = i
							startFounded = true
							t.Editor.Selection = true
							t.Editor.StartSelection(i)
							c.SelectedTextStart = t
							if t == c.FocusedTextInput {
								//t.CursorInd = ind
								//t.CursorLine = t.StartLine
							}
						}
					}
				}

				// If text start was founded, drag delta helps to find boundaries of selected texts
			} else if c.SelectedTextStart != nil && wnd.capturedTextSelection {
				b := t.BoundingBox()
				x := c.io.dragStarted.X + c.io.dragDelta.X - b[0]
				y := c.io.dragStarted.Y + c.io.dragDelta.Y - b[1]

				endfounded := false // Should use this because we need to find only the first hovered line
				for _, line := range t.Editor.Lines {
					for i := line.Begin; i < line.End; i++ {
						char := t.Editor.CharsInfo[i]
						if x >= char.Xpos-5 && x <= char.Xpos+float32(char.Info.Advance)*c.CurrentStyle.FontScale && !endfounded &&
							y <= line.Ypos+line.Height*c.CurrentStyle.FontScale {
							t.EndInd = i
							t.Editor.SelectionEnd = i
							c.SelectedTextEnd = t
							endfounded = true
							if i < t.Editor.SelectionBegin {
								t.Editor.SelectionPoint = t.Editor.SelectionBegin
							}
						}
					}
				}
			}
		}

		if c.SelectedTextStart != nil && c.SelectedTextEnd != nil {
			tmp := []*widgets.Text{} // It's a temporary slice which contains selected text widgets
			if c.SelectedTextStart == c.SelectedTextEnd {
				tmp = append(tmp, c.SelectedTextStart)
			} else { // If selected widgets have more than 1 widget, find the first, and add to tmp in order of creation until the last
				startFounded := false
				for _, text := range wnd.VisibleTexts {
					if text == c.SelectedTextStart {
						text.Editor.SelectionEnd = text.Editor.Buff.Len() - 1
						tmp = append(tmp, text)
						startFounded = true
						continue
					}
					if startFounded {
						tmp = append(tmp, text)
						if text == c.SelectedTextEnd {
							break
						}
					}
				}
			}
			rects := []utils.Rect{} // Boundaries for each selected line
			selectedString := ""
			for _, text := range tmp {
				msg, regions := text.Editor.GetTextSelection(text.BoundingBox()[0], text.BoundingBox()[1])
				selectedString += msg + " "
				rects = append(rects, regions...)
			}
			wnd.textRegions = rects
			c.SelectedText = selectedString
		}
	}
}

func EndWindow() {
	c := ctx()
	wnd := c.getPeekWindow()
	if wnd == nil {
		return
	}

	c.solveTextSelection(wnd)
	count := 0
	if len(wnd.VisibleTexts) != 0 {
		for _, text := range wnd.VisibleTexts {
			if utils.PointOutsideRect(c.io.MouseClickedPos[0], utils.NewRectS(text.BoundingBox())) {
				count++
			}
		}
		if count == len(wnd.VisibleTexts) {
			c.SelectedTextStart = nil
			c.SelectedTextEnd = nil
			c.SelectedText = ""
			wnd.textRegions = []utils.Rect{}
			//ToggleAllWidgets()
		}
	}

	for _, f := range wnd.delayedWidgets {
		f()
	}
	wnd.delayedWidgets = []func(){}

	wnd.VisibleTexts = wnd.VisibleTexts[:0]

	wnd = c.windowStack.Pop()
	wnd.mainWidgetSpace.checkVerScroll()

	//var clip = draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect)
	var clip = draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h})
	wnd.buffer.SeparateBuffer(0, clip) // Make sure that we didn't miss anything

	wnd.mainWidgetSpace.AddVirtualHeight(c.CurrentStyle.BotMargin)

	wnd.mainWidgetSpace.lastVirtualHeight = wnd.mainWidgetSpace.virtualHeight
	wnd.mainWidgetSpace.virtualHeight = 0
	wnd.mainWidgetSpace.lastVirtualWidth = wnd.mainWidgetSpace.virtualWidth
	wnd.mainWidgetSpace.virtualWidth = 0
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
