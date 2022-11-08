package draw

import (
	"fmt"
	"math"

	"github.com/Dmitry-dms/mgui/fonts"
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/Dmitry-dms/mgui/widgets"
)

type CmdBuffer struct {
	commands []Command

	displaySize *utils.Vec2

	Inf []Info
	ofs int

	Vertices []float32
	Indeces  []int32

	VertCount int
	lastIndc  int
	lastElems int

	//InnerWindowSpace [4]float32
}

type Info struct {
	Elems       int
	IndexOffset int
	TexId       uint32
	Type        string
	//Clip        ClipRectCompose
	ClipRect [4]float32
	//MainClipRect [4]float32
}

type ClipRectCompose struct {
	ClipRect     [4]float32 // only used in sub widget spaces
	MainClipRect [4]float32 // main window clipping rectangle
}

var EmptyClip = [4]float32{0, 0, 0, 0}

func NewClip(inner, main [4]float32) ClipRectCompose {
	c := ClipRectCompose{
		ClipRect:     inner,
		MainClipRect: main,
	}
	return c
}

func NewBuffer(size *utils.Vec2) *CmdBuffer {
	return &CmdBuffer{
		commands:    []Command{},
		Vertices:    []float32{},
		Indeces:     []int32{},
		displaySize: size,
		VertCount:   0,
	}
}

func (c *CmdBuffer) Clear() {
	c.commands = []Command{}
	c.Vertices = []float32{}
	c.Indeces = []int32{}
	c.VertCount = 0
	c.lastIndc = 0
	c.ofs = 0
	c.Inf = []Info{}
	c.lastElems = 0
}

func (c *CmdBuffer) SeparateBuffer(texId uint32, clip ClipRectCompose) {
	cmd := Command{
		Type: SeparateBuffer,
		sb: &separate_buff{
			texid:        texId,
			clipRect:     clip.ClipRect,
			mainClipRect: clip.MainClipRect,
		},
	}
	c.AddCommand2(cmd, clip)
}

func (c *CmdBuffer) CreateButtonT(x, y float32, btn *widgets.TextButton, font fonts.Font, clip ClipRectCompose) {
	c.CreateRect(x, y, btn.Button.Width(), btn.Button.Height(), 0, StraightCorners, 0, btn.Color(), clip)
	btn.UpdateTextPos(x, y)
	c.CreateText(btn.Text.BoundingBox()[0], btn.Text.BoundingBox()[1], btn.Text, font, clip)
}
func (c *CmdBuffer) CreateText2(x, y float32, txt *widgets.Text, scale float32, font fonts.Font) ([]float32, []int32, int, int) {

	//c.CreateRect(x, y, txt.Width(), txt.Height(), 0, StraightCorners, 0, txt.BackgroundColor(), clip)
	vert, ind, count := c.Text2(txt, font, x, c.displaySize.Y-y, scale, txt.CurrentColor)
	return vert, ind, count, c.lastIndc
	//c.SeparateBuffer(t.Font.TextureId, clip) // don't forget to slice buffer
	//c.CreateBorderBox(x, y, txt.Width(), txt.Height(), 2, [4]float32{255, 0, 0, 1})
	//c.SeparateBuffer(font.TextureId, clip)
}
func (c *CmdBuffer) CreateText(x, y float32, txt *widgets.Text, font fonts.Font, clip ClipRectCompose) {
	tcmd := Text_command{
		X:       x,
		Y:       y,
		Clr:     txt.CurrentColor,
		Text:    txt.Message,
		Font:    font,
		Scale:   txt.Scale,
		Padding: txt.Padding,
		Widget:  txt,
	}
	cmd := Command{
		Text:     &tcmd,
		Type:     Text,
		WidgetId: txt.WidgetId(),
	}

	c.CreateRect(x, y, txt.Width(), txt.Height(), 0, StraightCorners, 0, txt.BackgroundColor(), clip)
	c.AddCommand(cmd, clip)
	//c.CreateBorderBox(x, y, txt.Width(), txt.Height(), 2, [4]float32{255, 0, 0, 1})
	//c.SeparateBuffer(font.TextureId, clip)
}

