// Taken from: https://github.com/n2d/webgl/blob/master/webgl.go
// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file here: https://github.com/n2d/webgl/blob/master/LICENSE
package main

import (
	"errors"
	"reflect"
	"syscall/js"
)

type ContextAttributes struct {
	// If Alpha is true, the drawing buffer has an alpha channel for
	// the purposes of performing OpenGL destination alpha operations
	// and compositing with the page.
	Alpha bool

	// If Depth is true, the drawing buffer has a depth buffer of at least 16 bits.
	Depth bool

	// If Stencil is true, the drawing buffer has a stencil buffer of at least 8 bits.
	Stencil bool

	// If Antialias is true and the implementation supports antialiasing
	// the drawing buffer will perform antialiasing using its choice of
	// technique (multisample/supersample) and quality.
	Antialias bool

	// If PremultipliedAlpha is true the page compositor will assume the
	// drawing buffer contains colors with premultiplied alpha.
	// This flag is ignored if the alpha flag is false.
	PremultipliedAlpha bool

	// If the value is true the buffers will not be cleared and will preserve
	// their values until cleared or overwritten by the author.
	PreserveDrawingBuffer bool
}

// Returns a copy of the default WebGL context attributes.
func DefaultAttributes() *ContextAttributes {
	return &ContextAttributes{true, true, false, true, true, false}
}

