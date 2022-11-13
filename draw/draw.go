package draw

import (
	"math"

	"github.com/Dmitry-dms/mgui/fonts"
	"github.com/Dmitry-dms/mgui/utils"
	"github.com/Dmitry-dms/mgui/widgets"
)

type CmdBuffer struct {
	displaySize *utils.Vec2

	DrawCalls []DrawCall
	ofs       int

	Vertices []float32
	Indices  []int32

	VertCount int
	LastInd   int
	lastElems int

	//InnerWindowSpace [4]float32
}

type DrawCall struct {
	Elems       int
	IndexOffset int
	TexId       uint32
	Type        string
	ClipRect    [4]float32
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
		Vertices:    make([]float32, 0, 10000),
		Indices:     make([]int32, 0, 10000),
		displaySize: size,
		VertCount:   0,
	}
}

func (c *CmdBuffer) Clear() {
	c.Vertices = c.Vertices[:0] // Reuse slices to prevent huge amount of re-allocations
	c.Indices = c.Indices[:0]
	c.DrawCalls = c.DrawCalls[:0]
	c.VertCount = 0
	c.LastInd = 0
	c.ofs = 0
	c.lastElems = 0
}

func (c *CmdBuffer) SeparateBuffer(texId uint32, clip ClipRectCompose) {
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

	inf := DrawCall{
		Elems:       c.VertCount - c.lastElems,
		IndexOffset: c.ofs,
		TexId:       texId,
	}
	if !useInnerClip {
		inf.ClipRect = mainRect
	} else if overlapWidth && overlapHeigth {
		inf.ClipRect = mainRect
	} else if overlapWidth {
		inf.ClipRect = mainRect
	} else if overlapHeigth {
		var tmp = innerRect
		inf.ClipRect = [4]float32{tmp[0], tmp[1], tmp[2], mainRect[3] - (tmp[1] - mainRect[1])}
	} else {
		inf.ClipRect = innerRect
	}

	c.DrawCalls = append(c.DrawCalls, inf)
	c.ofs += c.VertCount - c.lastElems
	c.lastElems = c.VertCount
}

func (c *CmdBuffer) CreateButtonT(x, y float32, btn *widgets.TextButton, font fonts.Font, clip ClipRectCompose) {
	//c.CreateRect(x, y, btn.Button.Width(), btn.Button.Height(), 0, StraightCorners, 0, btn.Color(), clip)
	btn.UpdateTextPos(x, y)
	//c.CreateText(btn.Text.BoundingBox()[0], btn.Text.BoundingBox()[1], btn.Text, font, clip)
}
func (c *CmdBuffer) CreateText(x, y float32, txt *widgets.Text, scale float32, font fonts.Font) ([]float32, []int32, int, int) {
	vert, ind, count := c.text(txt, font, x, c.displaySize.Y-y, scale, txt.CurrentColor)
	return vert, ind, count, c.LastInd
}

func (c *CmdBuffer) CheckIndicesChange(w widgets.Widget) bool {
	vert, _, _, l := w.RenderInfo()
	return l != (c.LastInd + len(vert)/9)
}

func (c *CmdBuffer) CreateTexturedRect(x, y, w, h float32, texId uint32, coords, clr [4]float32) ([]float32, []int32, int, int) {
	vert, ind, count := c.rectangle(x, c.displaySize.Y-y, w, h, texId, coords, clr)
	return vert, ind, count, c.LastInd
}

func (c *CmdBuffer) CreateRoundedRect(x, y, w, h float32, radius int, shape RoundedRectShape, clr [4]float32) ([]float32, []int32, int, int) {
	info := c.roundedRectangle(x, c.displaySize.Y-y, w, h, radius, shape, clr)
	return info.Vertices, info.Indices, info.VertCount, c.LastInd
}

//func (c *CmdBuffer) DrawRect(x, y, w, h float32, clr [4]float32) {
//	vert, ind, count := c.rectangle(x, c.displaySize.Y-y, w, h, 0, emptyCoords, clr)
//	c.SendToBuffer(vert, ind, count, c.Indices)
//}

func (c *CmdBuffer) CreateRect(x, y, w, h float32, clr [4]float32) ([]float32, []int32, int, int) {
	vert, ind, count := c.rectangle(x, c.displaySize.Y-y, w, h, 0, emptyCoords, clr)
	return vert, ind, count, c.LastInd
}
func checkSliceForNull(s [4]float32) bool {
	return (s[0] == 0) && (s[1] == 0) && (s[2] == 0) && (s[3] == 0)
}

