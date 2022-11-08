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
	defx, defy, defw, defh = 300, 100, 400, 500
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
		if c.ActiveWindow == wnd && wnd.capturedWin && !c.wantResizeV &&
			!c.wantResizeH && !wnd.capturedInsideWin && !wnd.capturedTextSelection {
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
	wnd.widgetSpaceLogic(wnd.mainWidgetSpace, func() draw.ClipRectCompose {
		cl := [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}
		return draw.NewClip(draw.EmptyClip, cl)
	})

	//// Draw selected text regions. We doo it here because we don't want to draw it in front of text.
	//// Maybe in future I will change text selection algorithm and rework this.
	//for _, region := range wnd.textRegions {
	//	b := region.Min
	//	wnd.buffer.CreateRect(b.X, b.Y, region.Width(), region.Height(), 0, draw.StraightCorners, 0,
	//		softGreen, wnd.DefaultClip())
	//}
	wnd.buffer.CreateTexturedRect(wnd.x, wnd.y, wnd.w, wnd.h, texId, textCoords, whiteColor, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
	wnd.buffer.SeparateBuffer(texId, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
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
	//var wnd *Window
	var cmdw draw.Window_command
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
			c.wantResizeH = true
		} else if utils.PointInRect(c.io.MousePos, vResizeRect) && c.ActiveWindow == wnd {
			c.io.SetCursor(HResizeCursor)
			c.wantResizeV = true
		} else {
			c.io.SetCursor(ArrowCursor)
		}
		c.dragBehavior(vResizeRect, &wnd.capturedV)
		c.dragBehavior(hResizeRect, &wnd.capturedH)
		// Изменение размеров окна
		if c.wantResizeH && c.ActiveWindow == wnd && wnd.capturedH {
			n := newH
			n += c.io.MouseDelta.Y
			if n > wnd.minH {
				newH = n
				if wnd.mainWidgetSpace.scrlY != 0 {
					wnd.mainWidgetSpace.scrlY -= c.io.MouseDelta.Y
				}
			}
		} else if c.wantResizeV && c.ActiveWindow == wnd && wnd.capturedV {
			n := newW
			n += c.io.MouseDelta.X
			if n > wnd.minW {
				newW = n
			}
		}

		c.dragBehavior(wnd.outerRect, &wnd.capturedWin)
		// Изменение положения окна
		if c.ActiveWindow == wnd && wnd.capturedWin && !c.wantResizeV &&
			!c.wantResizeH && !wnd.capturedInsideWin && !wnd.capturedTextSelection {
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

	{
		cmdw = draw.Window_command{
			Active: wnd.active,
			Id:     wnd.Id,
			X:      wnd.x,
			Y:      wnd.y,
			H:      wnd.h,
			W:      wnd.w,
			Clr:    mainWindowClr2,
			Toolbar: draw.Toolbar_command{
				X:   newX,
				Y:   newY,
				W:   newW,
				H:   wnd.toolbar.h,
				Clr: wnd.toolbar.clr,
			},
		}
	}
	wnd.mainWidgetSpace.cursorX = wnd.mainWidgetSpace.X + uiCtx.CurrentStyle.LeftMargin
	wnd.mainWidgetSpace.cursorY = wnd.mainWidgetSpace.Y + uiCtx.CurrentStyle.TopMargin

	wnd.mainWidgetSpace.ClipRect = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y,
		wnd.mainWidgetSpace.W - wnd.mainWidgetSpace.verticalScrollbar.w, wnd.mainWidgetSpace.H}

	wnd.buffer.CreateWindow(cmdw, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
	wnd.widgetSpaceLogic(wnd.mainWidgetSpace, func() draw.ClipRectCompose {
		cl := [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}
		return draw.NewClip(draw.EmptyClip, cl)
	})

	// Draw selected text regions. We doo it here because we don't want to draw it in front of text.
	// Maybe in future I will change text selection algorithm and rework this.
	for _, region := range wnd.textRegions {
		b := region.Min
		wnd.buffer.CreateRect(b.X, b.Y, region.Width(), region.Height(), 0, draw.StraightCorners, 0,
			softGreen, wnd.DefaultClip())
	}
	c.windowStack.Push(wnd)

	if wnd.toolbar.h != 0 {
		b, hovered, clicked := c.button(windowName+"-btn", wnd, wnd.x+wnd.w-25, wnd.y, 25, 25)
		if hovered {
			b.SetColor(softGreen)
		} else {
			b.SetColor(red)
		}
		if clicked {
			*opened = !*opened
		}
		//box := b.BoundingBox()
		wnd.buffer.CreateRect(wnd.x+wnd.w-25, wnd.y, 25, 25, 0, draw.StraightCorners,
			0, b.Color(), draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
		wnd.endWidget(wnd.x+wnd.w-25, wnd.y, false, b)
		txt, _ := c.textHelper(windowName, wnd.x, wnd.y, 10, windowName, DefaultTextFlag)
		wnd.buffer.CreateText(wnd.x, wnd.y, txt, *c.font, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
	}
}

var step float32 = 100

func (wnd *Window) widgetSpaceLogic(ws *WidgetSpace, clip func() draw.ClipRectCompose) {
	c := uiCtx

	if c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect)) {
		c.ActiveWidgetSpaceId = ws.id
		if ws.flags&Scrollable != 0 {
			c.WantScrollFocusWidgetSpaceId = ws.id
		}
	}
	// Scrollbar behavior
	if ws.flags&Scrollable != 0 {
		ws.vertScrollBar()
		if c.ActiveWindow == wnd && c.HoveredWindow == wnd && c.io.ScrollY != 0 && c.WantScrollFocusWidgetSpaceLastId == ws.id && c.FocusedWidgetSpace == nil {
			ws.handleMouseScroll(float32(c.io.ScrollY))
		}
		if ws.flags&ShowScrollbar != 0 && ws.isVertScrollShown {
			wnd.addDelayedWidget(func() {
				scrlClip := clip()
				scrl := ws.verticalScrollbar
				wnd.buffer.CreateRect(scrl.x, scrl.y, scrl.w, scrl.h, 5, draw.AllRounded, 0, scrl.clr, scrlClip)
				wnd.buffer.CreateRect(scrl.bX, scrl.bY, scrl.bW, scrl.bH, 5, draw.AllRounded, 0, [4]float32{255, 0, 0, 1}, scrlClip)
				wnd.buffer.SeparateBuffer(0, scrlClip)
			})
		}
	}

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

func ButtonT(id string, msg string) bool {
	c := ctx()
	wnd := c.windowStack.Peek()
	var tBtn *widgets.TextButton
	var hovered, clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	tBtn, hovered, clicked = c.textButton(id, wnd, msg, x, y, widgets.Center)
	if hovered {
		tBtn.Button.SetColor(c.CurrentStyle.BtnHoveredColor)
	} else if tBtn.Active() {
		tBtn.Button.SetColor(c.CurrentStyle.BtnActiveColor)
	} else {
		tBtn.Button.SetColor(c.CurrentStyle.BtnColor)
	}

	clip := wnd.endWidget(x, y, isRow, tBtn)
	wnd.buffer.CreateButtonT(x, y, tBtn, *c.font, clip)

	return clicked
}

func Slider(id string, i *float32, min, max float32) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var slider *widgets.Slider
	//var hovered, clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	slider = c.getWidget(id, func() widgets.Widget {
		return widgets.NewSlider(id, x, y, 200, 50, min, max, c.CurrentStyle)
	}).(*widgets.Slider)

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

	//clip := draw.NewClip(slider.BoundingBox(), wnd.currentWidgetSpace.ClipRect)

	clip := wnd.endWidget(x, y, isRow, slider)

	wnd.buffer.CreateRect(slider.MainSliderPos()[0], slider.MainSliderPos()[1], slider.MainSliderPos()[2], slider.MainSliderPos()[3], 0,
		draw.StraightCorners, 0, softGreen,
		clip)

	wnd.buffer.CreateRect(slider.BtnSliderPos()[0], slider.BtnSliderPos()[1], slider.BtnSliderPos()[2], slider.BtnSliderPos()[3], 0,
		draw.StraightCorners, 0, black,
		clip)
	wnd.buffer.SeparateBuffer(0, clip)
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
	wnd.textRegions = []utils.Rect{}
}

