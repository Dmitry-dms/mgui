package main

import (
	ui "github.com/Dmitry-dms/mgui"
	"github.com/Dmitry-dms/mgui/examples/wasm/gltypes"
	"syscall/js"
)

var (
	width   int
	height  int
	glCtx   js.Value
	glTypes gltypes.GLTypes
	done    chan struct{}
	doc     js.Value
)

func setCursor(c ui.CursorType) {
	//cur := Cursor(c)
	//doc.Get("style").Set("cursor", cur)
}
func JSMouseKey(btn int) ui.MouseKey {
	switch btn {
	case 0:
		return ui.MouseBtnLeft
	case 2:
		return ui.MouseBtnRight
	case 1:
		return ui.MouseBtnMiddle
	default:
		return ui.MouseBtnUnknown
	}
}

func JSMouseAction(action int) ui.Action {
	switch action {
	case 0:
		return ui.Press
	case 1:
		return ui.Release
	case 2:
		return ui.Repeat
	default:
		return ui.UnknownAction
	}
}

func main() {
	// Init Canvas stuff
	doc = js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "gocanvas")
	width = doc.Get("body").Get("clientWidth").Int()
	height = doc.Get("body").Get("clientHeight").Int()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)

	//glCtx, _ := NewContext(canvasEl)

	glCtx = canvasEl.Call("getContext", "webgl2")
	if glCtx.IsUndefined() {
		glCtx = canvasEl.Call("getContext", "experimental-webgl")
	}
	// once again
	if glCtx.IsUndefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

	done = make(chan struct{}, 0)

	// Get some WebGL bindings
	glTypes.New(glCtx)

	js.Global().Get("console").Call("log", js.Global().Get("Cur"))

	js.Global().Get("console").Call("log", "VERSION: ", glCtx.Call("getParameter", glCtx.Get("SHADING_LANGUAGE_VERSION")))
	renderer := NewWebGLRenderer(glCtx)

	ui.AddRenderer(renderer)

	fncs := initListeners(doc)
	defer fncs()
	ui.SetChangeCursorFunc(setCursor)
	var open = true

	//glCtx.Viewport(0, 0, width, height)
	glCtx.Call("viewport", 0, 0, width, height) // Viewport size
	js.Global().Get("console").Call("log", "STARTUP!!!!")
	renderer.renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//glCtx.ClearColor(0.5, 0.5, 0.5, 0.9)
		//glCtx.ClearDepth(1.0)
		glCtx.Call("clearColor", 0.5, 0.5, 0.5, 0.9) // Color the screen is cleared to
		glCtx.Call("clearDepth", 1.0)                // Z value that is set to the Depth buffer every frame

		glCtx.Call("depthFunc", glTypes.LEqual)

		ui.NewFrame([2]float32{float32(width), float32(height)})
		//ui.GlobalWidgetSpace("sfdfds", 0, 0, 300, 500, ui.Default, func() {
		//	if ui.Button("btnuju") {
		//		js.Global().Get("console").Call("log", "PRESSED!!!!!")
		//	}
		//})
		ui.BeginWindow("Th", &open)

		//ui.GlobalWidgetSpace("sfdfds", 300, 300, 500, 500, ui.Default, func() {
		if ui.Button("btn") {
			js.Global().Get("console").Call("log", "PRESSED!!!!!")
		}
		//})

		ui.EndWindow()

		ui.EndFrame([2]float32{float32(width), float32(height)})

		return nil
	})

	js.Global().Call("requestAnimationFrame", renderer.renderFrame)

	//js.Global().Set("exit", js.FuncOf(exit))

	<-done
	renderer.renderFrame.Release()
}

