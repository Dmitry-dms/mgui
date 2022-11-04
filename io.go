package ui

import (
	"github.com/Dmitry-dms/mgui/utils"
)

// "fmt"

// "fmt"

const (
	MOUSE_BTN_COUNT = 5
)

type Io struct {
	MousePos    utils.Vec2
	DisplaySize *utils.Vec2

	keyPressed          [610]bool
	modPressed          [8]bool
	PressedKey          GuiKey
	KeyPressedThisFrame bool

	// scrollX, scrollY float64
	WantCaptureMouse bool

	MouseDown     [MOUSE_BTN_COUNT]bool
	MouseReleased [MOUSE_BTN_COUNT]bool

	MouseClicked            [MOUSE_BTN_COUNT]bool
	MouseDoubleClicked      [MOUSE_BTN_COUNT]bool
	MouseClickedTime        [MOUSE_BTN_COUNT]float32
	MouseDownDuration       [MOUSE_BTN_COUNT]float32
	MouseDownDurationPrev   [MOUSE_BTN_COUNT]float32
	MouseDragMaxDistanceSqr [MOUSE_BTN_COUNT]float32 //Squared maximum distance of how much mouse has traveled from the clicking point (used for moving thresholds)
	MouseClickedPos         [MOUSE_BTN_COUNT]utils.Vec2

	MouseClickedCount     [MOUSE_BTN_COUNT]uint8 // == 0 (not clicked), == 1 (same as MouseClicked[]), == 2 (double-clicked), == 3 (triple-clicked)
	MouseClickedLastCount [MOUSE_BTN_COUNT]uint8 //Count successive number of clicks. Stays valid after mouse release. Reset after another click is done.

	IsDragging      bool
	dragDelta       utils.Vec2
	dragStarted     utils.Vec2
	dragStartedMain utils.Vec2

	//refactor
	MousePosPrev utils.Vec2
	MouseDelta   utils.Vec2

	DeltaTime               float32 // = 1.0f/60.0f
	MouseDoubleClickTime    float32 // 0.3
	MouseDoubleClickMaxDist float32 // 6.0

	//callbacks
	SetCursor func(c CursorType)

	//scroll
	ScrollX, ScrollY float64
}

type Key int

func (io *Io) SetDisplaySize(w, h float32) {
	io.DisplaySize.X = w
	io.DisplaySize.Y = h
}

// TODO: Refactor this if possible, it looks messy
func (io *Io) keyToString(key GuiKey) string {
	if io.IsKeyModPressed(ModShift) {
		switch key {
		case GuiKey_0:
			return ")"
		case GuiKey_1:
			return "!"
		case GuiKey_2:
			return "@"
		case GuiKey_3:
			return "#"
		case GuiKey_4:
			return "$"
		case GuiKey_5:
			return "%"
		case GuiKey_6:
			return "^"
		case GuiKey_7:
			return "&"
		case GuiKey_8:
			return "*"
		case GuiKey_9:
			return "("
		case GuiKey_A:
			return "A"
		case GuiKey_B:
			return "B"
		case GuiKey_C:
			return "C"
		case GuiKey_D:
			return "D"
		case GuiKey_E:
			return "E"
		case GuiKey_F:
			return "F"
		case GuiKey_G:
			return "G"
		case GuiKey_H:
			return "H"
		case GuiKey_I:
			return "I"
		case GuiKey_J:
			return "J"
		case GuiKey_K:
			return "K"
		case GuiKey_L:
			return "L"
		case GuiKey_M:
			return "M"
		case GuiKey_N:
			return "N"
		case GuiKey_O:
			return "O"
		case GuiKey_P:
			return "P"
		case GuiKey_Q:
			return "Q"
		case GuiKey_R:
			return "R"
		case GuiKey_S:
			return "S"
		case GuiKey_T:
			return "T"
		case GuiKey_U:
			return "U"
		case GuiKey_V:
			return "V"
		case GuiKey_W:
			return "W"
		case GuiKey_X:
			return "X"
		case GuiKey_Y:
			return "Y"
		case GuiKey_Z:
			return "Z"
		case GuiKey_Backspace:
			return "backspace"
		case GuiKey_Space:
			return " "
		case GuiKey_Enter:
			return "\n"
		case GuiKey_Apostrophe: // '
			return "\""
		case GuiKey_Comma: // ,
			return "<"
		case GuiKey_Minus: // -
			return "_"
		case GuiKey_Period: // .
			return ">"
		case GuiKey_Slash: // /
			return "?"
		case GuiKey_Semicolon: // ;
			return ":"
		case GuiKey_Equal: // =
			return "+"
		case GuiKey_LeftBracket: // [
			return "{"
		case GuiKey_Backslash: // \  // (this text inhibit multiline comment caused by backslash)
			return "|"
		case GuiKey_RightBracket: // ]
			return "}"
		case GuiKey_GraveAccent: // `
			return "~"
		default:
			return ""
		}
	} else {
		switch key {
		case GuiKey_0:
			return "0"
		case GuiKey_1:
			return "1"
		case GuiKey_2:
			return "2"
		case GuiKey_3:
			return "3"
		case GuiKey_4:
			return "4"
		case GuiKey_5:
			return "5"
		case GuiKey_6:
			return "6"
		case GuiKey_7:
			return "7"
		case GuiKey_8:
			return "8"
		case GuiKey_9:
			return "9"
		case GuiKey_A:
			return "a"
		case GuiKey_B:
			return "b"
		case GuiKey_C:
			return "c"
		case GuiKey_D:
			return "d"
		case GuiKey_E:
			return "e"
		case GuiKey_F:
			return "f"
		case GuiKey_G:
			return "g"
		case GuiKey_H:
			return "h"
		case GuiKey_I:
			return "i"
		case GuiKey_J:
			return "j"
		case GuiKey_K:
			return "k"
		case GuiKey_L:
			return "l"
		case GuiKey_M:
			return "m"
		case GuiKey_N:
			return "n"
		case GuiKey_O:
			return "o"
		case GuiKey_P:
			return "p"
		case GuiKey_Q:
			return "q"
		case GuiKey_R:
			return "r"
		case GuiKey_S:
			return "s"
		case GuiKey_T:
			return "t"
		case GuiKey_U:
			return "u"
		case GuiKey_V:
			return "v"
		case GuiKey_W:
			return "w"
		case GuiKey_X:
			return "x"
		case GuiKey_Y:
			return "y"
		case GuiKey_Z:
			return "z"
		case GuiKey_Backspace:
			return "backspace"
		case GuiKey_Space:
			return " "
		case GuiKey_Enter:
			return "\n"
		case GuiKey_Apostrophe: // '
			return "'"
		case GuiKey_Comma: // ,
			return ","
		case GuiKey_Minus: // -
			return "-"
		case GuiKey_Period: // .
			return "."
		case GuiKey_Slash: // /
			return "/"
		case GuiKey_Semicolon: // ;
			return ";"
		case GuiKey_Equal: // =
			return "="
		case GuiKey_LeftBracket: // [
			return "["
		case GuiKey_Backslash: // \  // (this text inhibit multiline comment caused by backslash)
			return "\\"
		case GuiKey_RightBracket: // ]
			return "]"
		case GuiKey_GraveAccent: // `
			return "`"
		default:
			return ""
		}
	}

}

