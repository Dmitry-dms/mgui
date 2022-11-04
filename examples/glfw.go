package main

import (
	"github.com/Dmitry-dms/mgui"
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/go-gl/glfw/v3.3/glfw"
	//"github.com/go-gl/mathgl/mgl32"
)

type glfwHandler struct {
	GlfwWindow         *glfw.Window
	Time               float32
	MouseWindow        *glfw.Window
	MouseCursorse      [9]*glfw.Cursor
	LastValidMousePos  utils.Vec2
	InstalledCallbacks bool

	//callbacks
	// PrevUserCallbackWindowFocus GLFWwindowfocusfun
	// PrevUserCallbackCursorPos   GLFWcursorposfun
	// PrevUserCallbackCursorEnter GLFWcursorenterfun
	PrevUserCallbackMousebutton glfw.MouseButtonCallback
	PrevUserCallbackScroll      glfw.ScrollCallback
	PrevUserCallbackKey         glfw.KeyCallback
	// 	PrevUserCallbackChar        GLFWcharfun
	// 	PrevUserCallbackMonitor     GLFWmonitorfun
	SetCursor func(c ui.CursorType)
}

func NewData() *glfwHandler {
	d := glfwHandler{}
	f := func(c ui.CursorType) {
		defC := glfw.CreateStandardCursor(Cursor(c))
		d.GlfwWindow.SetCursor(defC)
	}
	d.SetCursor = f
	return &d
}

func GlfwKeyToGuiKey(key glfw.Key) ui.GuiKey {
	switch key {
	case glfw.KeyTab:
		return ui.GuiKey_Tab
	case glfw.KeyLeft:
		return ui.GuiKey_LeftArrow
	case glfw.KeyRight:
		return ui.GuiKey_RightArrow
	case glfw.KeyUp:
		return ui.GuiKey_UpArrow
	case glfw.KeyDown:
		return ui.GuiKey_DownArrow
	case glfw.KeyPageUp:
		return ui.GuiKey_PageUp
	case glfw.KeyPageDown:
		return ui.GuiKey_PageDown
	case glfw.KeyHome:
		return ui.GuiKey_Home
	case glfw.KeyEnd:
		return ui.GuiKey_End
	case glfw.KeyInsert:
		return ui.GuiKey_Insert
	case glfw.KeyDelete:
		return ui.GuiKey_Delete
	case glfw.KeyBackspace:
		return ui.GuiKey_Backspace
	case glfw.KeySpace:
		return ui.GuiKey_Space
	case glfw.KeyEnter:
		return ui.GuiKey_Enter
	case glfw.KeyEscape:
		return ui.GuiKey_Escape
	case glfw.KeyApostrophe:
		return ui.GuiKey_Apostrophe
	case glfw.KeyComma:
		return ui.GuiKey_Comma
	case glfw.KeyMinus:
		return ui.GuiKey_Minus
	case glfw.KeyPeriod:
		return ui.GuiKey_Period
	case glfw.KeySlash:
		return ui.GuiKey_Slash
	case glfw.KeySemicolon:
		return ui.GuiKey_Semicolon
	case glfw.KeyEqual:
		return ui.GuiKey_Equal
	case glfw.KeyLeftBracket:
		return ui.GuiKey_LeftBracket
	case glfw.KeyBackslash:
		return ui.GuiKey_Backslash
	case glfw.KeyRightBracket:
		return ui.GuiKey_RightBracket
	case glfw.KeyGraveAccent:
		return ui.GuiKey_GraveAccent
	case glfw.KeyCapsLock:
		return ui.GuiKey_CapsLock
	case glfw.KeyScrollLock:
		return ui.GuiKey_ScrollLock
	case glfw.KeyNumLock:
		return ui.GuiKey_NumLock
	case glfw.KeyPrintScreen:
		return ui.GuiKey_PrintScreen
	case glfw.KeyPause:
		return ui.GuiKey_Pause
	case glfw.KeyKP0:
		return ui.GuiKey_Keypad0
	case glfw.KeyKP1:
		return ui.GuiKey_Keypad1
	case glfw.KeyKP2:
		return ui.GuiKey_Keypad2
	case glfw.KeyKP3:
		return ui.GuiKey_Keypad3
	case glfw.KeyKP4:
		return ui.GuiKey_Keypad4
	case glfw.KeyKP5:
		return ui.GuiKey_Keypad5
	case glfw.KeyKP6:
		return ui.GuiKey_Keypad6
	case glfw.KeyKP7:
		return ui.GuiKey_Keypad7
	case glfw.KeyKP8:
		return ui.GuiKey_Keypad8
	case glfw.KeyKP9:
		return ui.GuiKey_Keypad9
	case glfw.KeyKPDecimal:
		return ui.GuiKey_KeypadDecimal
	case glfw.KeyKPDivide:
		return ui.GuiKey_KeypadDivide
	case glfw.KeyKPMultiply:
		return ui.GuiKey_KeypadMultiply
	case glfw.KeyKPSubtract:
		return ui.GuiKey_KeypadSubtract
	case glfw.KeyKPAdd:
		return ui.GuiKey_KeypadAdd
	case glfw.KeyKPEnter:
		return ui.GuiKey_KeypadEnter
	case glfw.KeyKPEqual:
		return ui.GuiKey_KeypadEqual
	case glfw.KeyLeftShift:
		return ui.GuiKey_LeftShift
	case glfw.KeyLeftControl:
		return ui.GuiKey_LeftCtrl
	case glfw.KeyLeftAlt:
		return ui.GuiKey_LeftAlt
	case glfw.KeyLeftSuper:
		return ui.GuiKey_LeftSuper
	case glfw.KeyRightShift:
		return ui.GuiKey_RightShift
	case glfw.KeyRightControl:
		return ui.GuiKey_RightCtrl
	case glfw.KeyRightAlt:
		return ui.GuiKey_RightAlt
	case glfw.KeyRightSuper:
		return ui.GuiKey_RightSuper
	case glfw.KeyMenu:
		return ui.GuiKey_Menu
	case glfw.Key0:
		return ui.GuiKey_0
	case glfw.Key1:
		return ui.GuiKey_1
	case glfw.Key2:
		return ui.GuiKey_2
	case glfw.Key3:
		return ui.GuiKey_3
	case glfw.Key4:
		return ui.GuiKey_4
	case glfw.Key5:
		return ui.GuiKey_5
	case glfw.Key6:
		return ui.GuiKey_6
	case glfw.Key7:
		return ui.GuiKey_7
	case glfw.Key8:
		return ui.GuiKey_8
	case glfw.Key9:
		return ui.GuiKey_9
	case glfw.KeyA:
		return ui.GuiKey_A
	case glfw.KeyB:
		return ui.GuiKey_B
	case glfw.KeyC:
		return ui.GuiKey_C
	case glfw.KeyD:
		return ui.GuiKey_D
	case glfw.KeyE:
		return ui.GuiKey_E
	case glfw.KeyF:
		return ui.GuiKey_F
	case glfw.KeyG:
		return ui.GuiKey_G
	case glfw.KeyH:
		return ui.GuiKey_H
	case glfw.KeyI:
		return ui.GuiKey_I
	case glfw.KeyJ:
		return ui.GuiKey_J
	case glfw.KeyK:
		return ui.GuiKey_K
	case glfw.KeyL:
		return ui.GuiKey_L
	case glfw.KeyM:
		return ui.GuiKey_M
	case glfw.KeyN:
		return ui.GuiKey_N
	case glfw.KeyO:
		return ui.GuiKey_O
	case glfw.KeyP:
		return ui.GuiKey_P
	case glfw.KeyQ:
		return ui.GuiKey_Q
	case glfw.KeyR:
		return ui.GuiKey_R
	case glfw.KeyS:
		return ui.GuiKey_S
	case glfw.KeyT:
		return ui.GuiKey_T
	case glfw.KeyU:
		return ui.GuiKey_U
	case glfw.KeyV:
		return ui.GuiKey_V
	case glfw.KeyW:
		return ui.GuiKey_W
	case glfw.KeyX:
		return ui.GuiKey_X
	case glfw.KeyY:
		return ui.GuiKey_Y
	case glfw.KeyZ:
		return ui.GuiKey_Z
	case glfw.KeyF1:
		return ui.GuiKey_F1
	case glfw.KeyF2:
		return ui.GuiKey_F2
	case glfw.KeyF3:
		return ui.GuiKey_F3
	case glfw.KeyF4:
		return ui.GuiKey_F4
	case glfw.KeyF5:
		return ui.GuiKey_F5
	case glfw.KeyF6:
		return ui.GuiKey_F6
	case glfw.KeyF7:
		return ui.GuiKey_F7
	case glfw.KeyF8:
		return ui.GuiKey_F8
	case glfw.KeyF9:
		return ui.GuiKey_F9
	case glfw.KeyF10:
		return ui.GuiKey_F10
	case glfw.KeyF11:
		return ui.GuiKey_F11
	case glfw.KeyF12:
		return ui.GuiKey_F12
	default:
		return ui.GuiKey_None
	}
}