func initListeners(doc js.Value) func() {
	resize := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//e := args[0]
		//width = int(e.Get("clientWidth").Float())
		//height = int(e.Get("clientHeight").Float())
		//ui.SetDisplaySize(float32(width), float32(height))
		js.Global().Get("console").Call("log", "RESIZE!!!!")
		return nil
	})

	doc.Call("addEventListener", "resize", resize)

	mousePos := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x := int(e.Get("clientX").Float())
		y := int(e.Get("clientY").Float())
		ui.GetIo().MousePosCallback(float32(x), float32(y))
		return nil
	})

	doc.Call("addEventListener", "mousemove", mousePos)

	mouseDownEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		//js.Global().Get("console").Call("log", e)
		ui.GetIo().MouseBtnCallback(JSMouseKey(e.Get("button").Int()), JSMouseAction(0))
		return nil
	})

	doc.Call("addEventListener", "mousedown", mouseDownEvt)

	mouseUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		//js.Global().Get("console").Call("log", e)
		ui.GetIo().MouseBtnCallback(JSMouseKey(e.Get("button").Int()), JSMouseAction(1))
		return nil
	})

	doc.Call("addEventListener", "mouseup", mouseUpEvt)

	keyPressEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		e.Call("preventDefault")
		//key := e.Get("key")
		//js.Global().Get("console").Call("log", e)
		onKey(e, true)
		return nil
	})
	doc.Call("addEventListener", "keypress", keyPressEvt)

	keyReleaseEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		e.Call("preventDefault")
		//key := e.Get("key")
		onKey(e, false)
		return nil
	})
	doc.Call("addEventListener", "keyup", keyReleaseEvt)

	return func() {
		keyPressEvt.Release()
		mouseUpEvt.Release()
		mouseDownEvt.Release()
		mousePos.Release()
		resize.Release()
		keyReleaseEvt.Release()
	}
}

func JsKeyToGuiKey(key string) ui.GuiKey {
	switch key {
	case "Tab":
		return ui.GuiKey_Tab
	case "ArrowLeft":
		return ui.GuiKey_LeftArrow
	case "ArrowRight":
		return ui.GuiKey_RightArrow
	case "ArrowUp":
		return ui.GuiKey_UpArrow
	case "ArrowDown":
		return ui.GuiKey_DownArrow
	case "PageUp":
		return ui.GuiKey_PageUp
	case "PageDown":
		return ui.GuiKey_PageDown
	case "Home":
		return ui.GuiKey_Home
	case "End":
		return ui.GuiKey_End
	case "Insert":
		return ui.GuiKey_Insert
	case "Delete":
		return ui.GuiKey_Delete
	case "Backspace":
		return ui.GuiKey_Backspace
	case "Space":
		return ui.GuiKey_Space
	case "Enter":
		return ui.GuiKey_Enter
	case "Escape":
		return ui.GuiKey_Escape
	//case glfw.KeyApostrophe:
	//	return ui.GuiKey_Apostrophe
	case "Comma":
		return ui.GuiKey_Comma
	case "Minus":
		return ui.GuiKey_Minus
	case "Period":
		return ui.GuiKey_Period
	case "Slash":
		return ui.GuiKey_Slash
	case "Semicolon":
		return ui.GuiKey_Semicolon
	case "Equal":
		return ui.GuiKey_Equal
	case "BracketLeft":
		return ui.GuiKey_LeftBracket
	case "Backslash":
		return ui.GuiKey_Backslash
	case "BracketRight":
		return ui.GuiKey_RightBracket
	case "Backquote":
		return ui.GuiKey_GraveAccent
	case "CapsLock":
		return ui.GuiKey_CapsLock
	case "ScrollLock":
		return ui.GuiKey_ScrollLock
	case "NumLock":
		return ui.GuiKey_NumLock
	case "PrintScreen":
		return ui.GuiKey_PrintScreen
	case "Pause":
		return ui.GuiKey_Pause
	case "Numpad0":
		return ui.GuiKey_Keypad0
	case "Numpad1":
		return ui.GuiKey_Keypad1
	case "Numpad2":
		return ui.GuiKey_Keypad2
	case "Numpad3":
		return ui.GuiKey_Keypad3
	case "Numpad4":
		return ui.GuiKey_Keypad4
	case "Numpad5":
		return ui.GuiKey_Keypad5
	case "Numpad6":
		return ui.GuiKey_Keypad6
	case "Numpad7":
		return ui.GuiKey_Keypad7
	case "Numpad8":
		return ui.GuiKey_Keypad8
	case "Numpad9":
		return ui.GuiKey_Keypad9
	case "NumpadDecimal":
		return ui.GuiKey_KeypadDecimal
	case "NumpadDivide":
		return ui.GuiKey_KeypadDivide
	case "NumpadMultiply":
		return ui.GuiKey_KeypadMultiply
	case "NumpadSubtract":
		return ui.GuiKey_KeypadSubtract
	case "NumpadAdd":
		return ui.GuiKey_KeypadAdd
	case "NumpadEnter":
		return ui.GuiKey_KeypadEnter
	//case "glfw.KeyKPEqual":
	//	return ui.GuiKey_KeypadEqual
	case "ShiftLeft":
		return ui.GuiKey_LeftShift
	case "ControlLeft":
		return ui.GuiKey_LeftCtrl
	case "AltLeft":
		return ui.GuiKey_LeftAlt
	case "MetaLeft":
		return ui.GuiKey_LeftSuper
	case "ShiftRight":
		return ui.GuiKey_RightShift
	case "ControlRight":
		return ui.GuiKey_RightCtrl
	case "AltRight":
		return ui.GuiKey_RightAlt
	case "MetaRight":
		return ui.GuiKey_RightSuper
	//case glfw.KeyMenu:
	//	return ui.GuiKey_Menu
	case "Digit0":
		return ui.GuiKey_0
	case "Digit1":
		return ui.GuiKey_1
	case "Digit2":
		return ui.GuiKey_2
	case "Digit3":
		return ui.GuiKey_3
	case "Digit4":
		return ui.GuiKey_4
	case "Digit5":
		return ui.GuiKey_5
	case "Digit6":
		return ui.GuiKey_6
	case "Digit7":
		return ui.GuiKey_7
	case "Digit8":
		return ui.GuiKey_8
	case "Digit9":
		return ui.GuiKey_9
	case "KeyA":
		return ui.GuiKey_A
	case "KeyB":
		return ui.GuiKey_B
	case "KeyC":
		return ui.GuiKey_C
	case "KeyD":
		return ui.GuiKey_D
	case "KeyE":
		return ui.GuiKey_E
	case "KeyF":
		return ui.GuiKey_F
	case "KeyG":
		return ui.GuiKey_G
	case "KeyH":
		return ui.GuiKey_H
	case "KeyI":
		return ui.GuiKey_I
	case "KeyJ":
		return ui.GuiKey_J
	case "KeyK":
		return ui.GuiKey_K
	case "KeyL":
		return ui.GuiKey_L
	case "KeyM":
		return ui.GuiKey_M
	case "KeyN":
		return ui.GuiKey_N
	case "KeyO":
		return ui.GuiKey_O
	case "KeyP":
		return ui.GuiKey_P
	case "KeyQ":
		return ui.GuiKey_Q
	case "KeyR":
		return ui.GuiKey_R
	case "KeyS":
		return ui.GuiKey_S
	case "KeyT":
		return ui.GuiKey_T
	case "KeyU":
		return ui.GuiKey_U
	case "KeyV":
		return ui.GuiKey_V
	case "KeyW":
		return ui.GuiKey_W
	case "KeyX":
		return ui.GuiKey_X
	case "KeyY":
		return ui.GuiKey_Y
	case "KeyZ":
		return ui.GuiKey_Z
	case "F1":
		return ui.GuiKey_F1
	case "F2":
		return ui.GuiKey_F2
	case "F3":
		return ui.GuiKey_F3
	case "F4":
		return ui.GuiKey_F4
	case "F5":
		return ui.GuiKey_F5
	case "F6":
		return ui.GuiKey_F6
	case "F7":
		return ui.GuiKey_F7
	case "F8":
		return ui.GuiKey_F8
	case "F9":
		return ui.GuiKey_F9
	case "F10":
		return ui.GuiKey_F10
	case "F11":
		return ui.GuiKey_F11
	case "F12":
		return ui.GuiKey_F12
	default:
		return ui.GuiKey_None
	}
}