func (c *CmdBuffer) CreateWindow(wnd Window_command, clip ClipRectCompose) {
	c.CreateRect(wnd.X, wnd.Y, wnd.W, wnd.H, 10, AllRounded, 0, wnd.Clr, clip)
	toolbar := wnd.Toolbar
	c.CreateRect(toolbar.X, toolbar.Y, toolbar.W, toolbar.H, 10, TopRect, 0, toolbar.Clr, clip)

	//c.SeparateBuffer(0, clip)
}
func (c *CmdBuffer) CreateTexturedRect2(x, y, w, h float32, texId uint32, coords, clr [4]float32, clip ClipRectCompose) ([]float32, []int32, int, int) {
	//cmd := Command{
	//	Type: RectType,
	//	Rect: &Rect_command{
	//		X:      x,
	//		Y:      y,
	//		W:      w,
	//		H:      h,
	//		Clr:    clr,
	//		TexId:  texId,
	//		coords: coords,
	//	},
	//}
	//return c.AddCommand2(cmd, clip)

	vert, ind, count := c.RectangleT2(x, c.displaySize.Y-y, w, h, texId, coords, clr)
	//c.SeparateBuffer(texId, clip)
	//fmt.Println(y, ind)
	return vert, ind, count, c.lastIndc
}
func (c *CmdBuffer) CreateTexturedRect(x, y, w, h float32, texId uint32, coords, clr [4]float32, clip ClipRectCompose) {
	cmd := Command{
		Type: RectType,
		Rect: &Rect_command{
			X:      x,
			Y:      y,
			W:      w,
			H:      h,
			Clr:    clr,
			TexId:  texId,
			coords: coords,
		},
	}
	c.AddCommand(cmd, clip)
}

