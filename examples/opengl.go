////go:build opengl

package main

import (
	"errors"
	"fmt"
	"github.com/Dmitry-dms/mgui/draw"
	"github.com/go-gl/gl/v4.2-core/gl"
	"image"
	"image/color"
	idraw "image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

type GLRender struct {
	vaoId, vboId, ebo uint32
	shaderProgram     *ShaderProgram
}

const (
	// pos     color       texCoords    texId
	// f,f     f,f,f,f     f,f          f
	posSize       = 2
	colorSize     = 4
	texCoordsSize = 2
	texIdSize     = 1

	vertexSize = posSize + colorSize + texCoordsSize + texIdSize

	posOffset       = 0
	colorOffset     = posOffset + posSize
	texCoordsOffset = colorOffset + colorSize
	texIdOffset     = texCoordsOffset + texCoordsSize
)

func NewGlRenderer() *GLRender {
	s, err := NewShaderProgram("examples/gui.glsl")
	if err != nil {
		panic(err)
	}
	r := GLRender{
		vaoId:         0,
		vboId:         0,
		ebo:           0,
		shaderProgram: s,
	}

	gl.GenBuffers(1, &r.vboId)
	gl.GenBuffers(1, &r.ebo)

	//включаем layout
	// gogl.SetVertexAttribPointer(0, posSize, gl.FLOAT, vertexSize*4, posOffset)
	// gogl.SetVertexAttribPointer(1, colorSize, gl.FLOAT, vertexSize*4, colorOffset)
	// gogl.SetVertexAttribPointer(2, texCoordsSize, gl.FLOAT, vertexSize*4, texCoordsOffset)
	// gogl.SetVertexAttribPointer(3, texIdSize, gl.FLOAT, vertexSize*4, texIdOffset)
	return &r
}

func (r *GLRender) NewFrame() {

}

func (b *GLRender) Draw(displaySize [2]float32, buffer draw.CmdBuffer) {

	displayWidth := displaySize[0]
	displayHeight := displaySize[1]

	b.shaderProgram.Use()
	vaoId := GenBindVAO()
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboId)

	SetVertexAttribPointer(0, posSize, gl.FLOAT, vertexSize*4, posOffset)
	SetVertexAttribPointer(1, colorSize, gl.FLOAT, vertexSize*4, colorOffset)
	SetVertexAttribPointer(2, texCoordsSize, gl.FLOAT, vertexSize*4, texCoordsOffset)
	SetVertexAttribPointer(3, texIdSize, gl.FLOAT, vertexSize*4, texIdOffset)

	gl.Enable(gl.SCISSOR_TEST)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboId)
	gl.BufferData(gl.ARRAY_BUFFER, len(buffer.Vertices)*4, gl.Ptr(buffer.Vertices), gl.STREAM_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.Indices)*4, gl.Ptr(buffer.Indices), gl.STREAM_DRAW)

	orthoProjection := [4][4]float32{
		{2.0 / displayWidth, 0.0, 0.0, 0.0},
		{0.0, 2.0 / displayHeight, 0.0, 0.0},
		{0.0, 0.0, -2.0, 0.0},
		{-1.0, -1.0, -1.0, 1.0},
	}

	b.shaderProgram.UploadMatslice("uProjection", orthoProjection)

	for _, cmd := range buffer.DrawCalls {
		//fmt.Println(cmd.TexId)
		//mainRect := cmd.Clip.MainClipRect
		clipRect := cmd.ClipRect

		x := int32(clipRect[0])
		y := int32(clipRect[1])
		w := int32(clipRect[2])
		h := int32(clipRect[3])

		y = int32(displayHeight) - (y + h)
		//_ = x
		//_ = w
		// fmt.Printf("type = %s, elems = %d, ofs = %d, texId = %d \n", cmd.Type, cmd.Elems, cmd.IndexOffset, cmd.TexId)
		if cmd.TexId != 0 {
			gl.ActiveTexture(gl.TEXTURE0 + cmd.TexId)
			gl.BindTexture(gl.TEXTURE_2D, cmd.TexId)
			b.shaderProgram.UploadTexture("Texture", int32(cmd.TexId))
			if cmd.Type == "msdf" {
				b.shaderProgram.UploadInt("textureType", 1)
			} else {
				b.shaderProgram.UploadInt("textureType", 0)
			}
		}
		gl.Scissor(x, y, w, h)
		if cmd.Type == "LINE_STRIP" {
			gl.LineWidth(3)
			gl.DrawElementsBaseVertexWithOffset(gl.LINE_STRIP, int32(cmd.Elems), gl.UNSIGNED_INT,
				uintptr(cmd.IndexOffset*4), 0)
		} else if cmd.Type == "LINE" {
			gl.LineWidth(3)
			gl.DrawElementsBaseVertexWithOffset(gl.LINES, int32(cmd.Elems), gl.UNSIGNED_INT,
				uintptr(cmd.IndexOffset*4), 0)
		} else {
			gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(cmd.Elems), gl.UNSIGNED_INT,
				uintptr(cmd.IndexOffset*4), 0)
		}

	}

	b.shaderProgram.Detach()
	gl.DeleteVertexArrays(1, &vaoId)
	gl.Disable(gl.SCISSOR_TEST)
}

