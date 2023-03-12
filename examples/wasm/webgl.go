package main

import (
	"github.com/Dmitry-dms/mgui/draw"
	"github.com/Dmitry-dms/mgui/examples/wasm/gltypes"
	"io/ioutil"
	"strings"
	"syscall/js"
	"unsafe"
)

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

// TODO:(@Dmitry-dms): Make this more readable, maybe like https://github.com/goxjs/gl?

type WebGLRenderer struct {
	programID, vbo, ebo js.Value
	gl                  js.Value
	renderFrame         js.Func
}

func CreateShader(source string, gl js.Value, typ js.Value) js.Value {
	//js.Global().Get("console").Call("log", source)
	shaderId := gl.Call("createShader", typ)
	//shaderId := gl.CreateShader(typ)
	//gl.ShaderSource(shaderId, source)
	//gl.CompileShader(shaderId)
	gl.Call("shaderSource", shaderId, source)
	gl.Call("compileShader", shaderId)

	js.Global().Get("console").Call("log", gl.Call("getShaderInfoLog", shaderId))
	return shaderId
}

const vertSrcOrig = `# version 300 es
layout (location=0) in vec3 aPos;
layout (location=1) in vec4 aColor;
layout (location=2) in vec2 aTexCoords;
layout (location=3) in float aTexId;

uniform mat4 uProjection;

out vec4 fColor;
out vec2 fTexCoords;
out float fTexId;

void main()
{
    fColor = aColor;
    fTexCoords = aTexCoords;
    fTexId = aTexId;
    gl_Position = uProjection * vec4(aPos,1.0);
}`

const fragSrcOrig = `# version 300 es
precision mediump float;
in vec4 fColor;
in vec2 fTexCoords;
in float fTexId;
out vec4 color;

uniform sampler2D Texture;
uniform int textureType;

float median(float r, float g, float b) {
    return max(min(r, g), min(max(r, g), b));
}

float screenPxRange() {
    return 4.5;
}

void main()
{
	color = fColor;

}`

func LoadShaders(path string, gl js.Value) (js.Value, js.Value) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		js.Global().Get("console").Call("log", "Hello world Go/wasm!")
	}

	shaders := string(file)
	spl := strings.Split(shaders, "#type")
	//var vertexSource, fragmentSource string

	for i := 1; i < len(spl); i++ {
		tmp := strings.Split(spl[i], "\r\n")
		shaderTypeStr := strings.TrimSpace(tmp[0])
		switch shaderTypeStr {
		case "vertex":
			//vertexSource = spl[i][len(tmp[0]):]
		case "fragment":
			//fragmentSource = spl[i][len(tmp[0]):]
		}
	}

	js.Global().Get("console").Call("log", "creating shaders")
	vertId := CreateShader(vertSrcOrig, gl, glTypes.VertexShader)
	fragId := CreateShader(fragSrcOrig, gl, glTypes.FragmentShader)

	return vertId, fragId
}

func CreateProgram(path string, gl js.Value) js.Value {

	vert, frag := LoadShaders(path, gl)
	// Create a shader program object to store
	// the combined shader program
	shaderProgram := gl.Call("createProgram")
	//shaderProgram := gl.CreateProgram()
	//gl.AttachShader(shaderProgram, vert)
	//gl.AttachShader(shaderProgram, frag)
	//gl.LinkProgram(shaderProgram)
	gl.Call("attachShader", shaderProgram, vert)
	gl.Call("attachShader", shaderProgram, frag)
	gl.Call("linkProgram", shaderProgram)

	js.Global().Get("console").Call("log", "program log", gl.Call("getProgramInfoLog", shaderProgram).String())
	//gl.DeleteShader(vert)
	//gl.DeleteShader(frag)

	return shaderProgram
}

func NewWebGLRenderer(gl js.Value) *WebGLRenderer {
	renderer := WebGLRenderer{}
	s := CreateProgram("C:/Users/dmitry/go/src/github.com/Dmitry-dms/mgui/examples/gui.glsl", gl)
	renderer.programID = s
	vertexBuffer := gl.Call("createBuffer")
	indexBuffer := gl.Call("createBuffer")
	//vertexBuffer := gl.CreateBuffer()
	//indexBuffer := gl.CreateBuffer()

	renderer.vbo = vertexBuffer
	renderer.ebo = indexBuffer
	renderer.gl = gl

	return &renderer
}

