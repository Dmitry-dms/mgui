package main

import (
	"fmt"
	ui "github.com/Dmitry-dms/mgui"
	"github.com/Dmitry-dms/mgui/fonts"
	"github.com/Dmitry-dms/mgui/sprite_packer"
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/png"
	"os"
	"runtime"
	"sort"
)

func init() {
	runtime.LockOSThread()
}

var uiCtx *ui.UiContext
var window *glfw.Window
var Width, Height int = 1280, 720
var steps int = 1

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	glfw.DefaultWindowHints()
	// glfw.WindowHint(glfw.OpenGLDebugContext, 1)

	window, err = glfw.CreateWindow(Width, Height, "example", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	size := func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		fmt.Println(width, height)
		Width = width
		Height = height
		// uiCtx.Io().SetDisplaySize(float32(width), float32(height))
	}

	window.SetSizeCallback(size)
	window.SetKeyCallback(onKey)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetScrollCallback(scrollCallback)

	// gogl.InitGLdebug()

	// font := fonts.NewFont("assets/fonts/rany.otf", 60)
	// font := fonts.NewFont("assets/fonts/mono.ttf", 60)
	// font := fonts.NewFont("assets/fonts/Roboto.ttf", 60)

	// ctx.Io.DefaultFont, _ = ui2.LoadFontFromFile("C:/Windows/Fonts/times.ttf", 40)

	// font := fonts.NewFont("C:/Windows/Fonts/times.ttf", 40, true)

	// batch := fonts.NewTextBatch(font)
	// batch.Init()
	gl.Init()

	uiCtx = ui.UiCtx

	front := NewGlRenderer()
	uiCtx.Initialize(front)
	// uiCtx = ui.NewContext(front, cam)
	uiCtx.Io().SetCursor = setCursor

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//sheet := sprite_packer.NewSpriteSheet(512, "mgui")
	////ui.UiCtx.UploadFont("C:/Windows/Fonts/times.ttf", 14)
	fontName := "C:/Windows/Fonts/arial.ttf"
	fontName2 := "C:/Windows/Fonts/times.ttf"
	f, _ := ui.UiCtx.UploadFont(fontName2, 24, 157.0, 32, 256)
	f2, _ := ui.UiCtx.UploadFont(fontName, 18, 157.0, 32, 256)
	//ConvertFontToAtlas(f, sheet, data)
	//ConvertFontToAtlas(f2, sheet, data2)

	////ui.UiCtx.UploadFont("assets/fonts/rany.otf", 14)
	////ui.UiCtx.UploadFont("assets/fonts/sans.ttf", 18)
	////ui.UiCtx.UploadFont("assets/fonts/mono.ttf", 14)
	//mario := openImage("assets/images/mario.png")
	//goomba := openImage("assets/images/goomba.png")
	//windowCustom := openImage("assets/images/window_wooden.png")
	//arrow := openImage("assets/images/arrow-down-filled.png")
	//ms := sheet.AddSprite("sprites", "mario", mario)
	//gs := sheet.AddSprite("sprites", "goomba", goomba)
	//ard := sheet.AddSprite("sprites", "arrow-down", arrow)
	//sheet.AddSprite("sprites", "window", windowCustom)

	//t2 := gogl.UploadRGBATextureFromMemory(sheet.Image())
	//f2.TextureId = t2.GetId()
	//CreateImage("debug.png", sheet.Image())

	sheet, err := sprite_packer.GetSpriteSheetFromFile("atlas.json", "debug.png")
	if err != nil {
		panic(err)
	}
	t2 := UploadRGBATextureFromMemory(sheet.Image())
	f.TextureId = t2.TextureId
	f2.TextureId = t2.TextureId

	//tex.TexCoords = ms.TextCoords
	//tex2.TextureId = t2.TextureId
	//tex.TextureId = t2.TextureId
	//tex2.TexCoords = gs.TextCoords
	//arrowDown.TextureId = t2.TextureId
	//arrowDown.TexCoords = ard.TextCoords
	//custWnd.TextureId = t2.TextureId
	//custWnd.TexCoords = ws.TextCoords

	s, ok := sheet.GetGroup("sprites")
	if !ok {
		panic("Could not find sprites group")
	}
	for _, info := range s.Contents {
		if info.Id == "mario" {
			tex.TexCoords = info.TextCoords
			tex2.TextureId = t2.TextureId
		}
		if info.Id == "goomba" {
			tex.TextureId = t2.TextureId
			tex2.TexCoords = info.TextCoords
		}
		if info.Id == "arrow-down" {
			arrowDown.TextureId = t2.TextureId
			arrowDown.TexCoords = info.TextCoords
		}
		if info.Id == "window" {
			custWnd.TextureId = t2.TextureId
			custWnd.TexCoords = info.TextCoords
		}
	}
	//
	rr, _ := sheet.GetGroup(f2.Filepath)

	for _, info := range rr.Contents {
		if info != nil {
			ll := []rune(info.Id)
			char := f2.GetCharacter(ll[0])
			char.TexCoords = [2]utils.Vec2{{info.TextCoords[0], info.TextCoords[1]},
				{info.TextCoords[2], info.TextCoords[3]}}
		}
	}

	beginTime := float32(glfw.GetTime())
	var endTime float32
	var dt float32
	dt = dt

	var time float32 = 0

	//err = sheet.SaveSpriteSheetInfo("atlas.json")
	//if err != nil {
	//	panic(err)
	//}
	//tex, _ = tex.Init("assets/images/mario.png")
	//tex2, _ = tex2.Init("assets/images/goomba.png")

	// fb, err := NewFramebuffer(200, 200)
	// if err != nil {
	// 	panic(err)
	// }
	// var p bool
	// gl.Enable(gl.SCISSOR_TEST)
	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

		//if uiCtx.Io().KeyPressedThisFrame {
		//fmt.Println(uiCtx.Io().PressedKey)
		//}

		uiCtx.NewFrame([2]float32{float32(Width), float32(Height)})

		firstWindow()
		//customWindow()

		if uiCtx.Io().IsKeyPressed(ui.GuiKey_Space) {
			//if uiCtx.SelectableText != nil {
			//	fmt.Println(uiCtx.SelectableText.WidgetId())
			//} else {
			//	fmt.Println("nil")
			//}
			//fmt.Println(uiCtx.SelectedText)
			fmt.Println(uiCtx.FocusedTextInput)
			//fmt.Println(string(uiCtx.SelectedTextStart.Chars[uiCtx.SelectedTextStart.StartInd].Char.Rune),
			//	string(uiCtx.SelectedTextEnd.Chars[uiCtx.SelectedTextStart.EndInd].Char.Rune))
			//fmt.Println(uiCtx.SelectedTextStart.WidgetId(), uiCtx.SelectedTextEnd.WidgetId())
			//fmt.Println()
			//for _, text := range uiCtx.SelectedTexts {
			//	fmt.Print(text.WidgetId() + " ")
		}
		//}

		//secondWindow()

		// fb.Bind()
		uiCtx.EndFrame([2]float32{float32(Width), float32(Height)})
		// fb.Unbind()
		// rend.NewFrame()

		window.SwapBuffers()

		endTime = float32(glfw.GetTime())
		dt = endTime - beginTime
		beginTime = endTime

		// fmt.Println(time / float32(steps))

		time += dt
		steps++
	}
}