func (c *UiContext) tHelper(id string, x, y, w float32, msg string, flag widgets.TextFlag) (txt *widgets.Text, hovered bool) {

	wnd := c.windowStack.Peek()
	txt = c.getWidget(id, func() widgets.Widget {
		//if flag&widgets.SplitChars != 0 {
		//	msg = c.FitTextToWidth(x, w, msg)
		//} else if flag&widgets.SplitWords != 0 {
		//	numChars := int(math.Floor(float64(w / c.font.XCharAdvance())))
		//	msg = wrap(msg, numChars)
		//}
		width, h, l, chars := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewText(id, msg, x, y, width, h, chars, l, c.CurrentStyle, flag)

	}).(*widgets.Text)
	hovered = c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox()))
	if txt.LastWidth != w {
		if flag&widgets.SplitChars != 0 {
			msg = c.FitTextToWidth(x, w, msg)
		} else if flag&widgets.SplitWords != 0 {
			numChars := int(math.Floor(float64(w / c.font.XCharAdvance())))
			msg = wrap(msg, numChars)
		}
		width, h, l, chars := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		txt.Lines = l
		txt.Chars = chars
		txt.Message = msg
		txt.SetWH(width, h)
		txt.LastWidth = w
	}
	return
}

func (c *UiContext) inputTextHelper(id string, x, y, w float32, msg string, inpmsg *string, key GuiKey, flag widgets.TextFlag) (txt *widgets.Text, hovered bool) {
	wnd := c.windowStack.Peek()
	txt = c.getWidget(id, func() widgets.Widget {
		width, h, l, chars := c.font.CalculateTextBounds(*inpmsg, c.CurrentStyle.FontScale)
		return widgets.NewText(id, *inpmsg, x, y, width, h, chars, l, c.CurrentStyle, flag)
	}).(*widgets.Text)
	hovered = c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox()))
	if key != GuiKey_None && flag&widgets.Editable != 0 && c.FocusedTextInput == txt {
		if key == GuiKey_Backspace {
			if c.SelectedText != "" {
				*inpmsg = strings.ReplaceAll(*inpmsg, c.SelectedText, "")
				c.ClearTextSelection(wnd)
			} else {
				tmp := ""
				if txt.CursorInd == 0 && txt.CursorLine == 0 {
				} else {
					txt.CursorHelper(-1)
					tmp = removeFromString(*inpmsg, txt.CursorLine, txt.CursorInd, txt.Lines)
					*inpmsg = tmp
				}
			}
		} else if key == GuiKey_RightArrow {
			txt.CursorHelper(1)
		} else if key == GuiKey_LeftArrow {
			txt.CursorHelper(-1)
		} else if key == GuiKey_Enter {
			if txt.Flag&widgets.MultiLine != 0 {
				tmp := insertIntoString2(*inpmsg, txt.CursorLine, txt.CursorInd, msg, txt.Lines)
				txt.CursorLine++
				txt.CursorInd = 0
				*inpmsg = tmp
			}
		} else if IsCommandKey(key) {

		} else if key != GuiKey_None {
			tmp := insertIntoString2(*inpmsg, txt.CursorLine, txt.CursorInd, msg, txt.Lines)
			txt.CursorInd++
			*inpmsg = tmp
		}
		width, h, l, chars := c.font.CalculateTextBounds(*inpmsg, c.CurrentStyle.FontScale)
		txt.Lines = l
		txt.Chars = chars
		txt.SetWH(width, h)
		txt.Message = *inpmsg
		txt.LastWidth = w
	}
	return
}
func (c *UiContext) textHelper(id string, x, y, w float32, msg string, flag widgets.TextFlag) (txt *widgets.Text, hovered bool) {
	return c.tHelper(id, x, y, w, msg, flag)
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
	wnd := c.windowStack.Peek()
	var txtGl *widgets.Text
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws := c.subWidgetSpaceHelperWithBackground(id, x, y, w, h, 0, 0, softGreen, draw.StraightCorners, Scrollable|FitWidth, func() {
		currWs := wnd.currentWidgetSpace
		x, y, isRow := currWs.getCursorPosition()
		msg, key := c.getTextInput()
		txt, hovered := c.inputTextHelper(id, x, y, currWs.W, msg, message, key, widgets.Editable|Selectable)
		txtGl = txt
		if y+txt.Height() > wnd.y && y <= wnd.y+wnd.h {
			wnd.VisibleTexts = append(wnd.VisibleTexts, txt)
		}
		y += wnd.currentWidgetSpace.resolveRowAlign(txt.Height())

		if hovered {

		} else {
			txt.CurrentColor = whiteColor
			txt.SetBackGroundColor(transparent)
		}

		wnd.endWidget(x, y, isRow, txt)
		ws := wnd.currentWidgetSpace
		if c.FocusedTextInput == txt {
			xc, yc, wc, hc := txt.CalculateCursorPos()
			wnd.buffer.CreateRect(x+xc, y+yc, wc, hc, 0,
				draw.StraightCorners, 0, red, wnd.DefaultClip())
		}
		//wnd.buffer.CreateRect(ws.X, ws.Y, ws.W, ws.H, 0, draw.StraightCorners, 0, whiteColor, wnd.DefaultClip())
		wnd.buffer.CreateText(x, y, txt, *c.font, draw.NewClip(draw.EmptyClip, ws.ClipRect))
	})
	if c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect)) && c.io.MouseClicked[0] {
		txtGl.CursorInd = len(txtGl.Chars)
		c.FocusedTextInput = txtGl
		pos := c.io.MouseClickedPos[0]
		startFounded := false
		for _, line := range txtGl.Lines {
			if pos.Y > line.StartY+y && pos.Y <= line.StartY+y+line.Height && !startFounded {
				startFounded = true
				txtGl.CursorInd = len(line.Text)
			}
		}
	}

	wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	wnd.addCursor(ws.W, ws.H)
}