//func (c *CmdBuffer) AddCommand(cmd Command, clip ClipRectCompose) {
//	c.commands = append(c.commands, cmd)
//
//	switch cmd.Type {
//	case SeparateBuffer:
//		mainRect := clip.MainClipRect
//		innerRect := clip.ClipRect
//
//		x, x2 := int32(mainRect[0]), int32(innerRect[0])
//		y, y2 := int32(mainRect[1]), int32(innerRect[1])
//		w, w2 := int32(mainRect[2]), int32(innerRect[2])
//		h, h2 := int32(mainRect[3]), int32(innerRect[3])
//
//		useInnerClip := !checkSliceForNull(innerRect)
//		xl := x+w < x2+w2
//		yl := y+h < y2+h2
//
//		overlapWidth := useInnerClip && xl
//		overlapHeigth := useInnerClip && yl
//
//		inf := DrawCall{
//			Elems:       c.VertCount - c.lastElems,
//			IndexOffset: c.ofs,
//			TexId:       cmd.sb.texid,
//		}
//		if !useInnerClip {
//			inf.ClipRect = cmd.sb.mainClipRect
//		} else if overlapWidth && overlapHeigth {
//			inf.ClipRect = cmd.sb.mainClipRect
//		} else if overlapWidth {
//			inf.ClipRect = cmd.sb.mainClipRect
//		} else if overlapHeigth {
//			var tmp = cmd.sb.clipRect
//			inf.ClipRect = [4]float32{tmp[0], tmp[1], tmp[2], cmd.sb.mainClipRect[3] - (tmp[1] - cmd.sb.mainClipRect[1])}
//		} else {
//			inf.ClipRect = cmd.sb.clipRect
//		}
//
//		c.DrawCalls = append(c.DrawCalls, inf)
//		c.ofs += c.VertCount - c.lastElems
//		c.lastElems = c.VertCount
//	case RectType:
//		r := cmd.Rect
//		if r.radius == 0 {
//			if r.TexId == 0 {
//				c.RectangleR(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.Clr)
//			} else {
//				c.RectangleT(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.TexId, r.coords, r.Clr)
//				c.SeparateBuffer(r.TexId, clip) // don't forget to slice buffer
//			}
//		} else {
//			if r.TexId == 0 {
//				c.roundedRectangle(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.radius, r.shape, r.Clr)
//			} else {
//				// TODO: Add textured rounded rect
//			}
//		}
//	case Text:
//		t := cmd.Text
//		c.Text(t.Widget, t.Font, t.X, c.displaySize.Y-(t.Y+float32(t.Padding)), t.Scale, t.Clr)
//		c.SeparateBuffer(t.Font.TextureId, clip) // don't forget to slice buffer
//	case BezierQuad:
//		b := cmd.Bezier
//		c.bezierQuad(b.StartX, b.StartY, b.SupportX, b.SupportY, b.EndX, b.EndY, b.Steps, b.Clr, clip)
//		c.sepBuf(clip, "LINE_STRIP")
//	case Line:
//		l := cmd.Line
//		c.line(l.StartX, c.displaySize.Y-l.StartY, l.EndX, c.displaySize.Y-l.EndY, l.Clr)
//		c.sepBuf(clip, "LINE")
//	case LineStrip:
//		l := cmd.Line
//		changed := make([]utils.Vec2, len(l.Points))
//		for i, p := range l.Points {
//			changed[i].Y = c.displaySize.Y - p.Y
//			changed[i].X = p.X
//		}
//		c.lineStrip(l.Clr, changed)
//		c.sepBuf(clip, "LINE_STRIP")
//	}
//}
func (c *CmdBuffer) sepBuf(clip ClipRectCompose, t string) {
	inf := DrawCall{
		Elems:       c.VertCount - c.lastElems,
		IndexOffset: c.ofs,
		TexId:       0,
		Type:        t,
	}
	inf.ClipRect = clip.MainClipRect
	c.DrawCalls = append(c.DrawCalls, inf)
	c.ofs += c.VertCount - c.lastElems
	c.lastElems = c.VertCount
}