var ish bool = false
var tW float32 = 400
var message = "hello \nworld"
var msg1 = "hello"

func customWindow() {
	uiCtx.BeginCustomWindow("cstm wnd", 500, 500, 650, 650,
		100, 160, 400, 480,
		custWnd.TextureId, custWnd.TexCoords, func() {
			uiCtx.Image("#im4kjdg464tht", 100, 100, tex.TextureId, tex.TexCoords)
			uiCtx.Text("text-ttp-4", "Обычная", ui.Selectable)
		})
}

func firstWindow() {
	uiCtx.BeginWindow("first wnd")
	//uiCtx.Selection("sel-1", &selection, sle, arrowDown)
	//uiCtx.Selection("sel-1", &selection, sle, arrowDown)
	//uiCtx.Text("text-ttp-2", "Обычная картинка \nи это то-же 1", ui.Selectable)
	//uiCtx.Text("text-ttp-3", "Обычная картинка и \nэто то-же 2", ui.Editable)
	uiCtx.Image("#im4kjdg464tht", 100, 100, tex.TextureId, tex.TexCoords)
	uiCtx.Text("text-ttp-4", "Обычная", ui.Selectable)
	//uiCtx.Text("tlorem", "Lorem Ipsum - это текст-\"рыба\", часто \nиспользуемый в печати и вэб-дизайне. Lorem Ipsum является \nстандартной \"рыбой\" для текстов на \nлатинице с начала XVI века.", ui.Selectable)
	//uiCtx.TextFitted("text-ttp-1", tW, "Lorem Ipsum - это текст-\"рыба\", часто используемый в печати и вэб-дизайне. Lorem Ipsum является стандартной \"рыбой\" для текстов на латинице с начала XVI века.")
	uiCtx.TextFitted("text-ttваы-1", tW, "Съешь ещё этих мягких французских булочек")
	uiCtx.Row("slider-row", func() {
		uiCtx.Slider("slds", &tW, 100, 1200)
		uiCtx.Text("sl-tex", fmt.Sprint(tW), ui.DefaultTextFlag)
	})
	uiCtx.MultiLineTextInput("inputr23", &message)
	uiCtx.Text("text-ttp-43", "Обычная картинка и это то-же 123", ui.Selectable)
	uiCtx.TextInput("tinp-121", 300, 50, &msg1)
	//uiCtx.Bezier()
	//uiCtx.Line(200)
	//uiCtx.Line(400)

	//if uiCtx.ButtonT("Нажать", "Press") {
	//	//	ish = !ish
	//	//
	//}
	//uiCtx.ContextMenu("ASsfdffb", func() {
	//	uiCtx.Text("#t3fdj", "Опция 1", 14)
	//	uiCtx.Text("#t3аваfdj", "Опция 2", 14)
	//	uiCtx.Text("#t3ававаfdj", "Опция 3", 14)
	//})
	////if ish {
	//uiCtx.Text("#er", "АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ", 14)
	//uiCtx.Text("#er", "АБВГДЕЖЗИЙКЛАМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ", 14)
	//uiCtx.Row("roe23", func() {
	//	uiCtx.Text("#eывr", "the quick brown fox", 14)
	//	uiCtx.Image("#im4kjdg464tht", 100, 100, tex)
	//})
	//uiCtx.Text("#eыfdвr", "Съешь ещё этих мягких", 14)

	//uiCtx.Text("#eывr", "A", 14)
	//uiCtx.Text("#eывr", "AVWAV", 14)

	//uiCtx.Text("#eы324fdвr", "1+4=5 (323_32) [A-Za-z]*^_2", 14)
	////uiCtx.Text("#fgfgd", "hello world! dfdgdfgfd 434554654 gf ", 14)
	//////}
	////
	//uiCtx.TreeNode("tree1", "Configuration", func() {
	//	uiCtx.Text("text-ttp-1", "Обычная картинка, которая  ничего не делает", ui.DefaultTextFlag)
	//	uiCtx.Text("#t3j", "hello world!", ui.DefaultTextFlag)
	//	uiCtx.TreeNode("tree1yuy2", "Настройки", func() {
	//		uiCtx.Text("texiyt-ttp-1", "Обычная картинка, которая  ничего не делает", ui.DefaultTextFlag)
	//		uiCtx.Text("#tiy3j", "hello world!", ui.DefaultTextFlag)
	//	})
	//})
	//////
	//uiCtx.VSpace("#vs1fdgdf")
	////
	//uiCtx.Row("row 13214", func() {
	//	uiCtx.Image("#im4kjdg464", tex)
	//	uiCtx.Column("col fdfd", func() {
	//		uiCtx.Image("#im76", tex2)
	//		uiCtx.Image("#im4", tex)
	//	})
	//
	//	uiCtx.Column("col fdfdвава", func() {
	//		uiCtx.Button("ASsfdffb")
	//		uiCtx.Button("ASsfdffbbb")
	//		uiCtx.Slider("slider-1", &slCounter, 0, 255)
	//	})
	//
	//	uiCtx.Image("#im4kj", tex)
	//})
	//if uiCtx.ActiveWidget == "#im4kj" {
	//	uiCtx.Tooltip("ttp-1", func() {
	//		uiCtx.Text("text-ttp-1", "Обычная картинка, которая  ничего не делает", 14)
	//		uiCtx.Text("text-ttp-2", "Hello World", 16)
	//		uiCtx.Text("text-ttp-3", "Hello World", 16)
	//	})
	//}
	//uiCtx.SubWidgetSpace("widhspdf-1", 100, 200, ui.NotResizable|ui.Scrollable|ui.ShowScrollbar, func() {
	//	uiCtx.Image("#im4kjdg464tht", 100, 100, tex2)
	//	uiCtx.Image("#im76erewr", 100, 100, tex)
	//	uiCtx.Text("#t3d79f", "world!", ui.DefaultTextFlag)
	//})
	//uiCtx.VSpace("#hhvs1")
	//uiCtx.Image("#imgj4", tex2)
	//uiCtx.VSpace("#dfff234")
	//uiCtx.ButtonT("sad3r3", "Hello!?")
	//uiCtx.Text("#t3dgdgdf", "world!", 24)
	//uiCtx.TabBar("bar1", func() {
	//	uiCtx.TabItem("Config", func() {
	//		uiCtx.Button("fdffdf")
	//		//uiCtx.PushStyleVar4f(ui.ButtonHoveredColor, [4]float32{100, 140, 76, 1})
	//		uiCtx.Button("fgfdffdf")
	//		//uiCtx.PopStyleVar()
	//		uiCtx.Text("textre-ttp-2", "Hello World", 16)
	//		uiCtx.Text("textrt-ttp-3", "Привет, мир!?", 16)
	//	})
	//	uiCtx.TabItem("Config 2", func() {
	//		uiCtx.SubWidgetSpace("widhswedf-1", 100, 200, ui.NotResizable|ui.Scrollable|ui.ShowScrollbar, func() {
	//			uiCtx.Image("#im4kjdg464tht", tex2)
	//			uiCtx.Image("#im76erewr", tex)
	//			uiCtx.Text("#t3df", "world!", 24)
	//		})
	//	})
	//	uiCtx.TabItem("Config 3", func() {
	//		uiCtx.Button("fdf4343545fdf")
	//		uiCtx.Text("te45xtаа", "Очень важная опция - ?", 16)
	//		uiCtx.VSpace("#hhvs1")
	//		uiCtx.Text("text23rtа", "2+2=4", 16)
	//		uiCtx.Image("#iваmgj4", tex)
	//	})
	//})
	//uiCtx.Image("#im4kjdg464tht", tex)
	//uiCtx.VSpace("#dfff234")

	uiCtx.EndWindow()
}