func GlfwMouseKey(btn glfw.MouseButton) ui.MouseKey {
	switch btn {
	case glfw.MouseButtonLeft:
		return ui.MouseBtnLeft
	case glfw.MouseButtonRight:
		return ui.MouseBtnRight
	case glfw.MouseButtonMiddle:
		return ui.MouseBtnMiddle
	default:
		return ui.MouseBtnUnknown
	}
}

func GlfwModKey(mod glfw.ModifierKey) ui.ModKey {
	switch mod {
	case glfw.ModAlt:
		return ui.ModAlt
	case glfw.ModControl:
		return ui.ModCtrl
	case glfw.ModShift:
		return ui.ModShift
	default:
		return ui.UnknownMod
	}
}

func GlfwAction(action glfw.Action) ui.Action {
	switch action {
	case glfw.Press:
		return ui.Press
	case glfw.Release:
		return ui.Release
	case glfw.Repeat:
		return ui.Repeat
	default:
		return ui.UnknownAction
	}
}

func Cursor(cur ui.CursorType) glfw.StandardCursor {
	switch cur {
	case ui.VResizeCursor:
		return glfw.VResizeCursor
	case ui.ArrowCursor:
		return glfw.ArrowCursor
	case ui.HResizeCursor:
		return glfw.HResizeCursor
	case ui.EditCursor:
		return glfw.CrosshairCursor
	default:
		return glfw.ArrowCursor
	}
}

type PlatformHandler interface {
	IsKeyPressed(key ui.Key)
	UpdateInputs()
	SetCursor(c ui.CursorType)
}