type Context struct {
	js.Value
	ARRAY_BUFFER                                 js.Value `js:"ARRAY_BUFFER"`
	ARRAY_BUFFER_BINDING                         js.Value `js:"ARRAY_BUFFER_BINDING"`
	ATTACHED_SHADERS                             js.Value `js:"ATTACHED_SHADERS"`
	BACK                                         js.Value `js:"BACK"`
	BLEND                                        js.Value `js:"BLEND"`
	BLEND_COLOR                                  js.Value `js:"BLEND_COLOR"`
	BLEND_DST_ALPHA                              js.Value `js:"BLEND_DST_ALPHA"`
	BLEND_DST_RGB                                js.Value `js:"BLEND_DST_RGB"`
	BLEND_EQUATION                               js.Value `js:"BLEND_EQUATION"`
	BLEND_EQUATION_ALPHA                         js.Value `js:"BLEND_EQUATION_ALPHA"`
	BLEND_EQUATION_RGB                           js.Value `js:"BLEND_EQUATION_RGB"`
	BLEND_SRC_ALPHA                              js.Value `js:"BLEND_SRC_ALPHA"`
	BLEND_SRC_RGB                                js.Value `js:"BLEND_SRC_RGB"`
	BLUE_BITS                                    js.Value `js:"BLUE_BITS"`
	BOOL                                         js.Value `js:"BOOL"`
	BOOL_VEC2                                    js.Value `js:"BOOL_VEC2"`
	BOOL_VEC3                                    js.Value `js:"BOOL_VEC3"`
	BOOL_VEC4                                    js.Value `js:"BOOL_VEC4"`
	BROWSER_DEFAULT_WEBGL                        js.Value `js:"BROWSER_DEFAULT_WEBGL"`
	BUFFER_SIZE                                  js.Value `js:"BUFFER_SIZE"`
	BUFFER_USAGE                                 js.Value `js:"BUFFER_USAGE"`
	BYTE                                         js.Value `js:"BYTE"`
	CCW                                          js.Value `js:"CCW"`
	CLAMP_TO_EDGE                                js.Value `js:"CLAMP_TO_EDGE"`
	COLOR_ATTACHMENT0                            js.Value `js:"COLOR_ATTACHMENT0"`
	COLOR_BUFFER_BIT                             js.Value `js:"COLOR_BUFFER_BIT"`
	COLOR_CLEAR_VALUE                            js.Value `js:"COLOR_CLEAR_VALUE"`
	COLOR_WRITEMASK                              js.Value `js:"COLOR_WRITEMASK"`
	COMPILE_STATUS                               js.Value `js:"COMPILE_STATUS"`
	COMPRESSED_TEXTURE_FORMATS                   js.Value `js:"COMPRESSED_TEXTURE_FORMATS"`
	CONSTANT_ALPHA                               js.Value `js:"CONSTANT_ALPHA"`
	CONSTANT_COLOR                               js.Value `js:"CONSTANT_COLOR"`
	CONTEXT_LOST_WEBGL                           js.Value `js:"CONTEXT_LOST_WEBGL"`
	CULL_FACE                                    js.Value `js:"CULL_FACE"`
	CULL_FACE_MODE                               js.Value `js:"CULL_FACE_MODE"`
	CURRENT_PROGRAM                              js.Value `js:"CURRENT_PROGRAM"`
	CURRENT_VERTEX_ATTRIB                        js.Value `js:"CURRENT_VERTEX_ATTRIB"`
	CW                                           js.Value `js:"CW"`
	DECR                                         js.Value `js:"DECR"`
	DECR_WRAP                                    js.Value `js:"DECR_WRAP"`
	DELETE_STATUS                                js.Value `js:"DELETE_STATUS"`
	DEPTH_ATTACHMENT                             js.Value `js:"DEPTH_ATTACHMENT"`
	DEPTH_BITS                                   js.Value `js:"DEPTH_BITS"`
	DEPTH_BUFFER_BIT                             js.Value `js:"DEPTH_BUFFER_BIT"`
	DEPTH_CLEAR_VALUE                            js.Value `js:"DEPTH_CLEAR_VALUE"`
	DEPTH_COMPONENT                              js.Value `js:"DEPTH_COMPONENT"`
	DEPTH_COMPONENT16                            js.Value `js:"DEPTH_COMPONENT16"`
	DEPTH_FUNC                                   js.Value `js:"DEPTH_FUNC"`
	DEPTH_RANGE                                  js.Value `js:"DEPTH_RANGE"`
	DEPTH_STENCIL                                js.Value `js:"DEPTH_STENCIL"`
	DEPTH_STENCIL_ATTACHMENT                     js.Value `js:"DEPTH_STENCIL_ATTACHMENT"`
	DEPTH_TEST                                   js.Value `js:"DEPTH_TEST"`
	DEPTH_WRITEMASK                              js.Value `js:"DEPTH_WRITEMASK"`
	DITHER                                       js.Value `js:"DITHER"`
	DONT_CARE                                    js.Value `js:"DONT_CARE"`
	DST_ALPHA                                    js.Value `js:"DST_ALPHA"`
	DST_COLOR                                    js.Value `js:"DST_COLOR"`
	DYNAMIC_DRAW                                 js.Value `js:"DYNAMIC_DRAW"`
	ELEMENT_ARRAY_BUFFER                         js.Value `js:"ELEMENT_ARRAY_BUFFER"`
	ELEMENT_ARRAY_BUFFER_BINDING                 js.Value `js:"ELEMENT_ARRAY_BUFFER_BINDING"`
	EQUAL                                        js.Value `js:"EQUAL"`
	FASTEST                                      js.Value `js:"FASTEST"`
	FLOAT                                        js.Value `js:"FLOAT"`
	FLOAT_MAT2                                   js.Value `js:"FLOAT_MAT2"`
	FLOAT_MAT3                                   js.Value `js:"FLOAT_MAT3"`
	FLOAT_MAT4                                   js.Value `js:"FLOAT_MAT4"`
	FLOAT_VEC2                                   js.Value `js:"FLOAT_VEC2"`
	FLOAT_VEC3                                   js.Value `js:"FLOAT_VEC3"`
	FLOAT_VEC4                                   js.Value `js:"FLOAT_VEC4"`
	FRAGMENT_SHADER                              js.Value `js:"FRAGMENT_SHADER"`
	FRAMEBUFFER                                  js.Value `js:"FRAMEBUFFER"`
	FRAMEBUFFER_ATTACHMENT_OBJECT_NAME           js.Value `js:"FRAMEBUFFER_ATTACHMENT_OBJECT_NAME"`
	FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE           js.Value `js:"FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE"`
	FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE js.Value `js:"FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE"`
	FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL         js.Value `js:"FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL"`
	FRAMEBUFFER_BINDING                          js.Value `js:"FRAMEBUFFER_BINDING"`
	FRAMEBUFFER_COMPLETE                         js.Value `js:"FRAMEBUFFER_COMPLETE"`
	FRAMEBUFFER_INCOMPLETE_ATTACHMENT            js.Value `js:"FRAMEBUFFER_INCOMPLETE_ATTACHMENT"`
	FRAMEBUFFER_INCOMPLETE_DIMENSIONS            js.Value `js:"FRAMEBUFFER_INCOMPLETE_DIMENSIONS"`
	FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT    js.Value `js:"FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT"`
	FRAMEBUFFER_UNSUPPORTED                      js.Value `js:"FRAMEBUFFER_UNSUPPORTED"`
	FRONT                                        js.Value `js:"FRONT"`
	FRONT_AND_BACK                               js.Value `js:"FRONT_AND_BACK"`
	FRONT_FACE                                   js.Value `js:"FRONT_FACE"`
	FUNC_ADD                                     js.Value `js:"FUNC_ADD"`
	FUNC_REVERSE_SUBTRACT                        js.Value `js:"FUNC_REVERSE_SUBTRACT"`
	FUNC_SUBTRACT                                js.Value `js:"FUNC_SUBTRACT"`
	GENERATE_MIPMAP_HINT                         js.Value `js:"GENERATE_MIPMAP_HINT"`
	GEQUAL                                       js.Value `js:"GEQUAL"`
	GREATER                                      js.Value `js:"GREATER"`
	GREEN_BITS                                   js.Value `js:"GREEN_BITS"`
	HIGH_FLOAT                                   js.Value `js:"HIGH_FLOAT"`
	HIGH_INT                                     js.Value `js:"HIGH_INT"`
	INCR                                         js.Value `js:"INCR"`
	INCR_WRAP                                    js.Value `js:"INCR_WRAP"`
	INFO_LOG_LENGTH                              js.Value `js:"INFO_LOG_LENGTH"`
	INT                                          js.Value `js:"INT"`
	INT_VEC2                                     js.Value `js:"INT_VEC2"`
	INT_VEC3                                     js.Value `js:"INT_VEC3"`
	INT_VEC4                                     js.Value `js:"INT_VEC4"`
	INVALID_ENUM                                 js.Value `js:"INVALID_ENUM"`
	INVALID_FRAMEBUFFER_OPERATION                js.Value `js:"INVALID_FRAMEBUFFER_OPERATION"`
	INVALID_OPERATION                            js.Value `js:"INVALID_OPERATION"`
	INVALID_VALUE                                js.Value `js:"INVALID_VALUE"`
	INVERT                                       js.Value `js:"INVERT"`
	KEEP                                         js.Value `js:"KEEP"`
	LEQUAL                                       js.Value `js:"LEQUAL"`
	LESS                                         js.Value `js:"LESS"`
	LINEAR                                       js.Value `js:"LINEAR"`
	LINEAR_MIPMAP_LINEAR                         js.Value `js:"LINEAR_MIPMAP_LINEAR"`
	LINEAR_MIPMAP_NEAREST                        js.Value `js:"LINEAR_MIPMAP_NEAREST"`
	LINES                                        js.Value `js:"LINES"`
	LINE_LOOP                                    js.Value `js:"LINE_LOOP"`
	LINE_STRIP                                   js.Value `js:"LINE_STRIP"`
	LINE_WIDTH                                   js.Value `js:"LINE_WIDTH"`
	LINK_STATUS                                  js.Value `js:"LINK_STATUS"`
	LOW_FLOAT                                    js.Value `js:"LOW_FLOAT"`
	LOW_INT                                      js.Value `js:"LOW_INT"`
	LUMINANCE                                    js.Value `js:"LUMINANCE"`
	LUMINANCE_ALPHA                              js.Value `js:"LUMINANCE_ALPHA"`
	MAX_COMBINED_TEXTURE_IMAGE_UNITS             js.Value `js:"MAX_COMBINED_TEXTURE_IMAGE_UNITS"`
	MAX_CUBE_MAP_TEXTURE_SIZE                    js.Value `js:"MAX_CUBE_MAP_TEXTURE_SIZE"`
	MAX_FRAGMENT_UNIFORM_VECTORS                 js.Value `js:"MAX_FRAGMENT_UNIFORM_VECTORS"`
	MAX_RENDERBUFFER_SIZE                        js.Value `js:"MAX_RENDERBUFFER_SIZE"`
	MAX_TEXTURE_IMAGE_UNITS                      js.Value `js:"MAX_TEXTURE_IMAGE_UNITS"`
	MAX_TEXTURE_SIZE                             js.Value `js:"MAX_TEXTURE_SIZE"`
	MAX_VARYING_VECTORS                          js.Value `js:"MAX_VARYING_VECTORS"`
	MAX_VERTEX_ATTRIBS                           js.Value `js:"MAX_VERTEX_ATTRIBS"`
	MAX_VERTEX_TEXTURE_IMAGE_UNITS               js.Value `js:"MAX_VERTEX_TEXTURE_IMAGE_UNITS"`
	MAX_VERTEX_UNIFORM_VECTORS                   js.Value `js:"MAX_VERTEX_UNIFORM_VECTORS"`
	MAX_VIEWPORT_DIMS                            js.Value `js:"MAX_VIEWPORT_DIMS"`
	MEDIUM_FLOAT                                 js.Value `js:"MEDIUM_FLOAT"`
	MEDIUM_INT                                   js.Value `js:"MEDIUM_INT"`
	MIRRORED_REPEAT                              js.Value `js:"MIRRORED_REPEAT"`
	NEAREST                                      js.Value `js:"NEAREST"`
	NEAREST_MIPMAP_LINEAR                        js.Value `js:"NEAREST_MIPMAP_LINEAR"`
	NEAREST_MIPMAP_NEAREST                       js.Value `js:"NEAREST_MIPMAP_NEAREST"`
	NEVER                                        js.Value `js:"NEVER"`
	NICEST                                       js.Value `js:"NICEST"`
	NONE                                         js.Value `js:"NONE"`
	NOTEQUAL                                     js.Value `js:"NOTEQUAL"`
	NO_ERROR                                     js.Value `js:"NO_ERROR"`
	NUM_COMPRESSED_TEXTURE_FORMATS               js.Value `js:"NUM_COMPRESSED_TEXTURE_FORMATS"`
	ONE                                          js.Value `js:"ONE"`
	ONE_MINUS_CONSTANT_ALPHA                     js.Value `js:"ONE_MINUS_CONSTANT_ALPHA"`
	ONE_MINUS_CONSTANT_COLOR                     js.Value `js:"ONE_MINUS_CONSTANT_COLOR"`
	ONE_MINUS_DST_ALPHA                          js.Value `js:"ONE_MINUS_DST_ALPHA"`
	ONE_MINUS_DST_COLOR                          js.Value `js:"ONE_MINUS_DST_COLOR"`
	ONE_MINUS_SRC_ALPHA                          js.Value `js:"ONE_MINUS_SRC_ALPHA"`
	ONE_MINUS_SRC_COLOR                          js.Value `js:"ONE_MINUS_SRC_COLOR"`
	OUT_OF_MEMORY                                js.Value `js:"OUT_OF_MEMORY"`
	PACK_ALIGNMENT                               js.Value `js:"PACK_ALIGNMENT"`
	POINTS                                       js.Value `js:"POINTS"`
	POLYGON_OFFSET_FACTOR                        js.Value `js:"POLYGON_OFFSET_FACTOR"`
	POLYGON_OFFSET_FILL                          js.Value `js:"POLYGON_OFFSET_FILL"`
	POLYGON_OFFSET_UNITS                         js.Value `js:"POLYGON_OFFSET_UNITS"`
	RED_BITS                                     js.Value `js:"RED_BITS"`
	RENDERBUFFER                                 js.Value `js:"RENDERBUFFER"`
	RENDERBUFFER_ALPHA_SIZE                      js.Value `js:"RENDERBUFFER_ALPHA_SIZE"`
	RENDERBUFFER_BINDING                         js.Value `js:"RENDERBUFFER_BINDING"`
	RENDERBUFFER_BLUE_SIZE                       js.Value `js:"RENDERBUFFER_BLUE_SIZE"`
	RENDERBUFFER_DEPTH_SIZE                      js.Value `js:"RENDERBUFFER_DEPTH_SIZE"`
	RENDERBUFFER_GREEN_SIZE                      js.Value `js:"RENDERBUFFER_GREEN_SIZE"`
	RENDERBUFFER_HEIGHT                          js.Value `js:"RENDERBUFFER_HEIGHT"`
	RENDERBUFFER_INTERNAL_FORMAT                 js.Value `js:"RENDERBUFFER_INTERNAL_FORMAT"`
	RENDERBUFFER_RED_SIZE                        js.Value `js:"RENDERBUFFER_RED_SIZE"`
	RENDERBUFFER_STENCIL_SIZE                    js.Value `js:"RENDERBUFFER_STENCIL_SIZE"`
	RENDERBUFFER_WIDTH                           js.Value `js:"RENDERBUFFER_WIDTH"`
	RENDERER                                     js.Value `js:"RENDERER"`
	REPEAT                                       js.Value `js:"REPEAT"`
	REPLACE                                      js.Value `js:"REPLACE"`
	RGB                                          js.Value `js:"RGB"`
	RGB5_A1                                      js.Value `js:"RGB5_A1"`
	RGB565                                       js.Value `js:"RGB565"`
	RGBA                                         js.Value `js:"RGBA"`
	RGBA4                                        js.Value `js:"RGBA4"`
	SAMPLER_2D                                   js.Value `js:"SAMPLER_2D"`
	SAMPLER_CUBE                                 js.Value `js:"SAMPLER_CUBE"`
	SAMPLES                                      js.Value `js:"SAMPLES"`
	SAMPLE_ALPHA_TO_COVERAGE                     js.Value `js:"SAMPLE_ALPHA_TO_COVERAGE"`
	SAMPLE_BUFFERS                               js.Value `js:"SAMPLE_BUFFERS"`
	SAMPLE_COVERAGE                              js.Value `js:"SAMPLE_COVERAGE"`
	SAMPLE_COVERAGE_INVERT                       js.Value `js:"SAMPLE_COVERAGE_INVERT"`
	SAMPLE_COVERAGE_VALUE                        js.Value `js:"SAMPLE_COVERAGE_VALUE"`
	SCISSOR_BOX                                  js.Value `js:"SCISSOR_BOX"`
	SCISSOR_TEST                                 js.Value `js:"SCISSOR_TEST"`
	SHADER_COMPILER                              js.Value `js:"SHADER_COMPILER"`
	SHADER_SOURCE_LENGTH                         js.Value `js:"SHADER_SOURCE_LENGTH"`
	SHADER_TYPE                                  js.Value `js:"SHADER_TYPE"`
	SHADING_LANGUAGE_VERSION                     js.Value `js:"SHADING_LANGUAGE_VERSION"`
	SHORT                                        js.Value `js:"SHORT"`
	SRC_ALPHA                                    js.Value `js:"SRC_ALPHA"`
	SRC_ALPHA_SATURATE                           js.Value `js:"SRC_ALPHA_SATURATE"`
	SRC_COLOR                                    js.Value `js:"SRC_COLOR"`
	STATIC_DRAW                                  js.Value `js:"STATIC_DRAW"`
	STENCIL_ATTACHMENT                           js.Value `js:"STENCIL_ATTACHMENT"`
	STENCIL_BACK_FAIL                            js.Value `js:"STENCIL_BACK_FAIL"`
	STENCIL_BACK_FUNC                            js.Value `js:"STENCIL_BACK_FUNC"`
	STENCIL_BACK_PASS_DEPTH_FAIL                 js.Value `js:"STENCIL_BACK_PASS_DEPTH_FAIL"`
	STENCIL_BACK_PASS_DEPTH_PASS                 js.Value `js:"STENCIL_BACK_PASS_DEPTH_PASS"`
	STENCIL_BACK_REF                             js.Value `js:"STENCIL_BACK_REF"`
	STENCIL_BACK_VALUE_MASK                      js.Value `js:"STENCIL_BACK_VALUE_MASK"`
	STENCIL_BACK_WRITEMASK                       js.Value `js:"STENCIL_BACK_WRITEMASK"`
	STENCIL_BITS                                 js.Value `js:"STENCIL_BITS"`
	STENCIL_BUFFER_BIT                           js.Value `js:"STENCIL_BUFFER_BIT"`
	STENCIL_CLEAR_VALUE                          js.Value `js:"STENCIL_CLEAR_VALUE"`
	STENCIL_FAIL                                 js.Value `js:"STENCIL_FAIL"`
	STENCIL_FUNC                                 js.Value `js:"STENCIL_FUNC"`
	STENCIL_INDEX                                js.Value `js:"STENCIL_INDEX"`
	STENCIL_INDEX8                               js.Value `js:"STENCIL_INDEX8"`
	STENCIL_PASS_DEPTH_FAIL                      js.Value `js:"STENCIL_PASS_DEPTH_FAIL"`
	STENCIL_PASS_DEPTH_PASS                      js.Value `js:"STENCIL_PASS_DEPTH_PASS"`
	STENCIL_REF                                  js.Value `js:"STENCIL_REF"`
	STENCIL_TEST                                 js.Value `js:"STENCIL_TEST"`
	STENCIL_VALUE_MASK                           js.Value `js:"STENCIL_VALUE_MASK"`
	STENCIL_WRITEMASK                            js.Value `js:"STENCIL_WRITEMASK"`
	STREAM_DRAW                                  js.Value `js:"STREAM_DRAW"`
	SUBPIXEL_BITS                                js.Value `js:"SUBPIXEL_BITS"`
	TEXTURE                                      js.Value `js:"TEXTURE"`
	TEXTURE0                                     js.Value `js:"TEXTURE0"`
	TEXTURE1                                     js.Value `js:"TEXTURE1"`
	TEXTURE2                                     js.Value `js:"TEXTURE2"`
	TEXTURE3                                     js.Value `js:"TEXTURE3"`
	TEXTURE4                                     js.Value `js:"TEXTURE4"`
	TEXTURE5                                     js.Value `js:"TEXTURE5"`
	TEXTURE6                                     js.Value `js:"TEXTURE6"`
	TEXTURE7                                     js.Value `js:"TEXTURE7"`
	TEXTURE8                                     js.Value `js:"TEXTURE8"`
	TEXTURE9                                     js.Value `js:"TEXTURE9"`
	TEXTURE10                                    js.Value `js:"TEXTURE10"`
	TEXTURE11                                    js.Value `js:"TEXTURE11"`
	TEXTURE12                                    js.Value `js:"TEXTURE12"`
	TEXTURE13                                    js.Value `js:"TEXTURE13"`
	TEXTURE14                                    js.Value `js:"TEXTURE14"`
	TEXTURE15                                    js.Value `js:"TEXTURE15"`
	TEXTURE16                                    js.Value `js:"TEXTURE16"`
	TEXTURE17                                    js.Value `js:"TEXTURE17"`
	TEXTURE18                                    js.Value `js:"TEXTURE18"`
	TEXTURE19                                    js.Value `js:"TEXTURE19"`
	TEXTURE20                                    js.Value `js:"TEXTURE20"`
	TEXTURE21                                    js.Value `js:"TEXTURE21"`
	TEXTURE22                                    js.Value `js:"TEXTURE22"`
	TEXTURE23                                    js.Value `js:"TEXTURE23"`
	TEXTURE24                                    js.Value `js:"TEXTURE24"`
	TEXTURE25                                    js.Value `js:"TEXTURE25"`
	TEXTURE26                                    js.Value `js:"TEXTURE26"`
	TEXTURE27                                    js.Value `js:"TEXTURE27"`
	TEXTURE28                                    js.Value `js:"TEXTURE28"`
	TEXTURE29                                    js.Value `js:"TEXTURE29"`
	TEXTURE30                                    js.Value `js:"TEXTURE30"`
	TEXTURE31                                    js.Value `js:"TEXTURE31"`
	TEXTURE_2D                                   js.Value `js:"TEXTURE_2D"`
	TEXTURE_BINDING_2D                           js.Value `js:"TEXTURE_BINDING_2D"`
	TEXTURE_BINDING_CUBE_MAP                     js.Value `js:"TEXTURE_BINDING_CUBE_MAP"`
	TEXTURE_CUBE_MAP                             js.Value `js:"TEXTURE_CUBE_MAP"`
	TEXTURE_CUBE_MAP_NEGATIVE_X                  js.Value `js:"TEXTURE_CUBE_MAP_NEGATIVE_X"`
	TEXTURE_CUBE_MAP_NEGATIVE_Y                  js.Value `js:"TEXTURE_CUBE_MAP_NEGATIVE_Y"`
	TEXTURE_CUBE_MAP_NEGATIVE_Z                  js.Value `js:"TEXTURE_CUBE_MAP_NEGATIVE_Z"`
	TEXTURE_CUBE_MAP_POSITIVE_X                  js.Value `js:"TEXTURE_CUBE_MAP_POSITIVE_X"`
	TEXTURE_CUBE_MAP_POSITIVE_Y                  js.Value `js:"TEXTURE_CUBE_MAP_POSITIVE_Y"`
	TEXTURE_CUBE_MAP_POSITIVE_Z                  js.Value `js:"TEXTURE_CUBE_MAP_POSITIVE_Z"`
	TEXTURE_MAG_FILTER                           js.Value `js:"TEXTURE_MAG_FILTER"`
	TEXTURE_MIN_FILTER                           js.Value `js:"TEXTURE_MIN_FILTER"`
	TEXTURE_WRAP_S                               js.Value `js:"TEXTURE_WRAP_S"`
	TEXTURE_WRAP_T                               js.Value `js:"TEXTURE_WRAP_T"`
	TRIANGLES                                    js.Value `js:"TRIANGLES"`
	TRIANGLE_FAN                                 js.Value `js:"TRIANGLE_FAN"`
	TRIANGLE_STRIP                               js.Value `js:"TRIANGLE_STRIP"`
	UNPACK_ALIGNMENT                             js.Value `js:"UNPACK_ALIGNMENT"`
	UNPACK_COLORSPACE_CONVERSION_WEBGL           js.Value `js:"UNPACK_COLORSPACE_CONVERSION_WEBGL"`
	UNPACK_FLIP_Y_WEBGL                          js.Value `js:"UNPACK_FLIP_Y_WEBGL"`
	UNPACK_PREMULTIPLY_ALPHA_WEBGL               js.Value `js:"UNPACK_PREMULTIPLY_ALPHA_WEBGL"`
	UNSIGNED_BYTE                                js.Value `js:"UNSIGNED_BYTE"`
	UNSIGNED_INT                                 js.Value `js:"UNSIGNED_INT"`
	UNSIGNED_SHORT                               js.Value `js:"UNSIGNED_SHORT"`
	UNSIGNED_SHORT_4_4_4_4                       js.Value `js:"UNSIGNED_SHORT_4_4_4_4"`
	UNSIGNED_SHORT_5_5_5_1                       js.Value `js:"UNSIGNED_SHORT_5_5_5_1"`
	UNSIGNED_SHORT_5_6_5                         js.Value `js:"UNSIGNED_SHORT_5_6_5"`
	VALIDATE_STATUS                              js.Value `js:"VALIDATE_STATUS"`
	VENDOR                                       js.Value `js:"VENDOR"`
	VERSION                                      js.Value `js:"VERSION"`
	VERTEX_ATTRIB_ARRAY_BUFFER_BINDING           js.Value `js:"VERTEX_ATTRIB_ARRAY_BUFFER_BINDING"`
	VERTEX_ATTRIB_ARRAY_ENABLED                  js.Value `js:"VERTEX_ATTRIB_ARRAY_ENABLED"`
	VERTEX_ATTRIB_ARRAY_NORMALIZED               js.Value `js:"VERTEX_ATTRIB_ARRAY_NORMALIZED"`
	VERTEX_ATTRIB_ARRAY_POINTER                  js.Value `js:"VERTEX_ATTRIB_ARRAY_POINTER"`
	VERTEX_ATTRIB_ARRAY_SIZE                     js.Value `js:"VERTEX_ATTRIB_ARRAY_SIZE"`
	VERTEX_ATTRIB_ARRAY_STRIDE                   js.Value `js:"VERTEX_ATTRIB_ARRAY_STRIDE"`
	VERTEX_ATTRIB_ARRAY_TYPE                     js.Value `js:"VERTEX_ATTRIB_ARRAY_TYPE"`
	VERTEX_SHADER                                js.Value `js:"VERTEX_SHADER"`
	VIEWPORT                                     js.Value `js:"VIEWPORT"`
	ZERO                                         js.Value `js:"ZERO"`
}