func (io *Io) KeyCallback(key GuiKey, mods ModKey, pressed bool) {
	if key == GuiKey_None {
		return
	}
	if pressed {
		io.PressedKey = key
		io.modPressed[mods] = true
		io.KeyPressedThisFrame = true
		io.keyPressed[key] = true
	} else {
		io.modPressed[mods] = false
		io.keyPressed[key] = false
	}
}

func (io *Io) DragStarted(rect utils.Rect) bool {
	return utils.PointInRect(io.dragStarted, rect)
}

func (io *Io) MousePosCallback(x, y float32) {
	io.MousePos.X = x
	io.MousePos.Y = y

	io.IsDragging = io.MouseDown[0] ||
		io.MouseDown[1] ||
		io.MouseDown[2]

	if io.IsDragging {
		io.dragDelta = io.MousePos.Sub(io.dragStarted)
	} else {
		io.dragDelta = utils.Vec2{}
	}
}

func (io *Io) MouseBtnCallback(mouseBtn MouseKey, action Action) {
	switch action {
	case Press:
		//если кнопок на мышке больше, игнорируем нажатие
		if int(mouseBtn) < len(io.MouseDown) {
			io.MouseDown[mouseBtn] = true
			io.dragStarted = io.MousePos
			io.dragStartedMain = io.MousePos
		}
	case Release:
		if int(mouseBtn) < len(io.MouseDown) {
			io.MouseDown[mouseBtn] = false
			io.IsDragging = false
			io.dragStarted = utils.Vec2{}
			io.dragStartedMain = utils.Vec2{}
		}
	}
}

func (io *Io) IsKeyPressed(key GuiKey) bool {
	return io.keyPressed[key]
}
func (io *Io) IsKeyModPressed(key ModKey) bool {
	return io.modPressed[key]
}

const (
	MOUSE_INVALID float32 = -256000.0
)

func (io *Io) IsMousePosValid(pos *utils.Vec2) bool {
	var r utils.Vec2
	if pos != nil {
		r = *pos
	} else {
		r = io.MousePos
	}
	return r.X >= MOUSE_INVALID && r.Y >= MOUSE_INVALID
}

func NewIo() *Io {
	i := Io{
		MousePos:                utils.Vec2{0, 0},
		DisplaySize:             &utils.Vec2{0, 0},
		keyPressed:              [610]bool{},
		WantCaptureMouse:        false,
		MouseDown:               [MOUSE_BTN_COUNT]bool{},
		MouseReleased:           [5]bool{},
		MouseClicked:            [5]bool{},
		MouseDoubleClicked:      [5]bool{},
		MouseClickedTime:        [5]float32{},
		MouseDownDuration:       [5]float32{},
		MouseDownDurationPrev:   [5]float32{},
		MouseDragMaxDistanceSqr: [5]float32{},
		MouseClickedPos:         [5]utils.Vec2{},
		MouseClickedCount:       [5]uint8{},
		MouseClickedLastCount:   [5]uint8{},
		IsDragging:              false,
		dragDelta:               utils.Vec2{},
		dragStarted:             utils.Vec2{},
		dragStartedMain:         utils.Vec2{},
		MousePosPrev:            utils.Vec2{},
		MouseDelta:              utils.Vec2{},
		DeltaTime:               1 / 60,
		MouseDoubleClickTime:    0.3,
		MouseDoubleClickMaxDist: 6,
	}
	return &i
}