// MultiLineTextInput TODO: Add want input flag
func MultiLineTextInput(id string, message *string) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var txtGl *widgets.Text
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws := c.subWidgetSpaceHelperWithBackground(id, x, y, wnd.mainWidgetSpace.W-(x-wnd.x)-wnd.mainWidgetSpace.verticalScrollbar.w, 200, 0, 0, softGreen, draw.StraightCorners, Scrollable|ShowScrollbar|FitWidth, func() {
		currWs := wnd.currentWidgetSpace
		x, y, isRow := currWs.getCursorPosition()
		msg, key := c.getTextInput()
		txt, hovered := c.inputTextHelper(id, x, y, currWs.W, msg, message, key, widgets.Editable|Selectable|widgets.MultiLine)
		txtGl = txt
		if y+txt.Height() > wnd.y && y <= wnd.y+wnd.h {
			wnd.VisibleTexts = append(wnd.VisibleTexts, txt)
		}
		y += wnd.currentWidgetSpace.resolveRowAlign(txt.Height())

		if hovered {
			//txt.SetBackGroundColor(softGreen)
			//txt.CurrentColor = [4]float32{167, 200, 100, 1}
		} else {
			txt.CurrentColor = whiteColor
			txt.SetBackGroundColor(transparent)
		}

		//txt.CurrentColor = [4]float32{255, 255, 255, 1}

		wnd.endWidget(x, y, isRow, txt)
		ws := wnd.currentWidgetSpace
		if c.FocusedTextInput == txt {
			xc, yc, wc, hc := txt.CalculateCursorPos()
			wnd.buffer.CreateRect(x+xc, y+yc, wc, hc, 0,
				draw.StraightCorners, 0, red, wnd.DefaultClip())
		}
		//wnd.buffer.CreateRect(ws.X, ws.Y, ws.W, ws.H, 0, draw.StraightCorners, 0, whiteColor, wnd.DefaultClip())
		wnd.buffer.CreateText(x, y, txt, *c.font, draw.NewClip(draw.EmptyClip, ws.ClipRect))
	})
	if c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect)) && c.io.MouseClicked[0] {
		txtGl.CursorInd = len(txtGl.Chars)
		c.FocusedTextInput = txtGl
		pos := c.io.MouseClickedPos[0]
		startFounded := false
		for _, line := range txtGl.Lines {
			if pos.Y > line.StartY+y && pos.Y <= line.StartY+y+line.Height && !startFounded {
				startFounded = true
				txtGl.CursorInd = len(line.Text)
			}
		}
	}

	wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	wnd.addCursor(ws.W, ws.H)
}

//func TextFitted(id string, w float32, msg string) {
//	c := ctx()
//	//wnd := c.windowStack.Peek()
//	//x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
//	//txt, hovered := c.textHelper(id, x, y, w, msg, widgets.SplitWords|widgets.Selectable)
//	//if y+txt.Height() > wnd.y && y <= wnd.y+wnd.h {
//	//	wnd.VisibleTexts = append(wnd.VisibleTexts, txt)
//	//}
//	//y += wnd.currentWidgetSpace.resolveRowAlign(txt.Height())
//	txt := c.getWidget(id, func() widgets.Widget {
//		width, h, l, chars := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
//		return widgets.NewText(id, msg, 0, 0, width, h, chars, l, c.CurrentStyle, widgets.SplitWords|widgets.Selectable)
//	}).(*widgets.Text)
//	wnd, x, y, isRow, out := CalculateWidgetInfo(txt.Width(), txt.Height())
//	//txt.UpdatePosition([4]float32{x, y, txt.Width(), txt.Height()})
//	if out || wnd == nil {
//		return
//	}
//	hovered := c.TextEX(txt, msg, w)
//	wnd.debugDrawS(txt.BoundingBox())
//	if hovered {
//
//	} else {
//		txt.CurrentColor = whiteColor
//		txt.SetBackGroundColor(transparent)
//	}
//
//	clip := wnd.endWidget(x, y, isRow, txt)
//
//	wnd.buffer.CreateText(x, y, txt, *c.font, clip)
//}