// NewContext takes an HTML5 canvas object and optional context attributes.
// If an error is returned it means you won't have access to WebGL
// functionality.
func NewContext(canvas js.Value) (*Context, error) {
	if js.Global().Get("WebGLRenderingContext").Equal(js.Undefined()) {
		return nil, errors.New("Your browser doesn't appear to support webgl.")
	}

	gl := canvas.Call("getContext", "webgl2")
	if gl.IsNull() {
		gl = canvas.Call("getContext", "experimental-webgl")
		if gl.IsNull() {
			return nil, errors.New("Creating a webgl context has failed.")
		}
	}
	ctx := new(Context)
	ctx.Value = gl

	rctx := reflect.TypeOf(*ctx)
	for i := 0; i < rctx.NumField(); i++ {
		field := rctx.Field(i)

		if _, ok := field.Tag.Lookup("js"); ok {
			val := gl.Get(field.Tag.Get("js"))
			reflect.ValueOf(ctx).Elem().Field(i).Set(reflect.ValueOf(val))
		}
	}
	return ctx, nil
}

// Returns the context attributes active on the context. These values might
// be different than what was requested on context creation if the
// browser's implementation doesn't support a feature.
func (c *Context) GetContextAttributes() ContextAttributes {
	ca := c.Call("getContextAttributes")
	return ContextAttributes{
		ca.Get("alpha").Bool(),
		ca.Get("depth").Bool(),
		ca.Get("stencil").Bool(),
		ca.Get("antialias").Bool(),
		ca.Get("premultipliedAlpha").Bool(),
		ca.Get("preservedDrawingBuffer").Bool(),
	}
}