func (w *WebGLRenderer) NewFrame() {

}
func (w *WebGLRenderer) SetVertexAttribPointer(index int, size int, xtype js.Value, stride, offset int) {
	var memSize int = 4

	//w.gl.VertexAttribPointer(index, size, xtype, false, stride, offset*memSize)
	//w.gl.EnableVertexAttribArray(index)
	w.gl.Call("vertexAttribPointer", index, size, xtype, false, int32(stride), offset*memSize)
	w.gl.Call("enableVertexAttribArray", index)
}

func (w *WebGLRenderer) Draw(displaySize [2]float32, buffer draw.CmdBuffer) {
	displayWidth := displaySize[0]
	displayHeight := displaySize[1]

	js.Global().Get("console").Call("log", displayWidth, displayHeight)

	gl := w.gl

	gl.Call("useProgram", w.programID)
	//gl.UseProgram(w.programID)
	//vaoId := gl.Call("createVertexArray")
	//gl.Call("bindVertexArray", vaoId)

	//vaoId := gl.CreateVertexArray()
	//gl.BindVertexArray(vaoId)

	//gl.BindBuffer(gl.ARRAY_BUFFER, w.vbo)
	gl.Call("bindBuffer", glTypes.ArrayBuffer, w.vbo)

	w.SetVertexAttribPointer(0, posSize, glTypes.Float, vertexSize*4, posOffset)
	w.SetVertexAttribPointer(1, colorSize, glTypes.Float, vertexSize*4, colorOffset)
	w.SetVertexAttribPointer(2, texCoordsSize, glTypes.Float, vertexSize*4, texCoordsOffset)
	w.SetVertexAttribPointer(3, texIdSize, glTypes.Float, vertexSize*4, texIdOffset)

	gl.Call("enable", gl.Get("SCISSOR_TEST"))
	//gl.Enable(gl.SCISSOR_TEST)

	//gl.BindBuffer(gl.ARRAY_BUFFER, w.vbo)
	//gl.BufferData(gl.ARRAY_BUFFER, gltypes.SliceToTypedArray(buffer.Vertices), gl.STREAM_DRAW)
	gl.Call("bindBuffer", glTypes.ArrayBuffer, w.vbo)
	gl.Call("bufferData", glTypes.ArrayBuffer, gltypes.SliceToTypedArray(buffer.Vertices), gl.Get("STREAM_DRAW"))

	gl.Call("bindBuffer", glTypes.ElementArrayBuffer, w.ebo)
	gl.Call("bufferData", glTypes.ElementArrayBuffer, gltypes.SliceToTypedArray(buffer.Indices), gl.Get("STREAM_DRAW"))

	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, w.ebo)
	//gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gltypes.SliceToTypedArray(buffer.Indices), gl.STREAM_DRAW)

	orthoProjection := [16]float32{
		2.0 / displayWidth, 0.0, 0.0, 0.0,
		0.0, 2.0 / displayHeight, 0.0, 0.0,
		0.0, 0.0, -2.0, 0.0,
		-1.0, -1.0, -1.0, 1.0,
	}
	//ProjectionMatrix := gl.GetUniformLocation(w.programID, "uProjection")
	ProjectionMatrix := gl.Call("getUniformLocation", w.programID, "uProjection")
	//gl.UseProgram(w.programID)
	gl.Call("useProgram", w.programID)

	var projMatrixBuffer *[16]float32
	projMatrixBuffer = (*[16]float32)(unsafe.Pointer(&orthoProjection))
	typedProjMatrixBuffer := gltypes.SliceToTypedArray([]float32((*projMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", ProjectionMatrix, false, typedProjMatrixBuffer)

	//gl.UniformMatrix4fv(ProjectionMatrix, false, typedProjMatrixBuffer)

	for _, cmd := range buffer.DrawCalls {

		clipRect := cmd.ClipRect

		x := int(clipRect[0])
		y := int(clipRect[1])
		w := int(clipRect[2])
		h := int(clipRect[3])

		y = int(displayHeight) - (y + h)
		//gl.Scissor(x, y, w, h)
		gl.Call("scissor", x, y, w, h)
		//gl.DrawElements(gl.TRIANGLES, int32(cmd.Elems), gl.UNSIGNED_INT, cmd.IndexOffset*4)
		gl.Call("drawElements", glTypes.Triangles, int32(cmd.Elems), gl.Get("UNSIGNED_INT"), cmd.IndexOffset*4)
	}
	js.Global().Call("requestAnimationFrame", w.renderFrame)

	//gl.UseProgram(js.Null())
	//gl.DeleteVertexArray(vaoId)
	//gl.Disable(gl.SCISSOR_TEST)
	//gl.Call("useProgram", 0)
	//gl.Call("deleteVertexArray", vaoId)
	gl.Call("disable", gl.Get("SCISSOR_TEST"))
}
