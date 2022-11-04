package ui

type GuiKey int

const (
	GuiKey_NamedKey_BEGIN = 512
	// Keyboard
	GuiKey_None GuiKey = 0
	GuiKey_Tab         = GuiKey_NamedKey_BEGIN + iota // == ImGuiKey_NamedKey_BEGIN
	GuiKey_LeftArrow
	GuiKey_RightArrow
	GuiKey_UpArrow
	GuiKey_DownArrow
	GuiKey_PageUp
	GuiKey_PageDown
	GuiKey_Home
	GuiKey_End
	GuiKey_Insert
	GuiKey_Delete
	GuiKey_Backspace
	GuiKey_Space
	GuiKey_Enter
	GuiKey_Escape
	GuiKey_LeftCtrl
	GuiKey_LeftShift
	GuiKey_LeftAlt
	GuiKey_LeftSuper
	GuiKey_RightCtrl
	GuiKey_RightShift
	GuiKey_RightAlt
	GuiKey_RightSuper
	GuiKey_Menu
	GuiKey_0
	GuiKey_1
	GuiKey_2
	GuiKey_3
	GuiKey_4
	GuiKey_5
	GuiKey_6
	GuiKey_7
	GuiKey_8
	GuiKey_9
	GuiKey_A
	GuiKey_B
	GuiKey_C
	GuiKey_D
	GuiKey_E
	GuiKey_F
	GuiKey_G
	GuiKey_H
	GuiKey_I
	GuiKey_J
	GuiKey_K
	GuiKey_L
	GuiKey_M
	GuiKey_N
	GuiKey_O
	GuiKey_P
	GuiKey_Q
	GuiKey_R
	GuiKey_S
	GuiKey_T
	GuiKey_U
	GuiKey_V
	GuiKey_W
	GuiKey_X
	GuiKey_Y
	GuiKey_Z

	GuiKey_F1
	GuiKey_F2
	GuiKey_F3
	GuiKey_F4
	GuiKey_F5
	GuiKey_F6
	GuiKey_F7
	GuiKey_F8
	GuiKey_F9
	GuiKey_F10
	GuiKey_F11
	GuiKey_F12
	GuiKey_Apostrophe   // '
	GuiKey_Comma        // ,
	GuiKey_Minus        // -
	GuiKey_Period       // .
	GuiKey_Slash        // /
	GuiKey_Semicolon    // ;
	GuiKey_Equal        // =
	GuiKey_LeftBracket  // [
	GuiKey_Backslash    // \ (this text inhibit multiline comment caused by backslash)
	GuiKey_RightBracket // ]
	GuiKey_GraveAccent  // `
	GuiKey_CapsLock
	GuiKey_ScrollLock
	GuiKey_NumLock
	GuiKey_PrintScreen
	GuiKey_Pause
	GuiKey_Keypad0
	GuiKey_Keypad1
	GuiKey_Keypad2
	GuiKey_Keypad3
	GuiKey_Keypad4
	GuiKey_Keypad5
	GuiKey_Keypad6
	GuiKey_Keypad7
	GuiKey_Keypad8
	GuiKey_Keypad9
	GuiKey_KeypadDecimal
	GuiKey_KeypadDivide
	GuiKey_KeypadMultiply
	GuiKey_KeypadSubtract
	GuiKey_KeypadAdd
	GuiKey_KeypadEnter
	GuiKey_KeypadEqual
)

type ModKey int

func IsCommandKey(key GuiKey) bool {
	if key >= 514 && key <= 537 || key >= 574 && key <= 601 {
		return true
	}
	return false
}

const (
	ModShift ModKey = iota
	ModCtrl
	ModAlt
	ModSuper
	ModCapsLock
	ModNumLock
	ModControl
	UnknownMod
)

type MouseKey int

const (
	MouseBtnLeft MouseKey = iota
	MouseBtnRight
	MouseBtnMiddle
	MouseBtnUnknown
)

type Action int

const (
	Press Action = iota
	Release
	Repeat
	UnknownAction
)

type CursorType int

const (
	ArrowCursor CursorType = iota
	HResizeCursor
	VResizeCursor
	EditCursor
	UnknownCursor
)