func (c *CmdBuffer) SendToBuffer(vert []float32, indices []int32, vertCount int, ind int) {
	c.Vertices = append(c.Vertices, vert...)
	c.Indices = append(c.Indices, indices...)
	c.VertCount += vertCount
	if ind != 0 {
		c.LastInd = ind
	}
}

func (c *CmdBuffer) CreateBorderBox(x, y, w, h, lineWidth float32, clr [4]float32) {
	c.CreateRect(x, y, w, lineWidth, clr)
	c.CreateRect(x+w-lineWidth, y, lineWidth, h, clr)
	c.CreateRect(x, y, lineWidth, h, clr)
	c.CreateRect(x, y+h-lineWidth, w, lineWidth, clr)
}

func (c *CmdBuffer) text(text *widgets.Text, font fonts.Font, x, y float32, scale float32, clr [4]float32) (vert []float32, ind []int32, cnt int) {
	texId := font.TextureId
	for i, r := range []rune(text.Message) {
		info := font.GetCharacter(r)

		if info.Rune == rune(127) { // '\n'
			continue
		}
		xPos := x + text.Chars[i].Pos.X
		yPos := y - text.Chars[i].Pos.Y
		v, idec, vc := c.addCharacter(xPos, yPos, scale, texId, *info, clr)
		vert = append(vert, v...)
		ind = append(ind, idec...)
		cnt += vc
	}
	return vert, ind, cnt
}
func (c *CmdBuffer) addCharacter(x, y float32, scale float32, texId uint32, info fonts.CharInfo, clr [4]float32) ([]float32, []int32, int) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	x0 := x
	y0 := y
	x1 := x + scale*float32(info.Width)
	y1 := y + scale*float32(info.Height)

	ux0, uy0 := info.TexCoords[0].X, info.TexCoords[0].Y
	ux1, uy1 := info.TexCoords[1].X, info.TexCoords[1].Y

	ind0 := c.LastInd
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

	c.LastInd = last + 1
	//c.Render(vert, ind, 6)
	return vert, ind, 6
}

func (c *CmdBuffer) rectangle(x, y, w, h float32, texId uint32, coords [4]float32, clr [4]float32) ([]float32, []int32, int) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	var ux0, uy0, ux1, uy1 float32
	ux0, uy0 = coords[2], coords[3]
	ux1, uy1 = coords[0], coords[1]

	ind0 := c.LastInd
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

	c.LastInd = last + 1
	//fmt.Println(ind)
	//c.Render(vert, ind, 6)
	return vert, ind, 6
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

func (c *CmdBuffer) Arc(x, y, radius float32, steps int, sector CircleSector, clr [4]float32) ([]float32, []int32, int) {
	ind0 := c.LastInd
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
	}
	fillVertices(vert, &offset, lastX, lastY, 0, 0, 0, clr)

	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)

	c.LastInd = ind2 + 1

	return vert, ind, (numV + 1) * 3
}
func (c *CmdBuffer) lineArc(x, y, radius, steps float32, sector CircleSector, clr [4]float32) (vert []float32, ind []int32, cnt int) {
	switch sector {
	case BotLeft:
		vert, ind, cnt = c.bezierQuad(x, y+radius, x-radius, y+radius, x-radius, y, steps, clr)
	case BotRight:
		vert, ind, cnt = c.bezierQuad(x+radius, y, x+radius, y+radius, x, y+radius, steps, clr)
	case TopLeft:
		vert, ind, cnt = c.bezierQuad(x, y-radius, x-radius, y-radius, x-radius, y, steps, clr)
	case TopRight:
		vert, ind, cnt = c.bezierQuad(x, y-radius, x+radius, y-radius, x+radius, y, steps, clr)
	}
	return
}

func (c *CmdBuffer) CreateLineStrip(p []utils.Vec2, clr [4]float32) ([]float32, []int32, int, int) {
	vert, ind, count := c.lineStrip(p, clr)
	return vert, ind, count, c.LastInd
}

func (c *CmdBuffer) CreateLine(startX, startY, endX, endY float32, clr [4]float32) ([]float32, []int32, int, int) {
	vert, ind, count := c.line(startX, startY, endX, endY, clr)
	return vert, ind, count, c.LastInd
}
func (c *CmdBuffer) CreateBezierQuad(startX, startY, supportX, supportY, endX, endY, steps float32, clr [4]float32) ([]float32, []int32, int, int) {
	vert, ind, count := c.bezierQuad(startX, startY, supportX, supportY, endX, endY, steps, clr)
	return vert, ind, count, c.LastInd
}

