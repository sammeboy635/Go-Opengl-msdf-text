package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Glyph struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}
type Glyphs struct {
	Glyph []Glyph `json:"chars"`
}

func Draw_Text(game *Game) {
	x := float32(0.0)
	y := float32(0.0)
	X := float32(0.017578125)
	Y := float32(0.095703125)
	x1 := float32(0)
	X1 := float32(16)
	y1 := float32(0)
	Y1 := float32(96)
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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, game.drawText.image)

	gl.BindBuffer(gl.ARRAY_BUFFER, game.drawText.VAO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))

	//gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, 1) //Count is the number of points | Instancecount is the number of points to draw
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
	drawData.image = loadImage("custom-msdf/custom.png")

	return drawData
}

func Set_Program_Matrix(_program uint32) {

	//Preparing for Projection Matrix
	prjCStr, free := gl.Strs("projection") //Needs a free called after
	defer free()
	glProjectionLocation := gl.GetUniformLocation(_program, *prjCStr)
	projection := mgl.Ortho2D(0, float32(width), 0.0, float32(height)) //Create a Ortho2d projection for

	gl.UseProgram(_program)                                             //Bind program to set uninform in GPU
	gl.UniformMatrix4fv(glProjectionLocation, 1, false, &projection[0]) //Setting Projections

	vertAttrib := uint32(gl.GetAttribLocation(_program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(_program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
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

	//gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

	return vao, vbo
}

func loadImage(_file string) uint32 {
	imgFile, err := os.Open(_file)
	if err != nil {
		println("Problem opening image:")
		panic(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		println("Problem decoding the image:")
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("incorret stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return loadTexture(rgba)
}

func loadTexture(_rgba *image.RGBA) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	//gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(_rgba.Rect.Size().X),
		int32(_rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(_rgba.Pix))

	return texture
}

func Text_Json_Parsing(_jsonFile string) {
	jsonFile, err := os.Open(_jsonFile)
	if err != nil {
		println("Problem opening image:")
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var glyph Glyphs

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &glyph)

	if err != nil {
		println("Problem Unmarshaling data:")
		panic(err)
	}
	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(glyph.Glyph); i++ {
		fmt.Println("x: ", glyph.Glyph[i].X)
		fmt.Println("y: ", glyph.Glyph[i].Y)
		fmt.Println("width: ", glyph.Glyph[i].Width)
		fmt.Println("height: ", glyph.Glyph[i].Height)
	}
}