var selection int = 0
var sle = []string{"Hello", "Как выбирать?", "Белый", "Что-то очень длинное"}
var slCounter float32 = 0

func secondWindow() {
	uiCtx.BeginWindow("second wnd")
	uiCtx.Row("row 1dfdf14", func() {
		uiCtx.Image("#im4", 100, 100, tex.TextureId, tex.TexCoords)
		uiCtx.Image("#im4", 100, 100, tex.TextureId, tex.TexCoords)
	})

	cl := fmt.Sprintf("%.0f", slCounter)
	//uiCtx.Text("text-1dff", "The quick brown fox jumps over the lazy dog", 16)
	uiCtx.Text("text-1dff", "Съешь еще этих мягких", 16)
	uiCtx.Text("text-1dfhjyf", cl, 16)
	uiCtx.Slider("slider-1", &slCounter, 0, 255)

	uiCtx.Row("row 13214", func() {
		uiCtx.Image("#im4kjdg464", 100, 100, tex.TextureId, tex.TexCoords)
		uiCtx.Column("col fdfd", func() {
			uiCtx.Image("#im76", 100, 100, tex2.TextureId, tex2.TexCoords)
			uiCtx.Image("#im4", 100, 100, tex.TextureId, tex.TexCoords)
		})

		uiCtx.Column("col fdfdвава", func() {
			uiCtx.Button("ASsfdffb")
			uiCtx.Button("ASsfdffbbb")
		})

		uiCtx.Image("#im4kj", 100, 100, tex.TextureId, tex.TexCoords)
	})

	uiCtx.EndWindow()
}

