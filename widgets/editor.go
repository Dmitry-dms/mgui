package widgets

import (
	"github.com/Dmitry-dms/mgui/fonts"
	"github.com/Dmitry-dms/mgui/utils"
	"strings"
	"unsafe"
)

// https://www.youtube.com/watch?v=w_yXlnjeAy4
type Editor struct {
	CharsInfo                                    []Character
	Buff                                         *Buffer
	Lines                                        []Line
	linesCount                                   int
	cursor                                       int
	font                                         *fonts.Font
	Scale                                        float32
	TextWidth, TextHeight                        float32
	Selection                                    bool
	SelectionEnd, SelectionBegin, SelectionPoint int
}

type Buffer struct {
	Data []rune
}

type Character struct {
	Xpos, Ypos    float32
	absoluteYpos  float32 //uses in case of selection
	Width, Height float32
	Info          *fonts.CharInfo
}

func newBuffer(initCap int) *Buffer {
	b := Buffer{Data: make([]rune, 0, initCap)}
	return &b
}

func (b *Buffer) replace(msg string) {
	b.Data = b.Data[:0]
	b.Data = append(b.Data, []rune(msg)...)
}

func (b *Buffer) grow() {
	buf := make([]rune, len(b.Data), 2*cap(b.Data))
	copy(buf, b.Data)
	b.Data = buf
}

func (e *Editor) String() string {
	return *(*string)(unsafe.Pointer(&e.Buff.Data))
}
func (b *Buffer) insertIntoBuffer(txt string, cursor int) {

	if cap(b.Data) <= len(b.Data)+len(txt) {
		b.grow()
	}
	b.Data = append(b.Data[:cursor], append([]rune(txt), b.Data[cursor:]...)...)
}
func (b *Buffer) Cap() int {
	return cap(b.Data)
}
func (b *Buffer) Len() int {
	return len(b.Data)
}

func NewEditor(f *fonts.Font, scale float32) *Editor {
	e := Editor{
		Buff:      newBuffer(10),
		Lines:     make([]Line, 1),
		CharsInfo: make([]Character, 0, 10),
		font:      f,
		Scale:     scale,
	}

	return &e
}

type Line struct {
	Begin, End    int
	Xpos, Ypos    float32
	Width, Height float32
}

var sb strings.Builder

func (e *Editor) StartSelection(i int) {
	e.SelectionBegin = i
}

func (e *Editor) GetTextSelection(x, y float32) (text string, selectedRegions []utils.Rect) {

	//if e.SelectionEnd < e.SelectionBegin {
	//	tmp := e.SelectionBegin
	//	e.SelectionBegin = e.SelectionEnd
	//	e.SelectionEnd = tmp
	//}
	//sb.Grow(e.SelectionEnd - e.SelectionBegin)
	var yPosChecker float32

	var (
		begin int
		end   int
	)
	if e.SelectionEnd < e.SelectionBegin {
		begin = e.SelectionEnd
		end = e.SelectionPoint
	} else {
		begin = e.SelectionBegin
		end = e.SelectionEnd
	}

	var selectedLine utils.Rect
	firstChar := e.CharsInfo[begin]
	selectedLine.Min = utils.Vec2{x + firstChar.Xpos, y + firstChar.absoluteYpos}
	selectedLine.Max.Y = e.Lines[0].Height * e.Scale
	//fmt.Println(begin, end)

	for i := begin; i <= end; i++ {
		char := e.CharsInfo[i]
		sb.WriteRune(char.Info.Rune)
		if yPosChecker != char.absoluteYpos {
			selectedLine.Max.X -= char.Width
			selectedRegions = append(selectedRegions, selectedLine)
			selectedLine.Min = utils.Vec2{x + char.Xpos, y + char.absoluteYpos}
			selectedLine.Max.X = 0
			yPosChecker = char.absoluteYpos
		}

		selectedLine.Max.X += char.Width
	}
	selectedRegions = append(selectedRegions, selectedLine)
	text = sb.String()
	sb.Reset()

	return
}

func (e *Editor) Backspace() {

	if e.cursor > e.Buff.Len() {
		e.cursor = e.Buff.Len()
	}

	if e.cursor == 0 {
		return
	}

	e.Buff.Data = append(e.Buff.Data[:e.cursor-1], e.Buff.Data[e.cursor:]...)
	e.cursor--
	e.retokenize()
}

func (e *Editor) ReplaceBuffer(msg string) (width, height float32) {
	e.Buff.replace(msg)
	e.retokenize()
	width = e.TextWidth
	height = e.TextHeight
	return
}

