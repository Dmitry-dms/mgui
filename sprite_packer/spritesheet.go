package sprite_packer

import (
	"encoding/json"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

type SpriteSheet struct {
	Width             int               `json:"atlas_width"`
	Height            int               `json:"atlas_height"`
	Filename          string            `json:"filename"`
	Groups            map[string]*Group `json:"groups"`
	image             *image.RGBA
	SrcX              int `json:"src_x"`
	SrcY              int `json:"src_y"`
	currentHeight     int
	wasResized        bool
	prevWidth         int
	prevHeight        int
	currentGroup      string
	currentSeparators []separator
}

type Group struct {
	Id       string                 `json:"group_id"`
	Contents map[string]*SpriteInfo `json:"sprites"`
}

func NewGroup(id string, initCap int) *Group {
	return &Group{
		Id:       id,
		Contents: make(map[string]*SpriteInfo, initCap),
	}
}

func NewSpriteSheet(initWidth int, filename string) *SpriteSheet {
	s := SpriteSheet{
		Width:    initWidth,
		Height:   initWidth,
		Filename: filename,
		Groups:   make(map[string]*Group),
	}
	im := image.NewRGBA(image.Rect(0, 0, initWidth, initWidth))
	s.image = im
	return &s
}

func (s *SpriteSheet) Image() *image.RGBA {
	return s.image
}

func (s *SpriteSheet) GetGroup(id string) (*Group, bool) {
	g, ok := s.Groups[id]
	return g, ok
}

func (s *SpriteSheet) AddSprite(groupId, spriteId string, data image.Image) *SpriteInfo {
	var spr *SpriteInfo
	s.BeginGroup(groupId, func() map[string]*SpriteInfo {
		spriteInfo := map[string]*SpriteInfo{}
		d := s.GetData(data)
		info := s.AddToSheet(spriteId, d)
		spr = info
		spriteInfo[spriteId] = spr
		return spriteInfo
	})
	return spr
}

func (s *SpriteSheet) increaseImage() int {
	s.prevWidth = s.Width
	s.prevHeight = s.Height
	newW := s.Width * 2
	im := image.NewRGBA(image.Rect(0, 0, newW, newW))
	draw.Draw(im, s.image.Rect, s.image, image.Point{}, draw.Src)
	s.image = im
	s.wasResized = true
	return newW
}
func printBorder(m *image.RGBA, x, y, w int) {

	//for i := y; i >= y-h; i-- {
	//	m.Set(x, i, colornames.Red)
	//}
	//for i := x; i <= x+w; i++ {
	//	m.Set(i, y-h, colornames.Red)
	//}
	//for i := y; i >= y-h; i-- {
	//	m.Set(x+w, i, colornames.Red)
	//}
	for i := x + w; x <= i; i-- {
		m.Set(i, y, colornames.Violet)
	}

}

func (s *SpriteSheet) findEmptySpace(srcX, srcY, width int) int {
	if srcY == 0 {
		return 0
	}
	rayStep := int(math.Ceil(float64(width / 10)))
	if rayStep == 0 {
		rayStep = 1
	}
	tmpH := 100000
	for j := srcX; j < srcX+width; j += rayStep {
		inner := 0
		for i := srcY; i < 0; i++ {
			if s.image.RGBAAt(j, -i).R != 0 {
				inner--
				break
			} else {
				inner++
			}
		}
		if inner < tmpH {
			tmpH = inner
		}
	}
	return tmpH
}

func (s *SpriteSheet) BeginGroup(id string, f func() map[string]*SpriteInfo) {
	s.currentGroup = id
	s.currentSeparators = []separator{}

	gr, exists := s.GetGroup(id)
	if !exists {
		gr = NewGroup(id, 1)
		gr.Contents = f()
		s.Groups[id] = gr
	} else {
		for _, info := range f() {
			gr.Contents[info.Id] = info
		}
	}
	s.cleanSeparators()
	s.currentGroup = ""
	if s.wasResized {
		for _, infos := range s.Groups {
			for _, info := range infos.Contents {
				if info != nil {
					info.calcTexCoords(s.Width, s.Height)
				}
			}
		}
		s.wasResized = false
	}
	s.currentSeparators = []separator{}
}

type separator struct {
	x, y, w int
}

func (s *SpriteSheet) cleanSeparators() {
	for _, i2 := range s.currentSeparators {
		for i := i2.x + i2.w; i2.x <= i; i-- {
			s.image.Set(i, i2.y, color.Transparent)
		}
	}
}

func (s *SpriteSheet) AddToSheet(id string, pixels [][]color.Color) *SpriteInfo {
	height := len(pixels)
	width := len(pixels[0])
	//fmt.Printf("%q x = %d, y = %d \n", char, s.SrcX, s.SrcY)
	//fmt.Println(width+s.SrcX, s.Width)
	if width+s.SrcX > s.Width {
		s.SrcY -= s.currentHeight
		if -s.SrcY+height > s.prevHeight {
			s.SrcX = 0
		} else {
			s.SrcX = s.prevWidth
		}
	}
	if -s.SrcY+height > s.Height {
		h := s.increaseImage()
		s.SrcX = s.prevWidth
		s.SrcY = 0
		s.Width = h
		s.Height = h
	}
	g1 := s.findEmptySpace(s.SrcX, s.SrcY, width)

	if height > s.currentHeight {
		s.currentHeight = height + 1
	}
	ypos := s.SrcY

	if g1 > 0 {
		ypos += g1
	}
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			p := pixels[y][x]
			s.image.Set(x+s.SrcX, y-ypos, p)
		}
	}
	printBorder(s.image, s.SrcX, -ypos+height, width-1)
	s.currentSeparators = append(s.currentSeparators, separator{
		x: s.SrcX,
		y: -ypos + height,
		w: width,
	})
	srcInfo := SpriteInfo{
		Id:     id,
		SrcX:   s.SrcX,
		SrcY:   -ypos + height,
		Width:  width,
		Height: height,
	}
	srcInfo.calcTexCoords(s.Width, s.Height)
	s.SrcX += width + 1
	return &srcInfo
}