func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	uiCtx.Io().MousePosCallback(float32(xpos), float32(ypos))
}

func mouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	uiCtx.Io().MouseBtnCallback(GlfwMouseKey(button), GlfwAction(action))
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	switch action {
	case glfw.Press:
		uiCtx.Io().KeyCallback(GlfwKeyToGuiKey(key), GlfwModKey(mods), true)
	case glfw.Release:
		uiCtx.Io().KeyCallback(GlfwKeyToGuiKey(key), GlfwModKey(mods), false)
	}
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	uiCtx.Io().ScrollX = xoff
	uiCtx.Io().ScrollY = yoff
}
func openImage(filepath string) image.Image {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		return nil
	}
	return img
}

func CreateImage(filename string, img image.Image) {
	pngFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	encoder.Encode(pngFile, img)
}

func ConvertFontToAtlas(f *fonts.Font, sheet *sprite_packer.SpriteSheet, srcImage *image.RGBA) {
	chars := make([]*fonts.CharInfo, len(f.CharMap))
	counter := 0
	for _, info := range f.CharMap {
		chars[counter] = info
		counter++
	}
	sort.Slice(chars, func(i, j int) bool {
		return chars[i].Width > chars[j].Width
	})

	sheet.BeginGroup(f.Filepath, func() map[string]*sprite_packer.SpriteInfo {
		spriteInfo := make(map[string]*sprite_packer.SpriteInfo, len(chars))
		//spriteInfo := make([]*sprite_packer.SpriteInfo, len(chars))
		for _, info := range chars {
			if info.Rune == ' ' || info.Rune == '\u00a0' {
				continue
			}
			ret := srcImage.SubImage(image.Rect(info.SrcX, info.SrcY, info.SrcX+info.Width, info.SrcY-info.Height)).(*image.RGBA)
			pixels := sheet.GetData(ret)
			spriteInfo[string(info.Rune)] = sheet.AddToSheet(string(info.Rune), pixels)
		}
		return spriteInfo
	})

	rr, ok := sheet.GetGroup(f.Filepath)
	if !ok {
		panic("doesnt exist")
	}

	for _, info := range rr.Contents {
		if info != nil {
			ll := []rune(info.Id)
			char := f.GetCharacter(ll[0])
			char.TexCoords = [2]utils.Vec2{{info.TextCoords[0], info.TextCoords[1]},
				{info.TextCoords[2], info.TextCoords[3]}}
		}
	}
}

func setCursor(c ui.CursorType) {
	defC := glfw.CreateStandardCursor(Cursor(c))
	window.SetCursor(defC)
}

var tex Texture
var tex2 Texture
var arrowDown Texture
var custWnd Texture

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}
