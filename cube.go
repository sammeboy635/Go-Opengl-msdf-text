package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Vec3 struct {
	x float32
	y float32
	z float32
}
type Cube struct {
	p1 Vec3
	p2 Vec3
	p3 Vec3
	p4 Vec3
}

const (
	size = 0.5
)

func Draw_Cube(game *Game) {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(game.drawBlock.program)

	cube := New_Cube(-0.5, -0.5, 0)
	cube = append(cube, New_Cube(0.0, 0.0, 0.0)...)
	cube = append(cube, New_Cube(0.5, 0.5, 0)...)
	//gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, game.drawBlock.VAO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(cube)*4, gl.Ptr(cube))
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
	//gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 8)
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 12, 3) //Count is the number of points | Instancecount is the number of points to draw
	gl.UseProgram(0)
	glfw.PollEvents()
	game.win.SwapBuffers()
}

func New_Create_DrawData_Block() DrawData {
	var drawData DrawData
	drawData.program = Create_Program("shader/frag.shadder", "shader/vert.shadder")
	drawData.VAO, drawData.VBO = Create_Dynamic_Vao_Block(1024)
	drawData.VAOSize = 1024

	return drawData
}

func Create_Dynamic_Vao_Block(_size int) (uint32, uint32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, _size, nil, gl.DYNAMIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao, vbo
}

func New_Cube(x float32, y float32, z float32) []float32 {
	X := x + size
	Y := y + size

	cube := []float32{
		x, Y, z, //top left
		x, y, z, //bottum left
		X, Y, z, //top right
		X, y, z, //bottum right
	}
	//Print_Cube(cube)

	return cube
}
func New_Triangle(x float32, y float32, z float32) []float32 {
	X := x + size
	Y := y + size
	triangle := []float32{
		x, Y, z,
		x, y, z,
		X, Y, z,
	}
	return triangle
}
func Print_Cube(cube []float32) {
	for c := range cube {
		fmt.Printf("%v", c)
	}
}
