package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type QuadRender struct {
	shader DrawData
}

const (
	size = 25
)

func (q *QuadRender) Draw_Quad() {

	cube := make([]float32, (4 * 12))

	New_Quad(cube, 0, 200, 200)
	New_Quad(cube, 1, 175, 175)

	first := []int32{0, 4}
	count := []int32{4, 4}
	num := int32(2) //int32(len(cube) / 12)
	gl.UseProgram(q.shader.program)

	gl.BindVertexArray(q.shader.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, q.shader.VAO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(cube)*4, gl.Ptr(cube))

	gl.MultiDrawArrays(gl.TRIANGLE_STRIP, &first[0], &count[0], num)

	gl.UseProgram(0)
	gl.BindVertexArray(0)
}

func (q *QuadRender) Init() {
	q.shader.Create_Program("shader/frag.shadder", "shader/vert.shadder")
	q.Create_Dynamic_VAO(1024)
	q.Set_Program_Matric()
}
func (q *QuadRender) Set_Program_Matric() {

	//Preparing for Projection Matrix
	prjCStr, free := gl.Strs("projection") //Needs a free called after
	defer free()
	glProjectionLocation := gl.GetUniformLocation(q.shader.program, *prjCStr)
	projection := mgl.Ortho2D(0, float32(win_width), 0.0, float32(win_height)) //Create a Ortho2d projection for

	gl.UseProgram(q.shader.program)                                     //Bind program to set uninform in GPU
	gl.UniformMatrix4fv(glProjectionLocation, 1, false, &projection[0]) //Setting Projections

}
func (q *QuadRender) Create_Dynamic_VAO(_size int) {
	var vao, vbo uint32
	//Gen the buffers
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	//Bind the Buffers
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	//Set the buffer data
	gl.BufferData(gl.ARRAY_BUFFER, _size, nil, gl.DYNAMIC_DRAW)

	//Set the attrib of the buffer array
	//	Location 0 vec3 position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)

	q.shader.VAO = vao
	q.shader.VBO = vbo
	q.shader.VAOSize = _size
}

func New_Quad(cube []float32, i int, x float32, y float32) {
	X := x + size
	Y := y + size
	index := (i * 12)
	copy(cube[index:index+12], []float32{
		x, Y, 0, //top left
		x, y, 0, //bottum left
		X, Y, 0, //top right
		X, y, 0, //bottum right
	})
	//Print_Cube(cube)

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