func (c *CmdBuffer) CreateRect(x, y, w, h float32, radius int, shape RoundedRectShape, texId uint32, clr [4]float32, clip ClipRectCompose) {
	cmd := Command{
		Type: RectType,
		Rect: &Rect_command{
			X:      x,
			Y:      y,
			W:      w,
			H:      h,
			Clr:    clr,
			radius: radius,
			TexId:  texId,
			shape:  shape,
		},
	}
	c.AddCommand(cmd, clip)
}
func checkSliceForNull(s [4]float32) bool {
	return (s[0] == 0) && (s[1] == 0) && (s[2] == 0) && (s[3] == 0)
}
func (c *CmdBuffer) AddCommand2(cmd Command, clip ClipRectCompose) (vert []float32, ind []int32, count int) {
	c.commands = append(c.commands, cmd)

	switch cmd.Type {
	case SeparateBuffer:
		mainRect := clip.MainClipRect
		innerRect := clip.ClipRect

		x, x2 := int32(mainRect[0]), int32(innerRect[0])
		y, y2 := int32(mainRect[1]), int32(innerRect[1])
		w, w2 := int32(mainRect[2]), int32(innerRect[2])
		h, h2 := int32(mainRect[3]), int32(innerRect[3])

		useInnerClip := !checkSliceForNull(innerRect)
		xl := x+w < x2+w2
		yl := y+h < y2+h2

		overlapWidth := useInnerClip && xl
		overlapHeigth := useInnerClip && yl

		inf := Info{
			Elems:       c.VertCount - c.lastElems,
			IndexOffset: c.ofs,
			TexId:       cmd.sb.texid,
		}
		if !useInnerClip {
			inf.ClipRect = cmd.sb.mainClipRect
		} else if overlapWidth && overlapHeigth {
			inf.ClipRect = cmd.sb.mainClipRect
		} else if overlapWidth {
			inf.ClipRect = cmd.sb.mainClipRect
		} else if overlapHeigth {
			var tmp = cmd.sb.clipRect
			inf.ClipRect = [4]float32{tmp[0], tmp[1], tmp[2], cmd.sb.mainClipRect[3] - (tmp[1] - cmd.sb.mainClipRect[1])}
		} else {
			inf.ClipRect = cmd.sb.clipRect
		}

		c.Inf = append(c.Inf, inf)
		c.ofs += c.VertCount - c.lastElems
		c.lastElems = c.VertCount
	case RectType:
		r := cmd.Rect
		//if r.radius == 0 {
		if r.TexId == 0 {
			c.RectangleR(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.Clr)

		} else {
			vert, ind, count = c.RectangleT2(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.TexId, r.coords, r.Clr)
			c.SeparateBuffer(r.TexId, clip) // don't forget to slice buffer
			fmt.Println(r.Y, ind)
			return
		}
		//} else {
		//	if r.TexId == 0 {
		//		c.RoundedRectangleR(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.radius, r.shape, r.Clr)
		//	} else {
		//		// TODO: Add textured rounded rect
		//	}
		//}
	case Text:
		t := cmd.Text
		c.Text(t.Widget, t.Font, t.X, c.displaySize.Y-(t.Y+float32(t.Padding)), t.Scale, t.Clr)
		c.SeparateBuffer(t.Font.TextureId, clip) // don't forget to slice buffer
		//case BezierQuad:
		//	b := cmd.Bezier
		//	c.DrawBezierQuad(b.StartX, b.StartY, b.SupportX, b.SupportY, b.EndX, b.EndY, b.Steps, b.Clr, clip)
		//	c.sepBuf(clip, "LINE_STRIP")
		//case Line:
		//	l := cmd.Line
		//	c.DrawLine(l.StartX, c.displaySize.Y-l.StartY, l.EndX, c.displaySize.Y-l.EndY, l.Clr)
		//	c.sepBuf(clip, "LINE")
		//case LineStrip:
		//	l := cmd.Line
		//	changed := make([]utils.Vec2, len(l.Points))
		//	for i, p := range l.Points {
		//		changed[i].Y = c.displaySize.Y - p.Y
		//		changed[i].X = p.X
		//	}
		//	c.DrawLineStrip(l.Clr, changed)
		//	c.sepBuf(clip, "LINE_STRIP")
	}
	//fmt.Println(len(vert), len(ind), count)
	return
}
func (c *CmdBuffer) AddCommand(cmd Command, clip ClipRectCompose) {
	c.commands = append(c.commands, cmd)

	switch cmd.Type {
	case SeparateBuffer:
		mainRect := clip.MainClipRect
		innerRect := clip.ClipRect

		x, x2 := int32(mainRect[0]), int32(innerRect[0])
		y, y2 := int32(mainRect[1]), int32(innerRect[1])
		w, w2 := int32(mainRect[2]), int32(innerRect[2])
		h, h2 := int32(mainRect[3]), int32(innerRect[3])

		useInnerClip := !checkSliceForNull(innerRect)
		xl := x+w < x2+w2
		yl := y+h < y2+h2

		overlapWidth := useInnerClip && xl
		overlapHeigth := useInnerClip && yl

		inf := Info{
			Elems:       c.VertCount - c.lastElems,
			IndexOffset: c.ofs,
			TexId:       cmd.sb.texid,
		}
		if !useInnerClip {
			inf.ClipRect = cmd.sb.mainClipRect
		} else if overlapWidth && overlapHeigth {
			inf.ClipRect = cmd.sb.mainClipRect
		} else if overlapWidth {
			inf.ClipRect = cmd.sb.mainClipRect
		} else if overlapHeigth {
			var tmp = cmd.sb.clipRect
			inf.ClipRect = [4]float32{tmp[0], tmp[1], tmp[2], cmd.sb.mainClipRect[3] - (tmp[1] - cmd.sb.mainClipRect[1])}
		} else {
			inf.ClipRect = cmd.sb.clipRect
		}

		c.Inf = append(c.Inf, inf)
		c.ofs += c.VertCount - c.lastElems
		c.lastElems = c.VertCount
	case RectType:
		r := cmd.Rect
		if r.radius == 0 {
			if r.TexId == 0 {
				c.RectangleR(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.Clr)
			} else {
				c.RectangleT(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.TexId, r.coords, r.Clr)
				c.SeparateBuffer(r.TexId, clip) // don't forget to slice buffer
			}
		} else {
			if r.TexId == 0 {
				c.RoundedRectangleR(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.radius, r.shape, r.Clr)
			} else {
				// TODO: Add textured rounded rect
			}
		}
	case Text:
		t := cmd.Text
		c.Text(t.Widget, t.Font, t.X, c.displaySize.Y-(t.Y+float32(t.Padding)), t.Scale, t.Clr)
		c.SeparateBuffer(t.Font.TextureId, clip) // don't forget to slice buffer
	case BezierQuad:
		b := cmd.Bezier
		c.DrawBezierQuad(b.StartX, b.StartY, b.SupportX, b.SupportY, b.EndX, b.EndY, b.Steps, b.Clr, clip)
		c.sepBuf(clip, "LINE_STRIP")
	case Line:
		l := cmd.Line
		c.DrawLine(l.StartX, c.displaySize.Y-l.StartY, l.EndX, c.displaySize.Y-l.EndY, l.Clr)
		c.sepBuf(clip, "LINE")
	case LineStrip:
		l := cmd.Line
		changed := make([]utils.Vec2, len(l.Points))
		for i, p := range l.Points {
			changed[i].Y = c.displaySize.Y - p.Y
			changed[i].X = p.X
		}
		c.DrawLineStrip(l.Clr, changed)
		c.sepBuf(clip, "LINE_STRIP")
	}
}
func (c *CmdBuffer) sepBuf(clip ClipRectCompose, t string) {
	inf := Info{
		Elems:       c.VertCount - c.lastElems,
		IndexOffset: c.ofs,
		TexId:       0,
		Type:        t,
	}
	inf.ClipRect = clip.MainClipRect
	c.Inf = append(c.Inf, inf)
	c.ofs += c.VertCount - c.lastElems
	c.lastElems = c.VertCount
}
func (c *CmdBuffer) Render(vert []float32, indeces []int32, vertCount int) {
	c.Vertices = append(c.Vertices, vert...)
	c.Indeces = append(c.Indeces, indeces...)
	c.VertCount += vertCount
}
func (c *CmdBuffer) Render2(vert []float32, indeces []int32, vertCount int, ind int) {
	c.Vertices = append(c.Vertices, vert...)
	c.Indeces = append(c.Indeces, indeces...)
	c.VertCount += vertCount
	if ind != 0 {
		c.lastIndc = ind
	}
}