// Specifies the active texture unit.
func (c *Context) ActiveTexture(texture int) {
	c.Call("activeTexture", texture)
}

// Attaches a WebGLShader object to a WebGLProgram object.
func (c *Context) AttachShader(program js.Value, shader js.Value) {
	c.Call("attachShader", program, shader)
}

// Binds a generic vertex index to a user-defined attribute variable.
func (c *Context) BindAttribLocation(program js.Value, index int, name string) {
	c.Call("bindAttribLocation", program, index, name)
}

// Associates a buffer with a buffer target.
func (c *Context) BindBuffer(target js.Value, buffer js.Value) {
	c.Call("bindBuffer", target, buffer)
}

// Associates a WebGLFramebuffer object with the FRAMEBUFFER bind target.
func (c *Context) BindFramebuffer(target js.Value, framebuffer js.Value) {
	c.Call("bindFramebuffer", target, framebuffer)
}

// Binds a WebGLRenderbuffer object to be used for rendering.
func (c *Context) BindRenderbuffer(target js.Value, renderbuffer js.Value) {
	c.Call("bindRenderbuffer", target, renderbuffer)
}

// Binds a named texture object to a target.
func (c *Context) BindTexture(target int, texture js.Value) {
	c.Call("bindTexture", target, texture)
}