// TODO: line drawing needs an optimization because now, each line takes 1 draw call. Maybe one buffer that will hold all lines?
func (c *CmdBuffer) line(startX, startY, endX, endY float32, clr [4]float32) ([]float32, []int32, int) {
	ind0 := c.LastInd
	offset := 0
	ind := make([]int32, 2)      // 1 - point
	vert := make([]float32, 9*2) //polygon

	fillVertices(vert, &offset, startX, startY, 0, 0, 0, clr)
	ind[0] = int32(ind0)
	ind0++
	fillVertices(vert, &offset, endX, endY, 0, 0, 0, clr)
	ind[1] = int32(ind0)

	c.LastInd = ind0 + 1
	return vert, ind, 2
}
func (c *CmdBuffer) lineStrip(points []utils.Vec2, clr [4]float32) ([]float32, []int32, int) {
	ind0 := c.LastInd
	offset := 0
	pointsLen := len(points)
	ind := make([]int32, pointsLen)      // 1 - point
	vert := make([]float32, 9*pointsLen) //polygon
	for i, point := range points {
		fillVertices(vert, &offset, point.X, point.Y, 0, 0, 0, clr)
		ind[i] = int32(ind0)
		ind0++
	}
	c.LastInd = ind0
	return vert, ind, pointsLen
}

func (c *CmdBuffer) bezierQuad(startX, startY, supportX, supportY, endX, endY, steps float32, clr [4]float32) ([]float32, []int32, int) {
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
	return c.lineStrip(points, clr)
}

var steps = 30