func (c *CmdBuffer) CreateBorderBox(x, y, w, h, lineWidth float32, clr [4]float32) {
	c.CreateRect(x, y, w, lineWidth, 0, StraightCorners, 0, clr, ClipRectCompose{})
	c.CreateRect(x+w-lineWidth, y, lineWidth, h, 0, StraightCorners, 0, clr, ClipRectCompose{})
	c.CreateRect(x, y, lineWidth, h, 0, StraightCorners, 0, clr, ClipRectCompose{})
	c.CreateRect(x, y+h-lineWidth, w, lineWidth, 0, StraightCorners, 0, clr, ClipRectCompose{})
}

func (c *CmdBuffer) RectangleR(x, y, w, h float32, clr [4]float32) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ind0 := c.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, x, y-h, 0, 0, 0, clr)
	fillVertices(vert, &offset, x+w, y-h, 0, 0, 0, clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, 0, 0, 0, clr)
	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	c.lastIndc = last + 1
	c.Render(vert, ind, 6)
}
func (c *CmdBuffer) Text2(text *widgets.Text, font fonts.Font, x, y float32, scale float32, clr [4]float32) ([]float32, []int32, int) {
	texId := font.TextureId
	vert := []float32{}
	ind := []int32{}
	cnt := 0
	for i, r := range []rune(text.Message) {
		info := font.GetCharacter(r)

		if info.Rune == rune(127) { // '\n'
			continue
		}
		xPos := x + text.Chars[i].Pos.X
		yPos := y - text.Chars[i].Pos.Y
		v, idec, vc := c.addCharacter2(xPos, yPos, scale, texId, *info, clr)
		vert = append(vert, v...)
		ind = append(ind, idec...)
		cnt += vc
	}
	return vert, ind, cnt
}
func (c *CmdBuffer) addCharacter2(x, y float32, scale float32, texId uint32, info fonts.CharInfo, clr [4]float32) ([]float32, []int32, int) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	x0 := x
	y0 := y
	x1 := x + scale*float32(info.Width)
	y1 := y + scale*float32(info.Height)

	ux0, uy0 := info.TexCoords[0].X, info.TexCoords[0].Y
	ux1, uy1 := info.TexCoords[1].X, info.TexCoords[1].Y

	ind0 := c.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x1, y0, ux1, uy0, float32(texId), clr)
	fillVertices(vert, &offset, x1, y1, ux1, uy1, float32(texId), clr)
	fillVertices(vert, &offset, x0, y1, ux0, uy1, float32(texId), clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x0, y0, ux0, uy0, float32(texId), clr)

	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	c.lastIndc = last + 1
	//c.Render(vert, ind, 6)
	return vert, ind, 6
}
func (c *CmdBuffer) Text(text *widgets.Text, font fonts.Font, x, y float32, scale float32, clr [4]float32) {
	texId := font.TextureId
	for i, r := range []rune(text.Message) {
		info := font.GetCharacter(r)

		if info.Rune == rune(127) { // '\n'
			continue
		}
		xPos := x + text.Chars[i].Pos.X
		yPos := y - text.Chars[i].Pos.Y
		c.addCharacter(xPos, yPos, scale, texId, *info, clr)
	}
}