type ShaderProgram struct {
	ProgramId uint32
	path      string
	beingUsed bool
}

func NewShaderProgram(path string) (*ShaderProgram, error) {
	id, err := CreateProgram(path)
	if err != nil {
		return nil, err
	}
	result := ShaderProgram{
		ProgramId: id,
		path:      path,
		beingUsed: false,
	}
	return &result, nil
}

func (s *ShaderProgram) Use() {
	if !s.beingUsed {
		useProgram(s.ProgramId)
		s.beingUsed = true
	}
}
func (s *ShaderProgram) Detach() {
	useProgram(0)
	s.beingUsed = false
}

func (s *ShaderProgram) UploadFloat(name string, f float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	gl.Uniform1f(location, f)
}
func (s *ShaderProgram) UploadTexture(name string, slot int32) {
	name_cstr := gl.Str(name + "\x00")
	// location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	// s.Use()
	gl.Uniform1i(gl.GetUniformLocation(s.ProgramId, name_cstr), slot)
}

func (s *ShaderProgram) UploadInt(name string, slot int32) {
	name_cstr := gl.Str(name + "\x00")
	// location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	// s.Use()
	gl.Uniform1i(gl.GetUniformLocation(s.ProgramId, name_cstr), slot)
}

func (s *ShaderProgram) UploadVec2(name string, vec []float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	gl.Uniform2f(location, vec[0], vec[1])
}

//func (s *ShaderProgram) UploadVec3(name string, vec mgl32.Vec3) {
//	name_cstr := gl.Str(name + "\x00")
//	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
//	s.Use()
//	v3 := [3]float32(vec)
//	gl.Uniform3fv(location, 1, &v3[0])
//}
//func (s *ShaderProgram) UploadVec4(name string, vec mgl32.Vec4) {
//	name_cstr := gl.Str(name + "\x00")
//	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
//	s.Use()
//	v4 := [4]float32(vec)
//	gl.Uniform4fv(location, 1, &v4[0])
//}
func (s *ShaderProgram) UploadMatslice(name string, mat [4][4]float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	// m4 := [16]float32(mat)
	gl.UniformMatrix4fv(location, 1, false, &mat[0][0])
}

//func (s *ShaderProgram) UploadMat4(name string, mat mgl32.Mat4) {
//	name_cstr := gl.Str(name + "\x00")
//	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
//	s.Use()
//	m4 := [16]float32(mat)
//	gl.UniformMatrix4fv(location, 1, false, &m4[0])
//}
//
//func (s *ShaderProgram) UploadIntArray(name string, array []int32) {
//	name_cstr := gl.Str(name + "\x00")
//	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
//	s.Use()
//	gl.Uniform1iv(location, int32(len(array)), &array[0])
//}