//func Text(id string, msg string, flag widgets.TextFlag) {
//	c := ctx()
//	txt := c.getWidget(id, func() widgets.Widget {
//		width, h, l, chars := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
//		return widgets.NewText(id, msg, 0, 0, width, h, chars, l, c.CurrentStyle, flag)
//	}).(*widgets.Text)
//	wnd, x, y, isRow, _ := CalculateWidgetInfo(txt.Width(), txt.Height())
//	//txt.UpdatePosition([4]float32{x, y, txt.Width(), txt.Height()})
//	//if out || wnd == nil {
//	//	return
//	//}
//	hovered := c.TextEX(txt, msg, txt.Width())
//	//wnd := c.windowStack.Peek()
//	//x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
//	//txt, hovered := c.textHelper(id, x, y, 0, msg, flag)
//	//if y+txt.Height() > wnd.y && y <= wnd.y+wnd.h {
//	//	wnd.VisibleTexts = append(wnd.VisibleTexts, txt)
//	//}
//	//y += wnd.currentWidgetSpace.resolveRowAlign(txt.Height())
//
//	//wnd.buffer.RoundedBorderRectangle(x, y, txt.Width(), txt.Height(), 30, 15, red, wnd.DefaultClip())
//	if hovered {
//		//txt.SetBackGroundColor(softGreen)
//		//txt.CurrentColor = [4]float32{167, 200, 100, 1}
//	} else {
//		txt.CurrentColor = whiteColor
//		txt.SetBackGroundColor(transparent)
//	}
//
//	//txt.CurrentColor = [4]float32{255, 255, 255, 1}
//
//	clip := wnd.endWidget(x, y, isRow, txt)
//
//	wnd.buffer.CreateText(x, y, txt, *c.font, clip)
//}

func (wnd *Window) DefaultClip() draw.ClipRectCompose {
	return draw.NewClip(wnd.currentWidgetSpace.ClipRect, wnd.mainWidgetSpace.ClipRect)
}