func (c *CmdBuffer) addCharacter(x, y float32, scale float32, texId uint32, info fonts.CharInfo, clr [4]float32) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	x0 := x
	y0 := y
	x1 := x + scale*float32(info.Width)
	y1 := y + scale*float32(info.Height)

	ux0, uy0 := info.TexCoords[0].X, info.TexCoords[0].Y
	ux1, uy1 := info.TexCoords[1].X, info.TexCoords[1].Y

	ind0 := c.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x1, y0, ux1, uy0, float32(texId), clr)
	fillVertices(vert, &offset, x1, y1, ux1, uy1, float32(texId), clr)
	fillVertices(vert, &offset, x0, y1, ux0, uy1, float32(texId), clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x0, y0, ux0, uy0, float32(texId), clr)

	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	c.lastIndc = last + 1
	c.Render(vert, ind, 6)
}
func (c *CmdBuffer) RectangleT2(x, y, w, h float32, texId uint32, coords [4]float32, clr [4]float32) ([]float32, []int32, int) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ux0, uy0 := coords[2], coords[3]
	ux1, uy1 := coords[0], coords[1]

	ind0 := c.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, ux1, uy0, float32(texId), clr)
	fillVertices(vert, &offset, x, y-h, ux1, uy1, float32(texId), clr)
	fillVertices(vert, &offset, x+w, y-h, ux0, uy1, float32(texId), clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, ux0, uy0, float32(texId), clr)

	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	c.lastIndc = last + 1
	//c.Render(vert, ind, 6)
	return vert, ind, 6
}

func (c *CmdBuffer) RectangleT(x, y, w, h float32, texId uint32, coords [4]float32, clr [4]float32) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ux0, uy0 := coords[2], coords[3]
	ux1, uy1 := coords[0], coords[1]

	ind0 := c.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, ux1, uy0, float32(texId), clr)
	fillVertices(vert, &offset, x, y-h, ux1, uy1, float32(texId), clr)
	fillVertices(vert, &offset, x+w, y-h, ux0, uy1, float32(texId), clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, ux0, uy0, float32(texId), clr)

	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	c.lastIndc = last + 1
	c.Render(vert, ind, 6)
}

func fillVertices(vert []float32, startOffset *int, x, y, uv0, uv1, texId float32, clr [4]float32) {
	offset := *startOffset
	vert[offset] = x
	vert[offset+1] = y

	vert[offset+2] = clr[0] / 255
	vert[offset+3] = clr[1] / 255
	vert[offset+4] = clr[2] / 255
	vert[offset+5] = clr[3]

	vert[offset+6] = uv0
	vert[offset+7] = uv1

	vert[offset+8] = texId

	*startOffset += 9
}

type CircleSector int
type RoundedRectShape int

const (
	TopLeftRect RoundedRectShape = 1 << iota
	TopRightRect
	BotLeftRect
	BotRightRect
	OnlyBorders
	StraightCorners

	TopRect = TopLeftRect | TopRightRect
	BotRect = BotLeftRect | BotRightRect

	AllRounded = TopRect | BotRect
)