func LoadShaders(path string) (uint32, uint32, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}
	shaders := string(file)
	spl := strings.Split(shaders, "#type")
	var vertexSource, fragmentSource string

	for i := 1; i < len(spl); i++ {
		tmp := strings.Split(spl[i], "\r\n")
		shaderTypeStr := strings.TrimSpace(tmp[0])
		switch shaderTypeStr {
		case "vertex":
			vertexSource = spl[i][len(tmp[0]):]
		case "fragment":
			fragmentSource = spl[i][len(tmp[0]):]
		}
	}
	vertId, err := CreateShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, 0, err
	}
	fragId, err := CreateShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, 0, err
	}
	return vertId, fragId, nil
}

func CreateShader(source string, shaderType uint32) (uint32, error) {
	shaderId := gl.CreateShader(shaderType)
	vsource, free := gl.Strs(source, "\x00")
	gl.ShaderSource(shaderId, 1, vsource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status) //logging
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		return 0, errors.New(log)
	}
	return shaderId, nil
}

func Str(src string) *uint8 {
	return gl.Str(src + "\x00")
}

func CreateProgram(path string) (uint32, error) {

	vert, frag, err := LoadShaders(path)
	if err != nil {
		return 0, err
	}
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vert)
	gl.AttachShader(shaderProgram, frag)
	gl.LinkProgram(shaderProgram)
	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status) //logging
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program: %s \n", log)
	}

	gl.DeleteShader(vert)
	gl.DeleteShader(frag)

	return shaderProgram, nil
}

func GenBindVAO() uint32 {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return VAO
}

func useProgram(progId uint32) {
	gl.UseProgram(progId)
}

type Texture struct {
	Filepath  string `json:"filepath"`
	TextureId uint32 `json:"texture_id"`
	Width     int32  `json:"texture_width"`
	Height    int32  `json:"texture_height"`
	TexCoords [4]float32
}

func (t *Texture) GetFilepath() string {
	return t.Filepath
}

func genBindTexture() uint32 {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	return texId
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.TextureId)
}
func (t *Texture) BindActive(texture uint32) {
	gl.ActiveTexture(texture)
	gl.BindTexture(gl.TEXTURE_2D, t.TextureId)
}
func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) GetWidth() int32 {
	return t.Width
}
func (t *Texture) GetHeight() int32 {
	return t.Height
}

func (t *Texture) GetId() uint32 {
	return t.TextureId
}

// TODO: Replace gl.RGBA with gl.RED (probably requires changing shaderProgram: separate font shaderProgram from general sahder gui.glsl)
func UploadRGBATextureFromMemory(data image.Image) *Texture {
	w := data.Bounds().Max.X
	h := data.Bounds().Max.Y
	pixels := make([]byte, w*h*4)
	bIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			//r, _, _, _ := data.At(x, y).RGBA()
			r, g, b, a := data.At(x, y).RGBA()
			rb := byte(r)
			gb := byte(g)
			bb := byte(b)
			ab := byte(a)
			pixels[bIndex] = rb
			bIndex++
			pixels[bIndex] = gb
			bIndex++
			pixels[bIndex] = bb
			bIndex++
			//if rb == 0 && gb == 0 && bb == 0 {
			if ab == 0 {
				pixels[bIndex] = byte(0)
				//} else if rb <= 150 && gb <= 150 && bb <= 150 { // removes char outlining
				//	pixels[bIndex] = byte(0)
			} else {
				//fmt.Println(r, g, b)
				pixels[bIndex] = ab
			}
			bIndex++
		}
	}
	texture := genBindTexture()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	textureStruct := Texture{
		TextureId: texture,
		Width:     int32(w),
		Height:    int32(h),
	}
	textureStruct.Unbind()
	pixels = nil
	return &textureStruct
}