func (s *SpriteSheet) GetData(data image.Image) [][]color.Color {
	b := data.Bounds()
	pixels := make([][]color.Color, b.Dy())
	for i := range pixels {
		pixels[i] = make([]color.Color, b.Dx())
	}
	ik := 0 // Needs of using internal counters because [data]
	// is probably subImage that points to the original image
	for i := b.Min.Y; i < b.Max.Y; i++ {
		jk := 0
		for j := b.Min.X; j < b.Max.X; j++ {
			pixels[ik][jk] = data.At(j, i)
			jk++
		}
		ik++
	}
	return pixels
}

func (s *SpriteSheet) SaveSpriteSheetInfo(filename string) error {
	fil, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(fil)
	return enc.Encode(s)
}
func GetSpriteSheetFromFile(filename, imageName string) (*SpriteSheet, error) {
	fil, err := os.OpenFile(filename, os.O_RDONLY, 0664)
	if err != nil {
		return nil, err
	}
	var ss SpriteSheet
	dec := json.NewDecoder(fil)
	err = dec.Decode(&ss)
	if err != nil {
		return nil, err
	}
	ss.image = openImage(imageName)
	return &ss, nil
}
func openImage(filepath string) *image.RGBA {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		return nil
	}

	rgba := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			clr := img.At(x, y)
			rgba.Set(x, y, clr)
		}
	}

	return rgba
}

type SpriteInfo struct {
	Id         string     `json:"sprite_id"`
	SrcX       int        `json:"sprite_src_x"`
	SrcY       int        `json:"sprite_src_y"`
	Width      int        `json:"sprite_w"`
	Height     int        `json:"sprite_h"`
	TextCoords [4]float32 `json:"sprite_tex_coords"` //uv0 uv0 uv1 uv1
}

func (c *SpriteInfo) calcTexCoords(atlasWidth, atlasHeight int) {
	x0 := float32(c.SrcX) / float32(atlasWidth)
	x1 := (float32(c.SrcX) + float32(c.Width)) / float32(atlasWidth)
	y0 := float32(c.SrcY) / float32(atlasHeight)
	y1 := (float32(c.SrcY) - float32(c.Height)) / float32(atlasHeight)

	c.TextCoords = [4]float32{x0, y0, x1, y1}
}