const (
	BotLeft CircleSector = iota
	BotRight
	TopLeft
	TopRight
)

func (c *CmdBuffer) DrawArc(x, y, radius float32, steps int, sector CircleSector, clr [4]float32) {
	ind0 := c.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0
	indOffset := 0

	angle := math.Pi * 2 / float32(steps)

	numV := int(math.Floor(1.57 / float64(angle)))

	ind := make([]int32, 3*(numV+1))    // 3 - triangle
	vert := make([]float32, 9*(3+numV)) //polygon

	var prevX, prevY, lastX, lastY float32

	var ang float32 = angle
	var sX func(x, radius float32) float32
	var sY func(y, radius float32) float32
	// counterTriangles := 0
	switch sector {
	case BotLeft:
		sX = func(x, ang float32) float32 {
			return x - float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y - float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y - radius
		lastX = x - radius
		lastY = y
	case BotRight:
		sX = func(x, ang float32) float32 {
			return x + float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y - float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y - radius
		lastX = x + radius
		lastY = y
	case TopLeft:
		sX = func(x, ang float32) float32 {
			return x - float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y + float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y + radius
		lastX = x - radius
		lastY = y
	case TopRight:
		sX = func(x, ang float32) float32 {
			return x + float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y + float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y + radius
		lastX = x + radius
		lastY = y
	}

	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, prevX, prevY, 0, 0, 0, clr)
	newx := sX(x, ang)
	newY := sY(y, ang)
	fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)
	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	indOffset += 3
	// ind = append(ind, int32(ind0), int32(ind1), int32(ind2))
	ind1++
	ind2++
	ang += angle

	vertC := 1
	for ang <= 1.57 { // 90 degress ~= 1.57 radians
		newx := sX(x, ang)
		newY := sY(y, ang)

		fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)

		ind[indOffset] = int32(ind0)
		ind[indOffset+1] = int32(ind1)
		ind[indOffset+2] = int32(ind2)
		indOffset += 3
		ind1++
		ind2++

		ang += angle
		vertC++
		// counterTriangles++
	}
	fillVertices(vert, &offset, lastX, lastY, 0, 0, 0, clr)

	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	// indOffset += 3

	c.lastIndc = ind2 + 1

	c.Render(vert, ind, (numV+1)*3)
}
func (c *CmdBuffer) CreateLineArc(x, y, radius, steps float32, sector CircleSector, clr [4]float32, clip ClipRectCompose) {
	switch sector {
	case BotLeft:
		c.CreateBezierQuad(x, y+radius, x-radius, y+radius, x-radius, y, steps, clr, clip)
	case BotRight:
		c.CreateBezierQuad(x+radius, y, x+radius, y+radius, x, y+radius, steps, clr, clip)
	case TopLeft:
		c.CreateBezierQuad(x, y-radius, x-radius, y-radius, x-radius, y, steps, clr, clip)
	case TopRight:
		c.CreateBezierQuad(x, y-radius, x+radius, y-radius, x+radius, y, steps, clr, clip)
	}
}

func (c *CmdBuffer) CreateBezierQuad(startX, startY, supportX, supportY, endX, endY, steps float32, clr [4]float32, clip ClipRectCompose) {
	cmd := Command{
		Type: BezierQuad,
		Bezier: &BezierQuad_command{
			StartX:   startX,
			StartY:   startY,
			SupportX: supportX,
			SupportY: supportY,
			EndX:     endX,
			EndY:     endY,
			Clr:      clr,
			Steps:    steps,
		},
	}
	c.AddCommand(cmd, clip)
}

// TODO: Line drawing needs an optimization beacuse now, each line takes 1 draw call
func (c *CmdBuffer) CreateLine(startX, startY, endX, endY float32, clr [4]float32, clip ClipRectCompose) {
	cmd := Command{
		Type: Line,
		Line: &Line_command{
			StartX: startX,
			StartY: startY,
			EndX:   endX,
			EndY:   endY,
			Clr:    clr,
		},
	}
	c.AddCommand(cmd, clip)
}
func (c *CmdBuffer) CreateLineStrip(p []utils.Vec2, clr [4]float32, clip ClipRectCompose) {
	cmd := Command{
		Type: LineStrip,
		Line: &Line_command{
			Points: p,
			Clr:    clr,
		},
	}
	c.AddCommand(cmd, clip)
}

func (c *CmdBuffer) DrawLine(startX, startY, endX, endY float32, clr [4]float32) {
	ind0 := c.lastIndc
	offset := 0

	ind := make([]int32, 2)      // 1 - point
	vert := make([]float32, 9*2) //polygon

	fillVertices(vert, &offset, startX, startY, 0, 0, 0, clr)
	ind[0] = int32(ind0)
	ind0++
	fillVertices(vert, &offset, endX, endY, 0, 0, 0, clr)
	ind[1] = int32(ind0)

	c.lastIndc = ind0 + 1
	c.Render(vert, ind, 2)
}
func (c *CmdBuffer) DrawLineStrip(clr [4]float32, points []utils.Vec2) {
	ind0 := c.lastIndc
	offset := 0
	pointsLen := len(points)
	ind := make([]int32, pointsLen)      // 1 - point
	vert := make([]float32, 9*pointsLen) //polygon
	for i, point := range points {
		fillVertices(vert, &offset, point.X, point.Y, 0, 0, 0, clr)
		ind[i] = int32(ind0)
		ind0++
	}
	c.lastIndc = ind0
	c.Render(vert, ind, pointsLen)
}

func (c *CmdBuffer) DrawBezierQuad(startX, startY, supportX, supportY, endX, endY, steps float32, clr [4]float32, clip ClipRectCompose) {
	bezierQuad := func(t float32) (float32, float32) {
		v1 := float32(math.Pow(float64(1-t), 2))
		v2 := 2 * t * (1 - t)
		v3 := float32(math.Pow(float64(t), 2))
		return v1*startX + v2*supportX + v3*endX, v1*startY + v2*supportY + v3*endY
	}
	acc := float64(1 / steps)
	points := make([]utils.Vec2, int(steps)+1)
	ind := 0
	for t := .0; t < 1.0; t += acc {
		x, y := bezierQuad(float32(t))
		points[ind] = utils.Vec2{x, y}
		ind++
	}
	points[ind] = utils.Vec2{endX, endY}
	c.CreateLineStrip(points, clr, clip)
}

var steps = 30

func (c *CmdBuffer) RoundedBorderRectangle(x, y, w, h, lineWidth float32, radius int, clr [4]float32, clip ClipRectCompose) {
	c.roundedLineRectangle(x, y, w, h, lineWidth, 10, radius, AllRounded, clr, clip)
}
func (c *CmdBuffer) roundedLineRectangle(x, y, w, h, lineWidth, steps float32, radius int, shape RoundedRectShape, clr [4]float32, clip ClipRectCompose) {

	topLeft := utils.Vec2{x + float32(radius), y + float32(radius)} //origin of arc
	topRight := utils.Vec2{x + w - float32(radius), y + float32(radius)}
	botLeft := utils.Vec2{x + float32(radius), y + h - float32(radius)}
	botRight := utils.Vec2{x + w - float32(radius), y + h - float32(radius)}

	switch shape {
	case TopLeftRect:
		c.CreateLineArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr, clip)
		c.CreateLineStrip([]utils.Vec2{
			{topLeft.X, y},
			{x + w, y},
			{x + w, y + h},
			{x, y + h},
			{x, topLeft.Y},
		}, clr, clip)
	case TopRightRect:
		c.CreateLineArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr, clip)
		c.CreateLineStrip([]utils.Vec2{
			{topRight.X + float32(radius), y + float32(radius)},
			{topRight.X + float32(radius), topRight.Y + h - float32(radius)},
			{x, y + h},
			{x, y},
			{topRight.X, y},
		}, clr, clip)

	case BotLeftRect:
		c.CreateLineArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr, clip)
		c.CreateLineStrip([]utils.Vec2{
			{x, botLeft.Y},
			{x, y},
			{x + w, y},
			{x + w, y + h},
			{botLeft.X, botLeft.Y + float32(radius)},
		}, clr, clip)
	case BotRightRect:
		c.CreateLineArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr, clip)
		c.CreateLineStrip([]utils.Vec2{
			{botRight.X, botRight.Y + float32(radius)},
			{x, y + h},
			{x, y},
			{x + w, y},
			{x + w, y + h - float32(radius)},
		}, clr, clip)
	case TopRect:
		c.CreateLineArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr, clip)
		c.CreateLine(topLeft.X, y, topRight.X, y, clr, clip)
		c.CreateLineArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr, clip)
		c.CreateLineStrip([]utils.Vec2{
			{x + w, topRight.Y},
			{x + w, y + h},
			{x, y + h},
			{x, y + float32(radius)},
		}, clr, clip)
	case BotRect:
		c.CreateLineArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr, clip)
		c.CreateLine(botRight.X, y+h, botLeft.X, y+h, clr, clip)
		c.CreateLineArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr, clip)
		c.CreateLineStrip([]utils.Vec2{
			{x, botLeft.Y},
			{x, y},
			{x + w, y},
			{x + w, botRight.Y},
		}, clr, clip)
	case AllRounded:
		c.CreateLineArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr, clip)
		c.CreateLine(topLeft.X, y, topRight.X, y, clr, clip)
		c.CreateLineArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr, clip)
		c.CreateLine(x+w, topRight.Y, x+w, botRight.Y, clr, clip)
		c.CreateLineArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr, clip)
		c.CreateLine(botRight.X, y+h, botLeft.X, y+h, clr, clip)
		c.CreateLineArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr, clip)
		c.CreateLine(x, botLeft.Y, x, topLeft.Y, clr, clip)
	}

}