func TextureFromPNG(filepath string) (*Texture, error) {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		return nil, err
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	texture := genBindTexture()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	//gl.GenerateMipmap(gl.TEXTURE_2D)

	textureStruct := Texture{
		Filepath:  filepath,
		TextureId: texture,
		Width:     int32(w),
		Height:    int32(h),
	}
	textureStruct.Unbind()
	return &textureStruct, nil
}

func (t *Texture) Init(filepath string) (*Texture, error) {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	pixels := make([]byte, w*h*4)
	i := 0
	for y := h - 1; y >= 0; y-- {
		for x := 0; x < w; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			pixels[i] = c.R
			pixels[i+1] = c.G
			pixels[i+2] = c.B
			pixels[i+3] = c.A

			i += 4
		}
	}

	texture := genBindTexture()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	// gl.GenerateMipmap(gl.TEXTURE_2D)

	textureStruct := Texture{
		Filepath:  filepath,
		TextureId: texture,
		Width:     int32(w),
		Height:    int32(h),
	}
	textureStruct.Unbind()
	return &textureStruct, nil
}

func ImageToBytes(img image.Image) []byte {
	size := img.Bounds().Size()
	w, h := size.X, size.Y

	switch img := img.(type) {
	case *image.Paletted:
		bs := make([]byte, 4*w*h)

		b := img.Bounds()
		x0 := b.Min.X
		y0 := b.Min.Y
		x1 := b.Max.X
		y1 := b.Max.Y

		palette := make([]uint8, len(img.Palette)*4)
		for i, c := range img.Palette {
			rgba := color.RGBAModel.Convert(c).(color.RGBA)
			palette[4*i] = rgba.R
			palette[4*i+1] = rgba.G
			palette[4*i+2] = rgba.B
			palette[4*i+3] = rgba.A
		}
		// Even img is a subimage of another image, Pix starts with 0-th index.
		idx0 := 0
		idx1 := 0
		d := img.Stride - (x1 - x0)
		for j := 0; j < y1-y0; j++ {
			for i := 0; i < x1-x0; i++ {
				p := int(img.Pix[idx0])
				bs[idx1] = palette[4*p]
				bs[idx1+1] = palette[4*p+1]
				bs[idx1+2] = palette[4*p+2]
				bs[idx1+3] = palette[4*p+3]
				idx0++
				idx1 += 4
			}
			idx0 += d
		}
		return bs
	case *image.RGBA:
		if len(img.Pix) == 4*w*h {
			return img.Pix
		}
		return imageToBytesSlow(img)
	default:
		return imageToBytesSlow(img)
	}
}

func imageToBytesSlow(img image.Image) []byte {
	size := img.Bounds().Size()
	w, h := size.X, size.Y
	bs := make([]byte, 4*w*h)

	dstImg := &image.RGBA{
		Pix:    bs,
		Stride: 4 * w,
		Rect:   image.Rect(0, 0, w, h),
	}
	idraw.Draw(dstImg, image.Rect(0, 0, w, h), img, img.Bounds().Min, idraw.Src)
	return bs
}

func flipImageY(stride, height int, pixels []byte) {
	// Flip image in y-direction. OpenGL's origin is in the lower
	// left corner.
	row := make([]uint8, stride)
	for y := 0; y < height/2; y++ {
		y1 := height - y - 1
		dest := y1 * stride
		src := y * stride
		copy(row, pixels[dest:])
		copy(pixels[dest:], pixels[src:src+len(row)])
		copy(pixels[src:], row)
	}
}

func SetVertexAttribPointer(index uint32, size int32, xtype uint32, stride, offset int) {
	var memSize int = 0
	switch xtype {
	case gl.INT:
		fallthrough
	case gl.FLOAT:
		memSize = 4
	}
	gl.VertexAttribPointer(index, size, xtype, false, int32(stride), gl.PtrOffset(offset*memSize))
	gl.EnableVertexAttribArray(index)
}
