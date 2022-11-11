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
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"sort"
	"time"
)

func init() {
	runtime.LockOSThread()
}

var window *glfw.Window
var Width, Height int = 1280, 720
var steps int = 1

func startPprof() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/p", func(w http.ResponseWriter, r *http.Request) {
			pprof.Index(w, r)
		})
		mux.HandleFunc("/profile", pprof.Profile)
		mux.Handle("/heap", pprof.Handler("heap"))
		mux.Handle("/allocs", pprof.Handler("allocs"))
		mux.Handle("/goroutine", pprof.Handler("goroutine"))

		if err := http.ListenAndServe(":8543", mux); err != nil {
			fmt.Println("Can't start monitoring on port 8543")
			os.Exit(1)
		}
	}()
}

func main() {
	startPprof()
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
		ui.SetDisplaySize(float32(width), float32(height))
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

	front := NewGlRenderer()
	ui.AddRenderer(front)
	// uiCtx = ui.NewContext(front, cam)
	ui.SetChangeCursorFunc(setCursor)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//sheet := sprite_packer.NewSpriteSheet(512, "mgui")
	////ui.uiCtx.UploadFont("C:/Windows/Fonts/times.ttf", 14)
	fontName := "C:/Windows/Fonts/arial.ttf"
	fontName2 := "C:/Windows/Fonts/times.ttf"
	f, _ := ui.UploadFont(fontName2, 24, 157.0, 32, 256)
	f2, _ := ui.UploadFont(fontName, 18, 157.0, 32, 256)
	//d2.Pix = []uint8{}
	//d.Pix = []uint8{}
	//d = nil
	//d2 = nil
	//_ = d
	//_ = d2
	//ConvertFontToAtlas(f, sheet, data)
	//ConvertFontToAtlas(f2, sheet, data2)

	////ui.uiCtx.UploadFont("assets/fonts/rany.otf", 14)
	////ui.uiCtx.UploadFont("assets/fonts/sans.ttf", 18)
	////ui.uiCtx.UploadFont("assets/fonts/mono.ttf", 14)
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
	sheet.ClearImage()
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
	go func() {
		for {
			<-tmer.C

		}
	}()
	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

		ui.NewFrame([2]float32{float32(Width), float32(Height)})

		firstWindow()
		//secondWindow()
		//customWindow()

		if ui.GetIo().IsKeyPressed(ui.GuiKey_Space) {
			fmt.Println(ui.GET_CONTEXT())
			opendW = true
		}

		ui.EndFrame([2]float32{float32(Width), float32(Height)})

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
	ui.BeginCustomWindow("cstm wnd", 500, 500, 650, 650,
		100, 160, 400, 480,
		custWnd.TextureId, custWnd.TexCoords, func() {
			ui.Image("#im4kjdg464tht", 100, 100, tex.TextureId, tex.TexCoords)
			ui.Text("text-ttp-4", "Обычная", ui.Selectable)
		})
}

var opendW = true

var counter = 0
var dur int64 = 0
var tmer = time.NewTicker(1 * time.Second)
var sh = false