func (c *UiContext) imageHelper(id string, x, y, w, h float32, f func() *widgets.Image) (img *widgets.Image, hovered, clicked bool) {
	wnd := c.windowStack.Peek()

	img = f()
	hovered = c.hoverBehavior(wnd, utils.NewRectS(img.BoundingBox()))
	{
		if hovered {
			c.setActiveWidget(img.WidgetId())
		}
		clicked = c.io.MouseClicked[0] && hovered
	}
	return
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

func (c *UiContext) TextEX(ws *WidgetSpace, id string, msg string, newWidth float32, flag widgets.TextFlag) (txt *widgets.Text, x, y float32, isRow, outOfWs bool) {
	txt = c.getWidget(id, func() widgets.Widget {
		width, h, l, chars := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		txt := widgets.NewText(id, msg, ws.cursorX, ws.cursorY, width, h, chars, l, c.CurrentStyle, flag)
		txt.Base.Updated = true
		return txt
	}).(*widgets.Text)
	x, y, isRow, outOfWs = CalculateWidgetInfo(txt.Width(), txt.Height(), ws)
	if outOfWs {
		ws.addCursor(txt.Width(), txt.Height())
		if !isRow {
			ws.AddVirtualWH(txt.Width(), txt.Height())
		}
		return
	}
	if txt.LastWidth != newWidth {
		fmt.Println(newWidth)
		if txt.Flag&widgets.SplitChars != 0 {
			msg = c.FitTextToWidth(txt.BoundingBox()[0], newWidth, msg)
		} else if txt.Flag&widgets.SplitWords != 0 {
			numChars := int(math.Floor(float64(newWidth / c.font.XCharAdvance())))
			msg = wrap(msg, numChars)
		}

		width, h, l, chars := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		txt.Lines = l
		txt.Chars = chars
		txt.Message = msg
		txt.SetWH(width, h)
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

func (c *UiContext) DrawImage(x, y float32, img *widgets.Image, buffer *draw.CmdBuffer, clip draw.ClipRectCompose) {
	if img.Updated {
		img.LastVert, img.LastInd, img.LastVertCount, img.Last = buffer.CreateTexturedRect2(x, y,
			img.Width(), img.Height(), img.TexId, img.TexCoords, img.Color(), clip)
		buffer.Render2(img.LastVert, img.LastInd, img.LastVertCount, 0)
		img.Updated = false
	} else {
		buffer.Render2(img.LastVert, img.LastInd, img.LastVertCount, img.Last)
	}
	buffer.SeparateBuffer(img.TexId, clip)
}

func GlobalImage(id string, x, y, w, h float32, texId uint32, texCoords [4]float32) bool {
	c := ctx()
	var ws *WidgetSpace
	var clicked bool
	ws = c.subWidgetSpaceHelper("wsp-"+id, x, y, w, h, IgnoreClipping, func() {
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
		ws.UpdateWidgetPosition(x, y, isRow, img)

		c.DrawImage(x, y, img, c.globalBuffer, ws.Clip())
	})

	return clicked
}

func Image(id string, w, h float32, texId uint32, texCoords [4]float32) bool {
	c := ctx()

	wnd := c.getPeekWindow()
	ws := wnd.currentWidgetSpace

	img, x, y, isRow, out := c.ImageEX(ws, id, w, h, texId, texCoords, whiteColor)
	if out {
		return false
	}
	//img := c.getWidget(id, func() widgets.Widget {
	//	img2 := widgets.NewImage2(id, ws.cursorX, ws.cursorY, w, h, texId, texCoords, whiteColor)
	//	img2.Updated = true
	//	return img2
	//}).(*widgets.Image)
	//
	//x, y, isRow, out := CalculateWidgetInfo(w, h, ws)
	//if out {
	//	ws.addCursor(w, h)
	//	if !isRow {
	//		ws.AddVirtualWH(w, h)
	//	}
	//	return false
	//}

	//hov, clicked := c.ImageEX(img)
	hovered := c.hoverBehavior(wnd, utils.NewRectS(img.BoundingBox()))
	clicked := hovered && c.io.MouseClicked[0]
	//img, hov, clicked := c.ImageEX(id, x, y, w, h, texId, texCoords, whiteColor, wnd.buffer, wnd.DefaultClip())
	//wnd.debugDrawS([4]float32{x, y, w, h})
	if hovered {
		img.SetColor(red)
	} else {
		img.SetColor(whiteColor)
	}
	ws.UpdateWidgetPosition(x, y, isRow, img)

	c.DrawImage(x, y, img, wnd.buffer, wnd.DefaultClip())
	return clicked
}

func Text(id string, msg string, flag widgets.TextFlag) {
	c := ctx()
	wnd := c.getPeekWindow()
	ws := wnd.currentWidgetSpace
	//txt := c.getWidget(id, func() widgets.Widget {
	//	width, h, l, chars := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
	//	txt := widgets.NewText(id, msg, ws.cursorX, ws.cursorY, width, h, chars, l, c.CurrentStyle, flag)
	//	txt.Base.Updated = true
	//	return txt
	//}).(*widgets.Text)
	//
	//x, y, isRow, out := CalculateWidgetInfo(txt.Width(), txt.Height(), ws)
	//if out {
	//	ws.addCursor(txt.Width(), txt.Height())
	//	if !isRow {
	//		ws.AddVirtualWH(txt.Width(), txt.Height())
	//	}
	//	return
	//}
	txt, x, y, isRow, out := c.TextEX(ws, id, msg, 0, flag)
	if out {
		return
	}

	hovered := c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox()))
	if hovered {
		txt.SetTextColor(softGreen)
	} else {
		txt.SetTextColor(c.CurrentStyle.TextColor)
	}
	ws.UpdateWidgetPosition(x, y, isRow, txt)
	c.DrawText(x, y, txt, c.font.TextureId, wnd.buffer, wnd.DefaultClip())
}
func (c *UiContext) DrawText(x, y float32, txt *widgets.Text, texid uint32, buffer *draw.CmdBuffer, clip draw.ClipRectCompose) {
	if txt.Base.Updated {
		txt.Base.LastVert, txt.Base.LastInd, txt.Base.LastVertCount, txt.Base.Last = buffer.CreateText2(x, y, txt, 1, *c.font)
		buffer.Render2(txt.Base.LastVert, txt.Base.LastInd, txt.Base.LastVertCount, 0)
		txt.Base.Updated = false
	} else {
		buffer.Render2(txt.Base.LastVert, txt.Base.LastInd, txt.Base.LastVertCount, txt.Base.Last)
	}
	buffer.SeparateBuffer(texid, clip)
}
func (ws *WidgetSpace) UpdateWidgetPosition(xPos, yPos float32, isRow bool, w widgets.Widget) {
	w.UpdatePosition([4]float32{xPos, yPos, w.Width(), w.Height()})
	ws.addCursor(w.Width(), w.Height())
	if !isRow {
		ws.AddVirtualWH(w.Width(), w.Height())
	}
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
	w.delayedWidgets = append(w.delayedWidgets, f)
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
	inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y, wnd.w, wnd.h)
	//inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y+wnd.toolbar.h, wnd.w, wnd.h-wnd.toolbar.h)
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
	wnd := c.windowStack.Peek()
	var tBtn *widgets.TextButton
	var _, hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	tBtn = c.getWidget(id, func() widgets.Widget {
		w, h, _, p := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewTextButton(id, x, y, w, h, msg, p, widgets.Left, widgets.AllPadding, c.CurrentStyle)
	}).(*widgets.TextButton)

	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.BoundingBox()))
		if hovered {
			c.setActiveWidget(tBtn.Id)
			tBtn.Button.SetColor(c.CurrentStyle.BtnHoveredColor)
			if c.io.MouseClicked[0] {
				tBtn.ChangeActive()
			}
		} else if tBtn.Active() {
			tBtn.Button.SetColor(c.CurrentStyle.BtnActiveColor)
		} else {
			tBtn.Button.SetColor(c.CurrentStyle.BtnColor)
		}
		//clicked = c.io.MouseClicked[0] && hovered
	}
	//
	tBtn.SetWidth(wnd.w)
	clip := wnd.endWidget(x, y, isRow, tBtn)
	wnd.buffer.CreateButtonT(x, y, tBtn, *c.font, clip)

	if tBtn.Active() {
		x += 50
		ws := c.subWidgetSpaceHelper(id, x, y+tBtn.Height(), 0, 0, NotScrollable|Resizable, widgFunc)
		wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
		wnd.addCursor(ws.W, ws.H)
	}

	return tBtn.Active()
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

func (c *UiContext) button(id string, wnd *Window, x, y, w, h float32) (btn *widgets.Button, hovered, clicked bool) {
	btn = c.getWidget(id, func() widgets.Widget {
		return widgets.NewButton(id, x, y, w, h, c.CurrentStyle.BtnColor)
	}).(*widgets.Button)

	hovered = c.hoverBehavior(wnd, utils.NewRectS(btn.BoundingBox()))
	if hovered {
		c.setActiveWidget(btn.WidgetId())
	}
	clicked = c.io.MouseClicked[0] && hovered
	if clicked {
		btn.ChangeActive()
	}

	return
}