func JSModKey(alt, ctr, shift bool) ui.ModKey {
	if alt {
		return ui.ModAlt
	}
	if ctr {
		return ui.ModCtrl
	}
	if shift {
		return ui.ModShift
	}
	return ui.UnknownMod
}
func onKey(key js.Value, pressed bool) {
	code := key.Get("code").String()
	alt := key.Get("altKey").Bool()
	shift := key.Get("shiftKey").Bool()
	ctrl := key.Get("ctrlKey").Bool()

	if pressed {
		ui.GetIo().KeyCallback(JsKeyToGuiKey(code), JSModKey(alt, ctrl, shift), true)
	} else {
		ui.GetIo().KeyCallback(JsKeyToGuiKey(code), JSModKey(alt, ctrl, shift), false)
	}
}

func Cursor(cur ui.CursorType) string {
	switch cur {
	case ui.VResizeCursor:
		return "w-resize"
	case ui.ArrowCursor:
		return "default"
	case ui.HResizeCursor:
		return "s-resize"
	case ui.EditCursor:
		return "crosshair"
	default:
		return "default"
	}
}

func exit(this js.Value, inputs []js.Value) interface{} {
	done <- struct{}{}
	js.Global().Get("console").Call("log", "EXIT!!!!!")
	return nil
}