func firstWindow() {
	//ui.GlobalWidgetSpace("glob", 200, 200, 300, 300, ui.IgnoreClipping|ui.Resizable, func() {
	//	ui.Image("-imgy5g", 100, 100, tex.TextureId, tex.TexCoords)
	//	ui.Text("-txt", "dfdgfg - это текст-\"рыба\"", ui.Selectable)
	//})
	//ui.GlobalImage("glob", 100, 100, 100, 100, tex.TextureId, tex.TexCoords)
	ui.BeginWindow("The first window", &opendW)
	//uiCtx.Selection("sel-1", &selection, sle, arrowDown)
	//uiCtx.Selection("sel-1", &selection, sle, arrowDown)
	//uiCtx.Text("text-ttp-2", "Обычная картинка \nи это то-же 1", ui.Selectable)
	//uiCtx.Text("text-ttp-3", "Обычная картинка и \nэто то-же 2", ui.Editable)
	//start := time.Now()
	//ui.SubWidgetSpace("dfапаd", 200, 300, ui.Scrollable, func() {
	//for i := 0; i < 5; i++ {
	//	ui.Image(fmt.Sprint(i)+"-imgy5g", 100, 100, tex.TextureId, tex.TexCoords)
	//}
	//	//ui.Image(fmt.Sprint(12)+"-imgy5g", 100, 100, tex.TextureId, tex.TexCoords)
	//	ui.Text(fmt.Sprint(2131)+"-txt", "dfdgfg - это текст-\"рыба\"", ui.Selectable)
	//})
	//for i := 0; i < 100; i++ {
	//	//fmt.Println(tex.TexCoords)
	if ui.Image("-iyimgy5g", 100, 100, tex.TextureId, tex.TexCoords) {
		sh = !sh
	}
	ui.Text("h787-txt", "Lorem Ipsum - это текст-\"рыба\"", ui.Selectable)
	////}
	ui.Slider("slds", &tW, 100, 1200)
	//elapsed := time.Since(start)
	//dur += elapsed.Microseconds()
	//counter++
	if sh {
		ui.TextFitted("text-ttваы-1", tW, "Съешь ещё этих мягких французских булочек")
	}
	ui.Image("-iyimgydf5g", 100, 100, tex2.TextureId, tex2.TexCoords)
	ui.TextInput("tirey21", 300, 50, &msg1)
	//fmt.Printf("Widgets took %d \n", int(dur)/counter)
	{
		//ui.Image("#im4kjdg464tht", 100, 100, tex.TextureId, tex.TexCoords)
		//ui.Text("text-ttp-4", "Lorem Ipsum - это текст-\"рыба\", часто \nиспользуемый в печати и вэб-дизайне.", ui.Selectable)
		//uiCtx.Text("tlorem", "Lorem Ipsum - это текст-\"рыба\", часто \nиспользуемый в печати и вэб-дизайне. Lorem Ipsum является \nстандартной \"рыбой\" для текстов на \nлатинице с начала XVI века.", ui.Selectable)
		//ui.Button("fd")
		//ui.ButtonT("ds", "Sas")
		//ui.Row("slider-row", func() {
		//ui.Slider("slds", &tW, 100, 1200)
		//	ui.Text("sl-tex", fmt.Sprint(tW), ui.DefaultTextFlag)
		//})
		//ui.MultiLineTextInput("inputr23", &message)
		//ui.Text("text-ttp-43", "Обычная картинка и это то-же 123", ui.Selectable)
		//ui.TextInput("tinp-121", 300, 50, &msg1)
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
	}
	ui.EndWindow()
}

var selection int = 0
var sle = []string{"Hello", "Как выбирать?", "Белый", "Что-то очень длинное"}
var slCounter float32 = 0

func secondWindow() {
	ui.BeginWindow("second wnd", &opendW)
	ui.SubWidgetSpace("dfd", 100, 100, ui.Default|ui.Resizable, func() {
		ui.Image("-iytuyumgy5g", 100, 100, tex.TextureId, tex.TexCoords)
		ui.Text("-txfgjt", "the quick brown fox \njumps over the lazy dog", ui.Selectable)
	})
	//ui.Row("row 1dfdf14", func() {
	//	ui.Image("#im4", 100, 100, tex.TextureId, tex.TexCoords)
	//	ui.Image("#im4", 100, 100, tex.TextureId, tex.TexCoords)
	//})

	//cl := fmt.Sprintf("%.0f", slCounter)
	//uiCtx.Text("text-1dff", "The quick brown fox jumps over the lazy dog", 16)
	//ui.Text("text-1dff", "Съешь еще этих мягких", 16)
	//ui.Text("text-1dfhjyf", cl, 16)
	//ui.Slider("slider-1", &slCounter, 0, 255)
	//
	//ui.Row("row 13214", func() {
	//	ui.Image("#im4kjdg464", 100, 100, tex.TextureId, tex.TexCoords)
	//	ui.Column("col fdfd", func() {
	//		ui.Image("#im76", 100, 100, tex2.TextureId, tex2.TexCoords)
	//		ui.Image("#im4", 100, 100, tex.TextureId, tex.TexCoords)
	//	})
	//
	//	ui.Column("col fdfdвава", func() {
	//		ui.Button("ASsfdffb")
	//		ui.Button("ASsfdffbbb")
	//	})
	//
	//	ui.Image("#im4kj", 100, 100, tex.TextureId, tex.TexCoords)
	//})

	ui.EndWindow()
}

func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	ui.GetIo().MousePosCallback(float32(xpos), float32(ypos))
}

func mouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	ui.GetIo().MouseBtnCallback(GlfwMouseKey(button), GlfwAction(action))
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	switch action {
	case glfw.Press:
		ui.GetIo().KeyCallback(GlfwKeyToGuiKey(key), GlfwModKey(mods), true)
	case glfw.Release:
		ui.GetIo().KeyCallback(GlfwKeyToGuiKey(key), GlfwModKey(mods), false)
	}
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	ui.GetIo().ScrollX = xoff
	ui.GetIo().ScrollY = yoff
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