func Button(id string) bool {
	c := ctx()
	wnd := c.windowStack.Peek()
	var btn *widgets.Button
	var clicked, hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	btn, hovered, clicked = c.button(id, wnd, x, y, 100, 100)

	if hovered {
		btn.SetColor(c.CurrentStyle.BtnHoveredColor)
	} else if btn.IsActive {
		btn.SetColor(c.CurrentStyle.BtnActiveColor)
	} else {
		btn.SetColor(c.CurrentStyle.BtnColor)
	}

	clip := wnd.endWidget(x, y, isRow, btn)
	wnd.buffer.CreateRect(x, y, btn.Width(), btn.Height(), 0, draw.StraightCorners, 0, btn.Color(), clip)

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
	wnd := c.windowStack.Peek()
	var s *widgets.Selection

	// Need to use WS because text may not fit into button, so it should be clipped
	SubWidgetSpace(id+"---", 0, 0, Resizable|NotScrollable, func() {
		Row(id+"row--", func() {
			x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
			s = c.getWidget(id, func() widgets.Widget {
				return widgets.NewSelection(id, x, y, 300, 40)
			}).(*widgets.Selection)

			wnd.buffer.CreateRect(x, y, s.Width(), s.Height(), 0, draw.StraightCorners, 0, whiteColor, wnd.DefaultClip())
			wnd.buffer.SeparateBuffer(0, wnd.DefaultClip())
			wnd.endWidget(x, y, isRow, s)

			x2, y2, isRow2 := wnd.currentWidgetSpace.getCursorPosition()
			img, _, clicked := c.imageHelper(id+"arrow", x2, y2, s.Height(), s.Height(), func() *widgets.Image {
				return nil
			})

			wnd.endWidget(x2-s.Height(), y2, isRow2, img)

			txt, _ := c.textHelper(data[*index]+"--"+id, x, y, 0, data[*index], widgets.Default)
			wnd.buffer.CreateText(x+c.CurrentStyle.Padding, y+(s.Height()-txt.Height())/2, txt,
				*c.font, draw.NewClip(draw.EmptyClip, [4]float32{x, y, s.Width(), s.Height()}))
			wnd.buffer.CreateTexturedRect(x2-s.Height(), y2, img.Width(), img.Height(), texId, texCoords, img.Color(), wnd.DefaultClip())
			if clicked {
				s.Opened = true
				c.setActiveWidget(id)
			}
			if c.ActiveWidget != id {
				s.Opened = false
			}
		})
	})

	ContextMenu(id, IgnoreClipping, func() {
		for i, datum := range data {
			x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
			tbt, _, clicked := c.textButton(datum+"_btnT_"+id, wnd, datum, x, y, widgets.Left)
			if clicked {
				*index = i
				c.FocusedWidgetSpace = nil
			}
			tbt.SetWidth(s.Width())
			clip := wnd.endWidget(x, y, false, tbt)
			wnd.buffer.CreateButtonT(x, y, tbt, *c.font, clip)
		}
	})
}

