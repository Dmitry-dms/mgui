package ui

import (
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/Dmitry-dms/mgui/widgets"
)

type WidgetSpaceFlag uint32

const (
	ShowScrollbar WidgetSpaceFlag = 1 << iota
	Resizable
	NotResizable
	HideScrollbar
	NotScrollable
	Scrollable
	IgnoreClipping
	FitWidth
	FillBackground

	Default = Resizable | ShowScrollbar | Scrollable
)

type WidgetSpace struct {
	X, Y, W, H float32
	id         string

	//inner widgets
	cursorX, cursorY float32
	widgetCounter    int
	//widgets                          []widgets.Widget
	virtualHeight, lastVirtualHeight float32
	virtualWidth, lastVirtualWidth   float32
	ClipRect                         [4]float32

	verticalScrollbar *Scrollbar
	captured          bool

	isVertScrollShown bool    // used to decide create draw command or not
	scrlY             float32 // main scroll state

	ratio float32

	rowStack utils.Stack[*widgets.HybridLayout]
	tabStack utils.Stack[*widgets.TabBar]

	flags WidgetSpaceFlag
}

var defScrollWidth float32 = 10

func (ws *WidgetSpace) vertScrollBar() {
	vB := ws.verticalScrollbar
	vB.x = ws.X + ws.W - vB.w
	vB.y = ws.Y
	vB.h = ws.H

	ws.ratio = ws.H / ws.lastVirtualHeight

	// this is here because we need to update scroll btn position every frame
	ws.handleMouseDrag()
}

func newWidgetSpace(id string, x, y, w, h float32, flags WidgetSpaceFlag) *WidgetSpace {
	vs := WidgetSpace{
		id:            id,
		X:             x,
		Y:             y,
		W:             w,
		H:             h,
		cursorX:       x,
		cursorY:       y,
		widgetCounter: 0,
		//widgets:       []widgets.Widget{},
		flags:         flags,
		virtualHeight: 0,
		verticalScrollbar: NewScrolBar(utils.NewRect(x+w-defScrollWidth, y, defScrollWidth, h),
			utils.NewRect(x+w-defScrollWidth, y, defScrollWidth, 50),
			[4]float32{150, 155, 155, 1}),
		rowStack: utils.NewStack[*widgets.HybridLayout](),
		tabStack: utils.NewStack[*widgets.TabBar](),
	}
	return &vs
}

// FIXME: is this right?
func (ws *WidgetSpace) setScrollY(scrollY float32) {
	if scrollY <= ws.lastVirtualHeight {
		ws.scrlY = scrollY * (1 - ws.ratio)
	}
}

func (ws *WidgetSpace) handleMouseDrag() {
	vB := ws.verticalScrollbar
	if ws.flags&ShowScrollbar != 0 && ws.isVertScrollShown {
		UiCtx.dragBehavior(utils.NewRect(vB.bX, vB.bY, vB.bW, vB.bH), &ws.captured)
		delta := UiCtx.io.MouseDelta.Y
		if ws.captured {
			// Предотвращение неправильного расчета позиции скроллинга при резком перемещении мыши (delta > 70)
			if ws.H+ws.scrlY+delta > ws.lastVirtualHeight {
				ws.scrlY = ws.lastVirtualHeight - ws.H
			} else {
				ws.scrlY += delta
			}

			if ws.scrlY < 0 {
				ws.scrlY = 0
			}
		}

	}
	vB.bH = ws.H * ws.ratio
	vB.bX = ws.X + ws.W - vB.w
	vB.bY = ws.Y + ws.scrlY*ws.ratio

}
func (ws *WidgetSpace) handleMouseScroll(scrollY float32) {
	var factor = -scrollY * step

	botBorder := ws.lastVirtualHeight

	// FIXED: uncorrect scroll position
	if ws.H+ws.scrlY+factor > botBorder {
		ws.scrlY = ws.lastVirtualHeight - ws.H
	} else {
		ws.scrlY += factor * ws.ratio
	}

	if ws.scrlY < 0 {
		ws.scrlY = 0
	}
}

func (ws *WidgetSpace) checkVerScroll() {
	if ws.cursorY > ws.Y+ws.H {
		ws.isVertScrollShown = true
	} else {
		ws.scrlY = 0
		ws.isVertScrollShown = false
	}
}

func (ws *WidgetSpace) AddVirtualHeight(height float32) {
	ws.virtualHeight += height
}
func (ws *WidgetSpace) AddVirtualWH(width, height float32) {
	if ws.virtualWidth == 0 {
		ws.virtualWidth += width
	} else {
		if ws.virtualWidth < width {
			ws.virtualWidth += width - ws.virtualWidth
		}
	}
	ws.virtualHeight += height
}
func (ws *WidgetSpace) getCurrentRow() (*widgets.HybridLayout, bool) {
	if ws.rowStack.Length() == 0 {
		return nil, false
	} else {
		return ws.rowStack.Peek(), true
	}
}

func (ws *WidgetSpace) getCurrentTabBar() (*widgets.TabBar, bool) {
	if ws.tabStack.Length() == 0 {
		return nil, false
	} else {
		return ws.tabStack.Peek(), true
	}
}

func (ws *WidgetSpace) resolveRowAlign(h float32) (yOffset float32) {
	row, ok := ws.getCurrentRow()
	if ok {
		if row.Align&widgets.VerticalAlign != 0 {
			yOffset = (row.Height() - h) / 2
		}
	}
	return
}
func (ws *WidgetSpace) getCursorPosition() (x float32, y float32, isRow bool) {
	row, ok := ws.getCurrentRow()
	if ok {
		x = row.CursorX
		y = row.CursorY
	} else {
		x = ws.cursorX
		y = ws.cursorY
	}
	isRow = ok
	y -= ws.scrlY
	return
}
