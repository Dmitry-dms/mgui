package fonts

import (
	"github.com/Dmitry-dms/mgui/utils"
)

const (
	CharNotFound = -25
)

type CharInfo struct {
	SrcX, SrcY                int
	Width, Height, Advance    int
	TexCoords                 [2]utils.Vec2
	Ascend, Descend           int
	LeftBearing, RightBearing int
	Rune                      rune
}

func (c *CharInfo) calcTexCoords(fontWidth, fontHeight int) {
	x0 := float32(c.SrcX) / float32(fontWidth)
	x1 := (float32(c.SrcX) + float32(c.Width)) / float32(fontWidth)
	y0 := float32(c.SrcY) / float32(fontHeight)
	y1 := (float32(c.SrcY) - float32(c.Height)) / float32(fontHeight)

	c.TexCoords[0] = utils.Vec2{X: x0, Y: y0}
	c.TexCoords[1] = utils.Vec2{X: x1, Y: y1}
}

func (f *Font) GetCharacter(chr rune) *CharInfo {
	c, ok := f.CharMap[int(chr)]
	if !ok {
		return f.CharMap[CharNotFound]
	}
	return c
}