func ContextMenu(ownerWidgetId string, flag WidgetSpaceFlag, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var bb [4]float32
	widg, ok := c.GetWidget(ownerWidgetId)
	if !ok {
		return
	}
	bb = widg.BoundingBox()
	id := ownerWidgetId + "-ws-context"
	ws := c.getWidgetSpace(id, 0, 0, wnd, Resizable|FitWidth|flag)
	if c.LastActiveWidget == widg.WidgetId() {
		c.FocusedWidgetSpace = ws
	}
	if c.FocusedWidgetSpace == ws {
		f := func() {
			ws.ClipRect = [4]float32{ws.X, ws.Y, ws.W, ws.H}
			clip := draw.NewClip(draw.EmptyClip, ws.ClipRect)
			wnd.buffer.CreateRect(bb[0], bb[1]+widg.Height(), ws.W, ws.H, 0, draw.StraightCorners, 0, black, clip)
			c.subWidgetSpaceHelper(id, bb[0], bb[1]+widg.Height(), widg.Width(), 0, Resizable|FitWidth|flag, widgFunc)
		}
		wnd.addDelayedWidget(f)
	}
}
func Tooltip(id string, widgFunc func()) {
	c := ctx()
	x, y := c.io.MousePos.X+10, c.io.MousePos.Y+5
	wnd := c.windowStack.Peek()

	ws := c.getWidgetSpace(id, 0, 0, wnd, Resizable|IgnoreClipping)

	wnd.addDelayedWidget(func() {
		wnd.buffer.CreateRect(x, y, ws.W, ws.H, 0, draw.StraightCorners, 0, black, draw.NewClip(draw.EmptyClip, ws.ClipRect))
		c.subWidgetSpaceHelper(id, x, y, 0, 0, Resizable|IgnoreClipping, widgFunc)
	})
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

func (c *UiContext) getWidgetSpace(id string, width, height float32, wnd *Window, flags WidgetSpaceFlag) *WidgetSpace {
	ws, ok := c.widgSpaceCache.Get(id)
	if !ok {
		ws = newWidgetSpace(id, 0, 0, width, height, flags)
		c.widgSpaceCache.Add(id, ws)
	}
	wnd.widgSpaces = append(wnd.widgSpaces, ws)
	return ws
}

func (c *UiContext) subWidgetSpaceHelperEx(id string, x, y, width, height float32, texId uint32, radius int, clr [4]float32, shape draw.RoundedRectShape, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	wnd := c.windowStack.Peek()

	ws := c.getWidgetSpace(id, width, height, wnd, flags)

	var prevWS = wnd.currentWidgetSpace
	wnd.currentWidgetSpace = ws

	wnd.debugDraw(ws.X, ws.Y, ws.W, ws.H, red)

	ws.X = x
	ws.Y = y
	ws.cursorY = y
	ws.cursorX = x

	outOfWindow := false
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
		ws.ClipRect = [4]float32{x, y, ws.W, (wnd.mainWidgetSpace.Y - y)}
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

	wnd.widgetSpaceLogic(ws, func() draw.ClipRectCompose {
		cl := [4]float32{ws.X, ws.Y, ws.W, ws.H}
		if outOfWindow {
			cl[1] = wnd.mainWidgetSpace.Y
		}
		return draw.NewClip(cl, wnd.mainWidgetSpace.ClipRect)
	})

	if flags&FillBackground != 0 {
		if texId == 0 {
			wnd.buffer.CreateRect(ws.X, ws.Y, ws.W, ws.H, radius, shape, 0, clr, wnd.DefaultClip())
		} else {
			//TODO: Add textured rect method
			//wnd.buffer.CreateTexturedRect()
		}
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

	wnd.currentWidgetSpace = prevWS
	return ws
}

func (c *UiContext) subWidgetSpaceHelperWithBackground(id string, x, y, width, height float32, texId uint32, radius int, clr [4]float32, shape draw.RoundedRectShape, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	return c.subWidgetSpaceHelperEx(id, x, y, width, height, texId, radius, clr, shape, flags|FillBackground, widgFunc)
}

func (c *UiContext) subWidgetSpaceHelper(id string, x, y, width, height float32, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	return c.subWidgetSpaceHelperEx(id, x, y, width, height, 0, 0, transparent, draw.StraightCorners, flags, widgFunc)
}

func SubWidgetSpace(id string, width, height float32, flags WidgetSpaceFlag, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var ws *WidgetSpace

	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws = c.subWidgetSpaceHelper(id, x, y, width, height, flags, widgFunc)

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
		ws = c.subWidgetSpaceHelper(wspId, x, y, 0, 0, Resizable|HideScrollbar, widgFunc)
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
	ws := c.subWidgetSpaceHelper(id, x, y, 0, 0, Resizable|NotScrollable, func() {
		//cr := wnd.currentWidgetSpace
		//wnd.buffer.CreateRect(cr.X, cr.Y, cr.W, cr.H, 10, draw.AllRounded, 0, softGreen, wnd.DefaultClip())
		Row("rowds", func() {
			row, _ := wnd.currentWidgetSpace.getCurrentRow()
			wnd.buffer.CreateRect(row.X, row.Y, row.Width(), row.Height(), 10, draw.AllRounded, 0, transparent, wnd.DefaultClip())
			for i, item := range tab.Bars {
				x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
				tbtn, hovered, clicked := c.textButton(fmt.Sprint(id, "-", i), wnd, item.Name, x, y, widgets.Center)
				if clicked {
					tab.CurrentTab = i
				}
				if hovered {
					tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnActiveColor)
					if clicked {
						tab.ChangeActive(item)
					}
				} else if item.Active {
					tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnActiveColor)
				} else {
					tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnColor)
				}

				tbtn.Text.CurrentColor = whiteColor
				tbtn.SetHeight(tbtn.Height() - (tbtn.Height() - tbtn.Text.Height()) + c.CurrentStyle.Padding)
				clip := wnd.endWidget(x, y, false, tbtn)
				//wnd.buffer.CreateRect(x, y, tbtn.Width(), tbtn.Height(), 10, draw.TopRect, 0, tbtn.Color(), clip)
				wnd.buffer.CreateRect(x, y, tbtn.Width(), tbtn.Height(), 10, draw.TopRect, 0, tbtn.Color(), clip)
				wnd.buffer.SeparateBuffer(0, clip)
				wnd.buffer.CreateText(tbtn.Text.BoundingBox()[0], tbtn.Text.BoundingBox()[1], tbtn.Text, *c.font, clip)

				row, _ := wnd.currentWidgetSpace.getCurrentRow()
				if i != len(tab.Bars)-1 {
					row.CursorX += 10
					row.W += 10
				}
				if row.Height() > rowHeight {
					rowHeight = row.Height()
				}
				rowWidth = row.Width()
			}
		})

		wnd.buffer.CreateRect(wnd.x, y+rowHeight, wnd.w, 2, 0, draw.StraightCorners, 0, c.CurrentStyle.TabBtnActiveColor, draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect))
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

	wnd.buffer.CreateRect(wnd.mainWidgetSpace.X, ws.H+y, wnd.w, 2,
		0, draw.StraightCorners, 0, c.CurrentStyle.WidgSpaceDividerColor, wnd.DefaultClip())

	wnd.addCursor(ws.W, ws.H)
	wnd.currentWidgetSpace.AddVirtualWH(ws.W, ws.H)
}

