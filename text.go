package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type ivec2 struct {
	x int
	y int
}

type Character struct {
	textureID uint16
	size      ivec2
	Bearing   ivec2
	advance   uint16
}

func Draw_Text(game *Game) {
	x := float32(0.0)
	y := float32(0.0)
	X := float32(1)
	Y := float32(1)
	x1 := float32(0)
	X1 := float32(250)
	y1 := float32(0)
	Y1 := float32(250)
	var vertices = []float32{
		// X, Y, U, V
		x1, y1, x, Y, // left-bottom
		x1, Y1, x, y, // left-top
		X1, y1, X, Y, // right-bottom
		X1, Y1, X, y, // right-top
	}
	gl.ClearColor(0.5, 1.0, 1.0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(game.drawText.program)
	//cube := New_Cube(0, 50, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, game.drawText.image)
	//gl.BindVertexArray(game.drawText.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, game.drawText.VAO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	//gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 0, 4) //Count is the number of points | Instancecount is the number of points to draw
	gl.UseProgram(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	glfw.PollEvents()
	game.win.SwapBuffers()
}

func New_Create_DrawData_Text() DrawData {
	var drawData DrawData
	drawData.program = Create_Program("shader/textFrag.shadder", "shader/textVert.shadder")
	drawData.VAO, drawData.VBO = Create_Dynamic_Vao_Text(drawData.program, 2048)
	drawData.VAOSize = 2048
	Set_Program_Matrix(drawData.program)
	drawData.image = loadImage("custom-msdf/square.png")

	return drawData
}
func Set_Program_Matrix(_program uint32) {
	prjCStr, free := gl.Strs("projection") //Needs a free called after
	defer free()
	glProjectionLocation := gl.GetUniformLocation(_program, *prjCStr)

	projection := mgl.Ortho2D(0, float32(width), 0.0, float32(height)) //Create a Ortho2d projection for

	gl.UseProgram(_program)                                             //Bind program to set projection in shader uniform
	gl.UniformMatrix4fv(glProjectionLocation, 1, false, &projection[0]) //Projection
	textureUniform := gl.GetUniformLocation(_program, gl.Str("text\x00"))
	gl.Uniform1i(textureUniform, 0)
	//gl.UseProgram(0) //Unbind program after.

}

func Create_Dynamic_Vao_Text(_program uint32, _size int) (uint32, uint32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, _size, nil, gl.DYNAMIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

	return vao, vbo
}