// The GL_BLEND_COLOR may be used to calculate the source and destination blending factors.
func (c *Context) BlendColor(r, g, b, a float64) {
	c.Call("blendColor", r, g, b, a)
}

// Sets the equation used to blend RGB and Alpha values of an incoming source
// fragment with a destination values as stored in the fragment's frame buffer.
func (c *Context) BlendEquation(mode int) {
	c.Call("blendEquation", mode)
}

// Controls the blending of an incoming source fragment's R, G, B, and A values
// with a destination R, G, B, and A values as stored in the fragment's WebGLFramebuffer.
func (c *Context) BlendEquationSeparate(modeRGB, modeAlpha int) {
	c.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

// Sets the blending factors used to combine source and destination pixels.
func (c *Context) BlendFunc(sfactor, dfactor js.Value) {
	c.Call("blendFunc", sfactor, dfactor)
}

// Sets the weighting factors that are used by blendEquationSeparate.
func (c *Context) BlendFuncSeparate(srcRGB, dstRGB, srcAlpha, dstAlpha int) {
	c.Call("blendFuncSeparate", srcRGB, dstRGB, srcAlpha, dstAlpha)
}

// Creates a buffer in memory and initializes it with array data.
// If no array is provided, the contents of the buffer is initialized to 0.
func (c *Context) BufferData(target js.Value, data js.Value, usage js.Value) {
	c.Call("bufferData", target, data, usage)
}

// Used to modify or update some or all of a data store for a bound buffer object.
func (c *Context) BufferSubData(target js.Value, offset int, data interface{}) {
	c.Call("bufferSubData", target, offset, data)
}

// Returns whether the currently bound WebGLFramebuffer is complete.
// If not complete, returns the reason why.
func (c *Context) CheckFramebufferStatus(target int) int {
	return c.Call("checkFramebufferStatus", target).Int()
}

// Sets all pixels in a specific buffer to the same value.
func (c *Context) Clear(flags js.Value) {
	c.Call("clear", flags)
}

// Specifies color values to use by the clear method to clear the color buffer.
func (c *Context) ClearColor(r, g, b, a float32) {
	c.Call("clearColor", r, g, b, a)
}

// Clears the depth buffer to a specific value.
func (c *Context) ClearDepth(depth float64) {
	c.Call("clearDepth", depth)
}

func (c *Context) ClearStencil(s int) {
	c.Call("clearStencil", s)
}

// Lets you set whether individual colors can be written when
// drawing or rendering to a framebuffer.
func (c *Context) ColorMask(r, g, b, a bool) {
	c.Call("colorMask", r, g, b, a)
}

// Compiles the GLSL shader source into binary data used by the WebGLProgram object.
func (c *Context) CompileShader(shader js.Value) {
	c.Call("compileShader", shader)
}

// Copies a rectangle of pixels from the current WebGLFramebuffer into a texture image.
func (c *Context) CopyTexImage2D(target, level, internal, x, y, w, h, border int) {
	c.Call("copyTexImage2D", target, level, internal, x, y, w, h, border)
}

// Replaces a portion of an existing 2D texture image with data from the current framebuffer.
func (c *Context) CopyTexSubImage2D(target, level, xoffset, yoffset, x, y, w, h int) {
	c.Call("copyTexSubImage2D", target, level, xoffset, yoffset, x, y, w, h)
}

// Creates and initializes a WebGLBuffer.
func (c *Context) CreateBuffer() js.Value {
	return c.Call("createBuffer")
}

// Returns a WebGLFramebuffer object.
func (c *Context) CreateFramebuffer() js.Value {
	return c.Call("createFramebuffer")
}

// Creates an empty WebGLProgram object to which vector and fragment
// WebGLShader objects can be bound.
func (c *Context) CreateProgram() js.Value {
	return c.Call("createProgram")
}

// Creates and returns a WebGLRenderbuffer object.
func (c *Context) CreateRenderbuffer() js.Value {
	return c.Call("createRenderbuffer")
}

// Returns an empty vertex or fragment shader object based on the type specified.
func (c *Context) CreateShader(typ js.Value) js.Value {
	return c.Call("createShader", typ)
}

// Used to generate a WebGLTexture object to which images can be bound.
func (c *Context) CreateTexture() js.Value {
	return c.Call("createTexture")
}

// Sets whether or not front, back, or both facing facets are able to be culled.
func (c *Context) CullFace(mode int) {
	c.Call("cullFace", mode)
}

// Delete a specific buffer.
func (c *Context) DeleteBuffer(buffer js.Value) {
	c.Call("deleteBuffer", buffer)
}

// Deletes a specific WebGLFramebuffer object. If you delete the
// currently bound framebuffer, the default framebuffer will be bound.
// Deleting a framebuffer detaches all of its attachments.
func (c *Context) DeleteFramebuffer(framebuffer js.Value) {
	c.Call("deleteFramebuffer", framebuffer)
}

// Flags a specific WebGLProgram object for deletion if currently active.
// It will be deleted when it is no longer being used.
// Any shader objects associated with the program will be detached.
// They will be deleted if they were already flagged for deletion.
func (c *Context) DeleteProgram(program js.Value) {
	c.Call("deleteProgram", program)
}

// Deletes the specified renderbuffer object. If the renderbuffer is
// currently bound, it will become unbound. If the renderbuffer is
// attached to the currently bound framebuffer, it is detached.
func (c *Context) DeleteRenderbuffer(renderbuffer js.Value) {
	c.Call("deleteRenderbuffer", renderbuffer)
}

// Deletes a specific shader object.
func (c *Context) DeleteShader(shader js.Value) {
	c.Call("deleteShader", shader)
}

// Deletes a specific texture object.
func (c *Context) DeleteTexture(texture js.Value) {
	c.Call("deleteTexture", texture)
}

// Sets a function to use to compare incoming pixel depth to the
// current depth buffer value.
func (c *Context) DepthFunc(fun js.Value) {
	c.Call("depthFunc", fun)
}

// Sets whether or not you can write to the depth buffer.
func (c *Context) DepthMask(flag bool) {
	c.Call("depthMask", flag)
}

// Sets the depth range for normalized coordinates to canvas or viewport depth coordinates.
func (c *Context) DepthRange(zNear, zFar float64) {
	c.Call("depthRange", zNear, zFar)
}

// Detach a shader object from a program object.
func (c *Context) DetachShader(program, shader js.Value) {
	c.Call("detachShader", program, shader)
}

// Turns off specific WebGL capabilities for this context.
func (c *Context) Disable(cap js.Value) {
	c.Call("disable", cap)
}

// Turns off a vertex attribute array at a specific index position.
func (c *Context) DisableVertexAttribArray(index int) {
	c.Call("disableVertexAttribArray", index)
}

// Render geometric primitives from bound and enabled vertex data.
func (c *Context) DrawArrays(mode, first, count int) {
	c.Call("drawArrays", mode, first, count)
}

// Renders geometric primitives indexed by element array data.
func (c *Context) DrawElements(mode js.Value, count int32, typ js.Value, offset int) {
	c.Call("drawElements", mode, count, typ, offset)
}

// Turns on specific WebGL capabilities for this context.
func (c *Context) Enable(cap js.Value) {
	c.Call("enable", cap)
}

// Turns on a vertex attribute at a specific index position in
// a vertex attribute array.
func (c *Context) EnableVertexAttribArray(index int) {
	c.Call("enableVertexAttribArray", index)
}

func (c *Context) Finish() {
	c.Call("finish")
}

func (c *Context) Flush() {
	c.Call("flush")
}

// Attaches a WebGLRenderbuffer object as a logical buffer to the
// currently bound WebGLFramebuffer object.
func (c *Context) FrameBufferRenderBuffer(target, attachment, renderbufferTarget int, renderbuffer js.Value) {
	c.Call("framebufferRenderBuffer", target, attachment, renderbufferTarget, renderbuffer)
}

// Attaches a texture to a WebGLFramebuffer object.
func (c *Context) FramebufferTexture2D(target, attachment, textarget int, texture js.Value, level int) {
	c.Call("framebufferTexture2D", target, attachment, textarget, texture, level)
}

// Sets whether or not polygons are considered front-facing based
// on their winding direction.
func (c *Context) FrontFace(mode int) {
	c.Call("frontFace", mode)
}

// Creates a set of textures for a WebGLTexture object with image
// dimensions from the original size of the image down to a 1x1 image.
func (c *Context) GenerateMipmap(target int) {
	c.Call("generateMipmap", target)
}

// Returns an WebGLActiveInfo object containing the size, type, and name
// of a vertex attribute at a specific index position in a program object.
func (c *Context) GetActiveAttrib(program js.Value, index int) js.Value {
	return c.Call("getActiveAttrib", program, index)
}

// Returns an WebGLActiveInfo object containing the size, type, and name
// of a uniform attribute at a specific index position in a program object.
func (c *Context) GetActiveUniform(program js.Value, index int) js.Value {
	return c.Call("getActiveUniform", program, index)
}

// Returns a slice of WebGLShaders bound to a WebGLProgram.
func (c *Context) GetAttachedShaders(program js.Value) []js.Value {
	objs := c.Call("getAttachedShaders", program)
	shaders := make([]js.Value, objs.Length())
	for i := 0; i < objs.Length(); i++ {
		shaders[i] = objs.Index(i)
	}
	return shaders
}

func (c *Context) CreateVertexArray() js.Value {
	return c.Call("createVertexArray")
}
func (c *Context) BindVertexArray(vaoId js.Value) {
	c.Call("bindVertexArray", vaoId)
}
func (c *Context) DeleteVertexArray(vaoId js.Value) {
	c.Call("deleteVertexArray", vaoId)
}

// Returns an index to the location in a program of a named attribute variable.
func (c *Context) GetAttribLocation(program js.Value, name string) int {
	return c.Call("getAttribLocation", program, name).Int()
}

// TODO: Create type specific variations.
// Returns the type of a parameter for a given buffer.
func (c *Context) GetBufferParameter(target, pname int) js.Value {
	return c.Call("getBufferParameter", target, pname)
}

// TODO: Create type specific variations.
// Returns the natural type value for a constant parameter.
func (c *Context) GetParameter(pname js.Value) js.Value {
	return c.Call("getParameter", pname)
}

// Returns a value for the WebGL error flag and clears the flag.
func (c *Context) GetError() int {
	return c.Call("getError").Int()
}

// TODO: Create type specific variations.
// Enables a passed extension, otherwise returns null.
func (c *Context) GetExtension(name string) js.Value {
	return c.Call("getExtension", name)
}

// TODO: Create type specific variations.
// Gets a parameter value for a given target and attachment.
func (c *Context) GetFramebufferAttachmentParameter(target, attachment, pname int) js.Value {
	return c.Call("getFramebufferAttachmentParameter", target, attachment, pname)
}

// Returns the value of the program parameter that corresponds to a supplied pname
// which is interpreted as an int.
func (c *Context) GetProgramParameteri(program js.Value, pname int) int {
	return c.Call("getProgramParameter", program, pname).Int()
}

// Returns the value of the program parameter that corresponds to a supplied pname
// which is interpreted as a bool.
func (c *Context) GetProgramParameterb(program js.Value, pname int) bool {
	return c.Call("getProgramParameter", program, pname).Bool()
}

// Returns information about the last error that occurred during
// the failed linking or validation of a WebGL program object.
func (c *Context) GetProgramInfoLog(program js.Value) string {
	return c.Call("getProgramInfoLog", program).String()
}

// TODO: Create type specific variations.
// Returns a renderbuffer parameter from the currently bound WebGLRenderbuffer object.
func (c *Context) GetRenderbufferParameter(target, pname int) js.Value {
	return c.Call("getRenderbufferParameter", target, pname)
}

// TODO: Create type specific variations.
// Returns the value of the parameter associated with pname for a shader object.
func (c *Context) GetShaderParameter(shader js.Value, pname int) js.Value {
	return c.Call("getShaderParameter", shader, pname)
}

// Returns the value of the parameter associated with pname for a shader object.
func (c *Context) GetShaderParameterb(shader js.Value, pname int) bool {
	return c.Call("getShaderParameter", shader, pname).Bool()
}

// Returns errors which occur when compiling a shader.
func (c *Context) GetShaderInfoLog(shader js.Value) string {
	return c.Call("getShaderInfoLog", shader).String()
}

// Returns source code string associated with a shader object.
func (c *Context) GetShaderSource(shader js.Value) string {
	return c.Call("getShaderSource", shader).String()
}

// Returns a slice of supported extension strings.
func (c *Context) GetSupportedExtensions() []string {
	ext := c.Call("getSupportedExtensions")
	extensions := make([]string, ext.Length())
	for i := 0; i < ext.Length(); i++ {
		extensions[i] = ext.Index(i).String()
	}
	return extensions
}

// TODO: Create type specific variations.
// Returns the value for a parameter on an active texture unit.
func (c *Context) GetTexParameter(target, pname int) js.Value {
	return c.Call("getTexParameter", target, pname)
}

// TODO: Create type specific variations.
// Gets the uniform value for a specific location in a program.
func (c *Context) GetUniform(program, location js.Value) js.Value {
	return c.Call("getUniform", program, location)
}

// Returns a WebGLUniformLocation object for the location
// of a uniform variable within a WebGLProgram object.
func (c *Context) GetUniformLocation(program js.Value, name string) js.Value {
	return c.Call("getUniformLocation", program, name)
}

// TODO: Create type specific variations.
// Returns data for a particular characteristic of a vertex
// attribute at an index in a vertex attribute array.
func (c *Context) GetVertexAttrib(index, pname int) js.Value {
	return c.Call("getVertexAttrib", index, pname)
}

// Returns the address of a specified vertex attribute.
func (c *Context) GetVertexAttribOffset(index, pname int) int {
	return c.Call("getVertexAttribOffset", index, pname).Int()
}

// public function hint(target:GLenum, mode:GLenum) : Void;

// Returns true if buffer is valid, false otherwise.
func (c *Context) IsBuffer(buffer js.Value) bool {
	return c.Call("isBuffer", buffer).Bool()
}

// Returns whether the WebGL context has been lost.
func (c *Context) IsContextLost() bool {
	return c.Call("isContextLost").Bool()
}

// Returns true if buffer is valid, false otherwise.
func (c *Context) IsFramebuffer(framebuffer js.Value) bool {
	return c.Call("isFramebuffer", framebuffer).Bool()
}

// Returns true if program object is valid, false otherwise.
func (c *Context) IsProgram(program js.Value) bool {
	return c.Call("isProgram", program).Bool()
}

// Returns true if buffer is valid, false otherwise.
func (c *Context) IsRenderbuffer(renderbuffer js.Value) bool {
	return c.Call("isRenderbuffer", renderbuffer).Bool()
}

// Returns true if shader is valid, false otherwise.
func (c *Context) IsShader(shader js.Value) bool {
	return c.Call("isShader", shader).Bool()
}

// Returns true if texture is valid, false otherwise.
func (c *Context) IsTexture(texture js.Value) bool {
	return c.Call("isTexture", texture).Bool()
}

// Returns whether or not a WebGL capability is enabled for this context.
func (c *Context) IsEnabled(capability int) bool {
	return c.Call("isEnabled", capability).Bool()
}

// Sets the width of lines in WebGL.
func (c *Context) LineWidth(width float64) {
	c.Call("lineWidth", width)
}

// Links an attached vertex shader and an attached fragment shader
// to a program so it can be used by the graphics processing unit (GPU).
func (c *Context) LinkProgram(program js.Value) {
	c.Call("linkProgram", program)
}

// Sets pixel storage modes for readPixels and unpacking of textures
// with texImage2D and texSubImage2D.
func (c *Context) PixelStorei(pname, param int) {
	c.Call("pixelStorei", pname, param)
}

// Sets the implementation-specific units and scale factor
// used to calculate fragment depth values.
func (c *Context) PolygonOffset(factor, units float64) {
	c.Call("polygonOffset", factor, units)
}

// TODO: Figure out if pixels should be a slice.
// Reads pixel data into an ArrayBufferView object from a
// rectangular area in the color buffer of the active frame buffer.
func (c *Context) ReadPixels(x, y, width, height, format, typ int, pixels js.Value) {
	c.Call("readPixels", x, y, width, height, format, typ, pixels)
}

// Creates or replaces the data store for the currently bound WebGLRenderbuffer object.
func (c *Context) RenderbufferStorage(target, internalFormat, width, height int) {
	c.Call("renderbufferStorage", target, internalFormat, width, height)
}

//func (c *Context) SampleCoverage(value float64, invert bool) {
//	c.Call("sampleCoverage", value, invert)
//}

// Sets the dimensions of the scissor box.
func (c *Context) Scissor(x, y, width, height int) {
	c.Call("scissor", x, y, width, height)
}

// Sets and replaces shader source code in a shader object.
func (c *Context) ShaderSource(shader js.Value, source string) {
	c.Call("shaderSource", shader, source)
}

// public function stencilFunc(func:GLenum, ref:GLint, mask:GLuint) : Void;
// public function stencilFuncSeparate(face:GLenum, func:GLenum, ref:GLint, mask:GLuint) : Void;
// public function stencilMask(mask:GLuint) : Void;
// public function stencilMaskSeparate(face:GLenum, mask:GLuint) : Void;
// public function stencilOp(fail:GLenum, zfail:GLenum, zpass:GLenum) : Void;
// public function stencilOpSeparate(face:GLenum, fail:GLenum, zfail:GLenum, zpass:GLenum) : Void;

// Loads the supplied pixel data into a texture.
func (c *Context) TexImage2D(target, level, internalFormat, format, kind int, image js.Value) {
	c.Call("texImage2D", target, level, internalFormat, format, kind, image)
}

// Sets texture parameters for the current texture unit.
func (c *Context) TexParameteri(target int, pname int, param int) {
	c.Call("texParameteri", target, pname, param)
}

// Replaces a portion of an existing 2D texture image with all of another image.
func (c *Context) TexSubImage2D(target, level, xoffset, yoffset, format, typ int, image js.Value) {
	c.Call("texSubImage2D", target, level, xoffset, yoffset, format, typ, image)
}

// Assigns a floating point value to a uniform variable for the current program object.
func (c *Context) Uniform1f(location js.Value, x float32) {
	c.Call("uniform1f", location, x)
}

// Assigns a integer value to a uniform variable for the current program object.
func (c *Context) Uniform1i(location js.Value, x int) {
	c.Call("uniform1i", location, x)
}

// Assigns 2 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform2f(location js.Value, x, y float32) {
	c.Call("uniform2f", location, x, y)
}