func (c *CmdBuffer) RoundedRectangleR(x, y, w, h float32, radius int, shape RoundedRectShape, clr [4]float32) {

	topLeft := utils.Vec2{x + float32(radius), y - float32(radius)} //origin of arc
	topRight := utils.Vec2{x + w - float32(radius), y - float32(radius)}
	botLeft := utils.Vec2{x + float32(radius), y - h + float32(radius)}
	botRight := utils.Vec2{x + w - float32(radius), y - h + float32(radius)}

	switch shape {
	case TopLeftRect:
		c.DrawArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		c.RectangleR(x, y-float32(radius), w, h-float32(radius), clr)               //main rect
		c.RectangleR(x+float32(radius), y, w-float32(radius), float32(radius), clr) //top rect
	case TopRightRect:
		c.DrawArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		c.RectangleR(x, y-float32(radius), w, h-float32(radius), clr) //main
		c.RectangleR(x, y, w-float32(radius), float32(radius), clr)
	case BotLeftRect:
		c.DrawArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		c.RectangleR(x, y, w, h-float32(radius), clr) //main
		c.RectangleR(botLeft.X, botLeft.Y, w-float32(radius), float32(radius), clr)
	case BotRightRect:
		c.DrawArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		c.RectangleR(x, y, w, h-float32(radius), clr) //main
		c.RectangleR(x, botLeft.Y, w-float32(radius), float32(radius), clr)
	case TopRect:
		c.DrawArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		c.DrawArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		c.RectangleR(x, y-float32(radius), w, h-float32(radius), clr)
		c.RectangleR(x+float32(radius), y, w-float32(radius)*2, float32(radius), clr)
	case BotRect:
		c.DrawArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		c.DrawArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		c.RectangleR(x, y, w, h-float32(radius), clr) //main
		c.RectangleR(botLeft.X, botLeft.Y, w-float32(radius)*2, float32(radius), clr)
	case AllRounded:
		c.DrawArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		c.DrawArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		c.DrawArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		c.DrawArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)

		c.RectangleR(topLeft.X, topLeft.Y+float32(radius), w-float32(radius)*2, float32(radius), clr) //top
		c.RectangleR(x, topLeft.Y, w, h-float32(radius)*2, clr)                                       //center
		c.RectangleR(botLeft.X, botLeft.Y, w-float32(radius)*2, float32(radius), clr)                 //bottom
	}

}