func Row(id string, widgFunc func()) {
	c := ctx()
	wnd := c.windowStack.Peek()
	var row *widgets.HybridLayout
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	// Need to return cursor back, because internal row cursor shouldn't know anything about outer
	y += wnd.currentWidgetSpace.scrlY

	row = c.getWidget(id, func() widgets.Widget {
		return widgets.NewHLayout(id, x, y, widgets.VerticalAlign, c.CurrentStyle)
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
				//x := c.io.dragStarted.X + c.io.dragDelta.X - b[0]
				//y := c.io.dragStarted.Y + c.io.dragDelta.Y - b[1]

				for i := 0; i < len(t.Lines); i++ {
					for ind, pos := range t.Lines[i].Text {
						if x >= pos.Pos.X-5 && x <= pos.Pos.X+float32(pos.Char.Width) && !startFounded &&
							y <= t.Lines[i].StartY+t.Lines[i].Height {
							t.StartLine = i
							t.StartInd = ind
							startFounded = true
							c.SelectedTextStart = t
							if t == c.FocusedTextInput {
								t.CursorInd = ind
								t.CursorLine = t.StartLine
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
				for i := 0; i < len(t.Lines); i++ {
					for ind, pos := range t.Lines[i].Text {
						if x >= pos.Pos.X-5 && x <= pos.Pos.X+float32(pos.Char.Width) && !endfounded &&
							y <= t.Lines[i].StartY+t.Lines[i].Height {
							t.EndLine = i
							t.EndInd = ind
							c.SelectedTextEnd = t
							endfounded = true
							//fmt.Println(ind)
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
		if len(tmp) == 1 {
			start := c.SelectedTextStart
			if start.StartLine == start.EndLine {
				r, msg := c.findOneLineTexRegion(start, start.StartLine, start.StartInd, start.EndInd)
				rects = append(rects, r...)
				c.SelectedText = msg
			} else {
				selectedText := ""
				for i := start.StartLine; i <= start.EndLine; i++ {
					if i == start.StartLine {
						r, msg := c.findFirstTexRegion(start, start.StartLine, start.StartInd)
						rects = append(rects, r...)
						selectedText += msg + "\n"
					} else if i == start.EndLine {
						r, msg := c.findLastTexRegion(start, start.EndLine, start.EndInd)
						rects = append(rects, r...)
						selectedText += msg
					} else {
						r, msg := c.findAllTexRegion(start, i, i)
						rects = append(rects, r...)
						selectedText += msg
					}
				}
				c.SelectedText = selectedText
			}
		} else if len(tmp) == 2 {
			selectedText := ""
			start := tmp[0]
			r, msg := c.findFirstTexRegion(start, start.StartLine, start.StartInd)
			rects = append(rects, r...)
			selectedText += msg
			if len(start.Lines) != 1 {
				r, msg = c.findAllTexRegion(start, start.StartLine+1, len(start.Lines)-1)
				rects = append(rects, r...)
				selectedText += msg
			}
			end := tmp[1]
			r, msg = c.findAllTexRegion(end, 0, end.EndLine-1)
			rects = append(rects, r...)
			selectedText += msg
			r, msg = c.findLastTexRegion(end, end.EndLine, end.EndInd)
			rects = append(rects, r...)
			selectedText += msg
			c.SelectedText = selectedText
		} else {
			selectedText := ""
			first := tmp[0]
			r, msg := c.findFirstTexRegion(first, first.StartLine, first.StartInd)
			rects = append(rects, r...)
			selectedText += msg
			if len(first.Lines) != 1 {
				r, msg = c.findAllTexRegion(first, first.StartLine+1, len(first.Lines)-1)
				rects = append(rects, r...)
				selectedText += msg
			}
			for i := 1; i < len(tmp)-1; i++ {
				txt := tmp[i]
				r, msg = c.findAllTexRegion(txt, 0, len(txt.Lines)-1)
				rects = append(rects, r...)
				selectedText += msg
			}
			last := tmp[len(tmp)-1]
			r, msg = c.findAllTexRegion(last, 0, last.EndLine-1)
			rects = append(rects, r...)
			selectedText += msg
			r, msg = c.findLastTexRegion(last, last.EndLine, last.EndInd)
			rects = append(rects, r...)
			selectedText += msg
			c.SelectedText = selectedText
		}
		wnd.textRegions = rects
	}

}

type selectionFlag uint

const (
	firstLine selectionFlag = iota
	lastLine
	innerLines // lines between the first and the last
)

// findOneLineTexRegion is used when you have 1 selected line and it's boundaries are placed between 0 and length of the line.
// For example start at 4 and end at 9 (hell[o wo]rld)
func (c *UiContext) findOneLineTexRegion(txt *widgets.Text, line, startIdx, endIndex int) ([]utils.Rect, string) {
	return c.findTextRegion(txt, line, 0, startIdx, endIndex, firstLine)
}
func (c *UiContext) findFirstTexRegion(txt *widgets.Text, line, startIdx int) ([]utils.Rect, string) {
	return c.findTextRegion(txt, line, 0, startIdx, len(txt.Lines[line].Text)-1, firstLine)
}
func (c *UiContext) findLastTexRegion(txt *widgets.Text, line, endIdx int) ([]utils.Rect, string) {
	return c.findTextRegion(txt, line, 0, 0, endIdx, lastLine)
}
func (c *UiContext) findInnerTexRegion(txt *widgets.Text) ([]utils.Rect, string) {
	return c.findTextRegion(txt, txt.StartLine+1, txt.EndLine-1, 0, 0, innerLines)
}
func (c *UiContext) findAllTexRegion(txt *widgets.Text, startLine, endLine int) ([]utils.Rect, string) {
	return c.findTextRegion(txt, startLine, endLine, 0, 0, innerLines)
}

func (c *UiContext) findTextRegion(txt *widgets.Text, startLine, endLine, startIdx, endIdx int, flag selectionFlag) ([]utils.Rect, string) {
	var x, y, w, h float32
	var msg string
	results := []utils.Rect{}
	b := txt.BoundingBox()
	w = 0
	y = txt.Lines[startLine].StartY + b[1]
	h = txt.Lines[startLine].Height
	switch flag {
	case firstLine:
		x = txt.Lines[startLine].Text[startIdx].Pos.X + b[0]
		for i := startIdx; i <= endIdx; i++ {
			w += txt.Lines[startLine].Text[i].Width
			msg += string(txt.Lines[startLine].Text[i].Char.Rune)
		}
		r := utils.NewRect(x, y, w, h)
		results = append(results, r)
	case lastLine:
		x = txt.Lines[startLine].StartX + b[0]
		for i := startIdx; i <= endIdx; i++ {
			w += txt.Lines[startLine].Text[i].Width
			msg += string(txt.Lines[startLine].Text[i].Char.Rune)
		}
		r := utils.NewRect(x, y, w, h)
		results = append(results, r)
	case innerLines:
		for i := startLine; i <= endLine; i++ {
			x = txt.Lines[i].StartX + b[0]
			y = txt.Lines[i].StartY + b[1]
			w = txt.Lines[i].Width
			h = txt.Lines[i].Height
			msg += txt.Lines[i].Msg + "\n"
			r := utils.NewRect(x, y, w, h)
			results = append(results, r)
		}
	}
	return results, msg
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
		}
	}

	for _, f := range wnd.delayedWidgets {
		f()
	}
	wnd.delayedWidgets = []func(){}

	wnd.VisibleTexts = []*widgets.Text{}

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