// Assigns 2 integer values to a uniform variable for the current program object.
func (c *Context) Uniform2i(location js.Value, x, y int) {
	c.Call("uniform2i", location, x, y)
}

// Assigns 3 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform3f(location js.Value, x, y, z float32) {
	c.Call("uniform3f", location, x, y, z)
}

// Assigns 3 integer values to a uniform variable for the current program object.
func (c *Context) Uniform3i(location js.Value, x, y, z int) {
	c.Call("uniform3i", location, x, y, z)
}

// Assigns 4 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform4f(location js.Value, x, y, z, w float32) {
	c.Call("uniform4f", location, x, y, z, w)
}

// Assigns 4 integer values to a uniform variable for the current program object.
func (c *Context) Uniform4i(location js.Value, x, y, z, w int) {
	c.Call("uniform4i", location, x, y, z, w)
}

// public function uniform1fv(location:WebGLUniformLocation, v:ArrayAccess<Float>) : Void;
// public function uniform1iv(location:WebGLUniformLocation, v:ArrayAccess<Long>) : Void;
// public function uniform2fv(location:WebGLUniformLocation, v:ArrayAccess<Float>) : Void;
// public function uniform2iv(location:WebGLUniformLocation, v:ArrayAccess<Long>) : Void;
// public function uniform3fv(location:WebGLUniformLocation, v:ArrayAccess<Float>) : Void;
// public function uniform3iv(location:WebGLUniformLocation, v:ArrayAccess<Long>) : Void;
// public function uniform4fv(location:WebGLUniformLocation, v:ArrayAccess<Float>) : Void;
// public function uniform4iv(location:WebGLUniformLocation, v:ArrayAccess<Long>) : Void;

