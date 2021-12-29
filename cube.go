package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type CubeRender struct {
	shader        DrawData
	cameraUniform int32
	modelUniform  int32
}

func (c *CubeRender) Draw_Cube() {
	gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)

	gl.UseProgram(c.shader.program)

	//Blend Enable
	//gl.Enable(gl.BLEND)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	model := mgl.HomogRotate3D(float32(20), mgl.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(c.modelUniform, 1, false, &model[0])
	//Binding textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, c.shader.image)

	//Binding VAO and applying subdata
	gl.BindVertexArray(c.shader.VAO)
	//gl.BindBuffer(gl.ARRAY_BUFFER, c.shader.VAO)
	//gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(cubeVertices)*4, gl.Ptr(cubeVertices))

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

	//vertices = nil
	//Unbiding Everything
	gl.Disable(gl.DEPTH_TEST)
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
func (c *CubeRender) Init() {
	c.shader.Create_Program("shader/cube.fg", "shader/cube.vt")
	c.shader.Load_Image("custom-msdf/square.png", gl.TEXTURE0)
	c.Set_Program_Matric()
	c.Create_Dynamic_VAO(1024)
}
func (c *CubeRender) Create_Dynamic_VAO(_size int) {
	var vao, vbo uint32
	//Gen the buffers
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	//Bind the Buffers
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	//Set the buffer data
	//gl.BufferData(gl.ARRAY_BUFFER, _size, nil, gl.DYNAMIC_DRAW)

	//Set the attrib of the buffer array
	//	Location 0 vec4 vec2 positon & vec2 texture position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 5*4, 0)
	// 	Location 1 vec2 texcords
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4) //(location in shadder), (number of variables), (Normalize), (Size of this location(vec2 floats)), (offset in the array)

	c.shader.VAO = vao
	c.shader.VBO = vbo
	c.shader.VAOSize = _size
}
func (c *CubeRender) Set_Program_Matric() {
	gl.UseProgram(c.shader.program)
	projection := mgl.Perspective(mgl.DegToRad(45.0), float32(width)/height, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(c.shader.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl.LookAtV(mgl.Vec3{3, 3, 3}, mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})
	c.cameraUniform = gl.GetUniformLocation(c.shader.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(c.cameraUniform, 1, false, &camera[0])

	model := mgl.Ident4()
	c.modelUniform = gl.GetUniformLocation(c.shader.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(c.modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(c.shader.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