func (e *Editor) Delete() {

	if e.cursor >= e.Buff.Len() {
		return
	}

	e.Buff.Data = append(e.Buff.Data[:e.cursor], e.Buff.Data[e.cursor+1:]...)

	e.retokenize()
}
func (e *Editor) MoveCharLeft() {
	if e.cursor > 0 {
		e.cursor--
	}
}
func (e *Editor) MoveCharRight() {
	if e.cursor < e.Buff.Len() {
		e.cursor++
	}
}

func (e *Editor) insertIntoBuffer(txt string) {
	e.Buff.insertIntoBuffer(txt, e.cursor)
}
func (e *Editor) InsertText(txt string) {
	if e.cursor > e.Buff.Len() {
		e.cursor = e.Buff.Len()
	}

	e.insertIntoBuffer(txt)
	e.cursor++

	e.retokenize()
}

func (e *Editor) retokenize() {
	e.linesCount = 1

	fontSize := e.font.XHeight() * e.Scale
	var line Line
	line.Begin = 0
	lines := e.Lines[:0]
	var boundWidth, boundHeight float32
	var maxDescend, baseline, maxWidth float32
	baseline = fontSize
	boundHeight = fontSize
	var dx, ypos float32 = 0, 0
	{
		line.Xpos = 0
		line.Ypos = ypos
		line.Height = fontSize
	}
	chars := e.CharsInfo[:0]
	var char Character

	srcString := e.Buff.Data

	prevR := rune(-1)
	for i, r := range srcString {
		var charWidth float32 = 0
		if r == '\n' {
			e.linesCount++
			line.End = i
			line.Width = dx
			lines = append(lines, line)
			dx = 0
			line.Xpos = 0
			line.Ypos = baseline
			baseline += fontSize * e.Scale
			line.Width = 0
			boundHeight += fontSize
			line.Begin = i + 1
			chars = append(chars, char)

			prevR = rune(-1)
			continue
		}

		if prevR >= 0 {
			kern := e.font.Face.Kern(prevR, r).Ceil()
			dx += float32(kern)
			charWidth += float32(kern)
		}

		char.Info = e.font.GetCharacter(r)

		//if r != ' ' {
		dx += float32(char.Info.LeftBearing)
		charWidth += float32(char.Info.LeftBearing)
		//}

		yPos := baseline
		xPos := dx
		if char.Info.Descend != 0 {
			d := float32(char.Info.Descend) * e.Scale
			yPos += d
			if d > maxDescend {
				maxDescend = d
			}
		}

		char.Xpos = xPos
		char.Ypos = yPos
		char.absoluteYpos = line.Ypos
		char.Width = charWidth + float32(char.Info.Width)*e.Scale + float32(char.Info.RightBearing)

		chars = append(chars, char)
		dx += float32(char.Info.Width) * e.Scale
		//if r != ' ' {
		dx += float32(char.Info.RightBearing)
		//}

		prevR = r
		boundWidth = dx
		line.Width = boundWidth
		if e.linesCount > 1 {
			if boundWidth > maxWidth {
				maxWidth = boundWidth
			}
		} else {
			maxWidth = boundWidth
		}

	}
	line.End = e.Buff.Len()
	lines = append(lines, line)
	e.Lines = lines
	//boundHeight += maxDescend
	e.TextWidth = maxWidth
	e.TextHeight = boundHeight + maxDescend
	e.CharsInfo = chars
}

func (e *Editor) cursorRow() int {
	for row := 0; row < e.linesCount; row++ {
		line := e.Lines[row]

		if line.Begin <= e.cursor && e.cursor <= line.End {
			return row
		}
	}
	return e.linesCount - 1
}
func (e *Editor) MoveLineUp() {
	cursorRow := e.cursorRow()
	cursorCol := e.cursor - e.Lines[cursorRow].Begin

	if cursorRow > 0 {
		nextLine := e.Lines[cursorRow-1]
		nextLineSize := nextLine.End - nextLine.Begin

		if cursorCol > nextLineSize {
			cursorCol = nextLineSize
		}
		e.cursor = nextLine.Begin + cursorCol
	}
}
func (e *Editor) MoveLineDown() {
	cursorRow := e.cursorRow()
	cursorCol := e.cursor - e.Lines[cursorRow].Begin

	if cursorRow < e.linesCount-1 {
		nextLine := e.Lines[cursorRow+1]
		nextLineSize := nextLine.End - nextLine.Begin

		if cursorCol > nextLineSize {
			cursorCol = nextLineSize
		}
		e.cursor = nextLine.Begin + cursorCol
	}
}
