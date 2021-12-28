package main

import (
	"encoding/json"
	"image"
	"image/draw"
	"io/ioutil"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

func Draw_Text(game *Game) {
	gl.UseProgram(game.drawText.program)

	if textRendered == false {
		var vertices = make([]float32, (len("Dildosers Test{69}") * 16))
		Text_Render_Text(vertices, "Dildosers Test{69}")
		textLength = len(vertices) / 16
		gl.BindVertexArray(game.drawText.VAO)
		gl.BindBuffer(gl.ARRAY_BUFFER, game.drawText.VAO)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
		textRendered = true
	}

	//Use Program

	//Blend Enable
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//Binding textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, game.drawText.image)

	//Binding VAO and applying subdata
	gl.BindVertexArray(game.drawText.VAO)
	//gl.BindBuffer(gl.ARRAY_BUFFER, game.drawText.VAO)
	//gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))

	i := 0
	for ; i <= textLength; i += 1 {
		gl.DrawArrays(gl.TRIANGLE_STRIP, int32(i*4), 4)
	}

	//vertices = nil
	//Unbiding Everything
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
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

}

func Create_Dynamic_Vao_Text(_program uint32, _size int) (uint32, uint32) {

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
	//	Location 0 vec4 vec2 positon & vec2 texture position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, nil)

	return vao, vbo
}

//Loading image from file
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

//-----JSON FOR TEXT ------

type Glyph struct {
	x        float32
	y        float32
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	W        float32 `json:"width"`
	H        float32 `json:"height"`
	Xadvance float32 `json:"xadvance"`
	Id       byte    `json:"id"`
}
type Glyphs struct {
	Glyph []Glyph `json:"chars"`
}

func Text_Json_Parsing(_jsonFile string) map[byte]Glyph {
	//todo: take out const of img size
	//Parses custom-msdf.json for information to how each charater is mapped inside the custom.png
	jsonFile, err := os.Open(_jsonFile)
	if err != nil {
		println("Problem opening image:")
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Glyphs array
	var glyph Glyphs

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'glyphs' which we defined above
	err = json.Unmarshal(byteValue, &glyph)

	if err != nil {
		println("Problem Unmarshaling data:")
		panic(err)
	}
	mapping := make(map[byte]Glyph)
	//Change the constants with actuall data from image like the 512
	for i := 0; i < len(glyph.Glyph); i++ {
		if glyph.Glyph[i].X == 0 {
			glyph.Glyph[i].x = 0
		} else {
			glyph.Glyph[i].x = glyph.Glyph[i].X / 512 //Divid by png size
		}

		if glyph.Glyph[i].Y == 0 {
			glyph.Glyph[i].y = 0
		} else {
			glyph.Glyph[i].y = glyph.Glyph[i].Y / 512
		}

		glyph.Glyph[i].X = (glyph.Glyph[i].X + glyph.Glyph[i].W) / 512
		glyph.Glyph[i].Y = (glyph.Glyph[i].Y + glyph.Glyph[i].H) / 512
		mapping[glyph.Glyph[i].Id] = glyph.Glyph[i]
	}
	return mapping
	/*for _, v := range mapping {
		fmt.Println("x: ", v.x)
		fmt.Println("y: ", v.y)
		fmt.Println("X: ", v.X)
		fmt.Println("Y: ", v.Y)
		fmt.Println("width: ", v.W)
		fmt.Println("height: ", v.H)
		fmt.Println("Id: ", v.Id)
	}*/
	/* //Debugging print
	for i := 0; i < len(glyph.Glyph); i++ {
		fmt.Println("x: ", glyph.Glyph[i].x)
		fmt.Println("y: ", glyph.Glyph[i].y)
		fmt.Println("X: ", glyph.Glyph[i].X)
		fmt.Println("Y: ", glyph.Glyph[i].Y)
		fmt.Println("width: ", glyph.Glyph[i].W)
		fmt.Println("height: ", glyph.Glyph[i].H)
		fmt.Println("Id: ", glyph.Glyph[i].Id)
	}
	*/
}

func Text_Render_Text(vert []float32, _sentence string) {
	var sx, sy int
	sx = 10
	sy = 10

	for i, v := range _sentence {
		if byte(v) == 32 {
			sx += 16
		} else {
			g := mapping[byte(v)]
			x := float32(sx)
			y := float32(sy)
			X := g.W + x
			Y := g.H + y
			index := (i * 16)
			copy(vert[index:index+16], []float32{
				x, y, g.x, g.Y, // left-bottom
				x, Y, g.x, g.y, // left-top
				X, y, g.X, g.Y, // right-bottom
				X, Y, g.X, g.y, // right-top
			})
			sx += int(g.Xadvance)
		}
	}
}