func (c *CmdBuffer) CreateRoundedBorderRectangle(x, y, w, h float32, radius int, clr [4]float32) ([]float32, []int32, int, int) {
	info := c.roundedLineRectangle(x, y, w, h, 10, radius, AllRounded, clr)
	return info.Vertices, info.Indices, info.VertCount, c.LastInd
}
func (c *CmdBuffer) roundedLineRectangle(x, y, w, h, steps float32, radius int, shape RoundedRectShape, clr [4]float32) (info primitiveInfo) {

	topLeft := utils.Vec2{x + float32(radius), y + float32(radius)} //origin of arc
	topRight := utils.Vec2{x + w - float32(radius), y + float32(radius)}
	botLeft := utils.Vec2{x + float32(radius), y + h - float32(radius)}
	botRight := utils.Vec2{x + w - float32(radius), y + h - float32(radius)}

	switch shape {
	case TopLeftRect:
		vert, ind, cnt := c.lineArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineStrip([]utils.Vec2{
			{topLeft.X, y},
			{x + w, y},
			{x + w, y + h},
			{x, y + h},
			{x, topLeft.Y},
		}, clr)
		info.update(vert, ind, cnt)
	case TopRightRect:
		vert, ind, cnt := c.lineArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineStrip([]utils.Vec2{
			{topRight.X + float32(radius), y + float32(radius)},
			{topRight.X + float32(radius), topRight.Y + h - float32(radius)},
			{x, y + h},
			{x, y},
			{topRight.X, y},
		}, clr)
		info.update(vert, ind, cnt)
	case BotLeftRect:
		vert, ind, cnt := c.lineArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineStrip([]utils.Vec2{
			{x, botLeft.Y},
			{x, y},
			{x + w, y},
			{x + w, y + h},
			{botLeft.X, botLeft.Y + float32(radius)},
		}, clr)
		info.update(vert, ind, cnt)
	case BotRightRect:
		vert, ind, cnt := c.lineArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineStrip([]utils.Vec2{
			{botRight.X, botRight.Y + float32(radius)},
			{x, y + h},
			{x, y},
			{x + w, y},
			{x + w, y + h - float32(radius)},
		}, clr)
		info.update(vert, ind, cnt)
	case TopRect:
		vert, ind, cnt := c.lineArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.line(topLeft.X, y, topRight.X, y, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineStrip([]utils.Vec2{
			{x + w, topRight.Y},
			{x + w, y + h},
			{x, y + h},
			{x, y + float32(radius)},
		}, clr)
		info.update(vert, ind, cnt)
	case BotRect:
		vert, ind, cnt := c.lineArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.line(botRight.X, y+h, botLeft.X, y+h, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineStrip([]utils.Vec2{
			{x, botLeft.Y},
			{x, y},
			{x + w, y},
			{x + w, botRight.Y},
		}, clr)
		info.update(vert, ind, cnt)
	case AllRounded:
		vert, ind, cnt := c.lineArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.line(topLeft.X, y, topRight.X, y, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.line(x+w, topRight.Y, x+w, botRight.Y, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.line(botRight.X, y+h, botLeft.X, y+h, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.lineArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		info.update(vert, ind, cnt)
		vert, ind, cnt = c.line(x, botLeft.Y, x, topLeft.Y, clr)
		info.update(vert, ind, cnt)
	}
	return
}

var emptyCoords = [4]float32{0, 0, 0, 0}

// TODO(@Dmitry-dms): Is there is a better way to handle this?
// primitiveInfo was created as a helper. It helps to prevent boilerplate append().
type primitiveInfo struct {
	Vertices  []float32
	Indices   []int32
	VertCount int
}

func (i *primitiveInfo) update(v []float32, ind []int32, cnt int) {
	i.Vertices = append(i.Vertices, v...)
	i.Indices = append(i.Indices, ind...)
	i.VertCount += cnt
}

func (c *CmdBuffer) roundedRectangle(x, y, w, h float32, radius int, shape RoundedRectShape, clr [4]float32) (info primitiveInfo) {
	topLeft := utils.Vec2{X: x + float32(radius), Y: y - float32(radius)} //origin of arc
	topRight := utils.Vec2{X: x + w - float32(radius), Y: y - float32(radius)}
	botLeft := utils.Vec2{X: x + float32(radius), Y: y - h + float32(radius)}
	botRight := utils.Vec2{X: x + w - float32(radius), Y: y - h + float32(radius)}
	switch shape {
	case TopLeftRect:
		v, i, cnt := c.Arc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, y-float32(radius), w, h-float32(radius), 0, emptyCoords, clr) //main rect
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x+float32(radius), y, w-float32(radius), float32(radius), 0, emptyCoords, clr) //top rect
		info.update(v, i, cnt)
	case TopRightRect:
		v, i, cnt := c.Arc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, y-float32(radius), w, h-float32(radius), 0, emptyCoords, clr) //main
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, y, w-float32(radius), float32(radius), 0, emptyCoords, clr)
		info.update(v, i, cnt)
	case BotLeftRect:
		v, i, cnt := c.Arc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, y, w, h-float32(radius), 0, emptyCoords, clr) //main
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(botLeft.X, botLeft.Y, w-float32(radius), float32(radius), 0, emptyCoords, clr)
		info.update(v, i, cnt)
	case BotRightRect:
		v, i, cnt := c.Arc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, y, w, h-float32(radius), 0, emptyCoords, clr) //main
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, botLeft.Y, w-float32(radius), float32(radius), 0, emptyCoords, clr)
		info.update(v, i, cnt)
	case TopRect:
		v, i, cnt := c.Arc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.Arc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, y-float32(radius), w, h-float32(radius), 0, emptyCoords, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x+float32(radius), y, w-float32(radius)*2, float32(radius), 0, emptyCoords, clr)
		info.update(v, i, cnt)
	case BotRect:
		v, i, cnt := c.Arc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.Arc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, y, w, h-float32(radius), 0, emptyCoords, clr) //main
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(botLeft.X, botLeft.Y, w-float32(radius)*2, float32(radius), 0, emptyCoords, clr)
		info.update(v, i, cnt)
	case AllRounded:
		v, i, cnt := c.Arc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.Arc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.Arc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		info.update(v, i, cnt)
		v, i, cnt = c.Arc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		info.update(v, i, cnt)

		v, i, cnt = c.rectangle(topLeft.X, topLeft.Y+float32(radius), w-float32(radius)*2, float32(radius), 0, emptyCoords, clr) //top
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(x, topLeft.Y, w, h-float32(radius)*2, 0, emptyCoords, clr) //center
		info.update(v, i, cnt)
		v, i, cnt = c.rectangle(botLeft.X, botLeft.Y, w-float32(radius)*2, float32(radius), 0, emptyCoords, clr) //bottom
		info.update(v, i, cnt)
	}
	return
}