// Sets values for a 2x2 floating point vector matrix into a
// uniform location as a matrix or a matrix array.
func (c *Context) UniformMatrix2fv(location js.Value, transpose bool, value []float32) {
	c.Call("uniformMatrix2fv", location, transpose, value)
}

// Sets values for a 3x3 floating point vector matrix into a
// uniform location as a matrix or a matrix array.
func (c *Context) UniformMatrix3fv(location js.Value, transpose bool, value []float32) {
	c.Call("uniformMatrix3fv", location, transpose, value)
}

// Sets values for a 4x4 floating point vector matrix into a
// uniform location as a matrix or a matrix array.
func (c *Context) UniformMatrix4fv(location js.Value, transpose bool, value js.Value) {
	c.Call("uniformMatrix4fv", location, transpose, value)
}

// Set the program object to use for rendering.
func (c *Context) UseProgram(program js.Value) {
	c.Call("useProgram", program)
}

// Returns whether a given program can run in the current WebGL state.
func (c *Context) ValidateProgram(program js.Value) {
	c.Call("validateProgram", program)
}

func (c *Context) VertexAttribPointer(index, size int, typ js.Value, normal bool, stride int, offset int) {
	c.Call("vertexAttribPointer", index, size, typ, normal, stride, offset)
}

// public function vertexAttrib1f(indx:GLuint, x:GLfloat) : Void;
// public function vertexAttrib2f(indx:GLuint, x:GLfloat, y:GLfloat) : Void;
// public function vertexAttrib3f(indx:GLuint, x:GLfloat, y:GLfloat, z:GLfloat) : Void;
// public function vertexAttrib4f(indx:GLuint, x:GLfloat, y:GLfloat, z:GLfloat, w:GLfloat) : Void;
// public function vertexAttrib1fv(indx:GLuint, values:ArrayAccess<Float>) : Void;
// public function vertexAttrib2fv(indx:GLuint, values:ArrayAccess<Float>) : Void;
// public function vertexAttrib3fv(indx:GLuint, values:ArrayAccess<Float>) : Void;
// public function vertexAttrib4fv(indx:GLuint, values:ArrayAccess<Float>) : Void;

// Represents a rectangular viewable area that contains
// the rendering results of the drawing buffer.
func (c *Context) Viewport(x, y, width, height int) {
	c.Call("viewport", x, y, width, height)
}
