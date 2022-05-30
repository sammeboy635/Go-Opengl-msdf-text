package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)
var abcs = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var whatTime = "Time is "
type TextRender struct {
	shader  DrawData
	charMap map[byte]Glyph
}

func (t *TextRender) Draw_Text() {
	//Use Program
	gl.UseProgram(t.shader.program)

	if textRendered == false {
		newtime := fmt.Sprintf("%s%s",whatTime, time.Now().String())
		var vertices = make([]float32, (len(newtime) * 16))
		t.Render_Text(vertices, newtime)

		textLength = len(abcs) // len(vertices) / 16 
		gl.BindVertexArray(t.shader.VAO)
		gl.BindBuffer(gl.ARRAY_BUFFER, t.shader.VAO)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
		//textRendered = true
	}

	//Blend Enable
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//Binding textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.shader.image)

	//Binding VAO and applying subdata
	gl.BindVertexArray(t.shader.VAO)
	//gl.BindBuffer(gl.ARRAY_BUFFER, game.drawText.VAO)
	//gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))

	
	first := make([]int32,textLength)
	count := make([]int32,textLength)
	for i := 0; i < textLength; i++ {
		first[i] = int32(i) * 4
		count[i] = 4
	}
	gl.MultiDrawArrays(gl.TRIANGLE_STRIP, &first[0], &count[0], int32(textLength))

	// i := 0
	// for ; i <= textLength; i += 1 {
	// 	gl.DrawArrays(gl.TRIANGLE_STRIP, int32(i*4), 4)
	// }
	// first := []int32{0, 4, 8, 12}
	// count := []int32{4, 4, 4, 4}
	//num := int32(4)
	

	//vertices = nil
	//Unbiding Everything
	gl.UseProgram(0)
	gl.Disable(gl.BLEND)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
func (t *TextRender) Init() {
	t.shader.Create_Program("shader/textFrag.shadder", "shader/textVert.shadder")
	t.shader.Load_Image("custom-msdf/custom.png", gl.TEXTURE0)
	t.Set_Program_Matric()
	t.Create_Dynamic_VAO(4096)

	t.charMap = Text_Json_Parsing("custom-msdf/custom-msdf.json")
}

func (t *TextRender) Create_Dynamic_VAO(_size int) {
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

	t.shader.VAO = vao
	t.shader.VBO = vbo
	t.shader.VAOSize = _size
}

func (t *TextRender) Set_Program_Matric() {

	//Preparing for Projection Matrix
	prjCStr, free := gl.Strs("projection") //Needs a free called after
	defer free()
	glProjectionLocation := gl.GetUniformLocation(t.shader.program, *prjCStr)
	projection := mgl.Ortho2D(0, float32(win_width), 0.0, float32(win_height)) //Create a Ortho2d projection for

	gl.UseProgram(t.shader.program)                                     //Bind program to set uninform in GPU
	gl.UniformMatrix4fv(glProjectionLocation, 1, false, &projection[0]) //Setting Projections

}

//-----JSON FOR TEXT ------

type Glyph struct {
	x        float32
	y        float32
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	OffY     float32 `json:"offy"`
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

}
func (g *Glyph) print() {
	println("x: ", g.x)
	println("y: ", g.y)
	println("X: ", g.X)
	println("Y: ", g.Y)
	println("width: ", g.W)
	println("height: ", g.H)
	println("Id: ", g.Id)
}

func (t *TextRender) Render_Text(vert []float32, _sentence string) {
	var scale float32 = 0.5
	var sx, sy float32
	sx = 10
	sy = 10

	for i, v := range _sentence {
		if byte(v) == 32{ // Space so move forward
			sx += 16 * scale // TODO * by scale
			
			continue
		}
		g := t.charMap[byte(v)]
		x := sx
		y := (sy + (g.OffY * scale))
		X := ((g.W * scale) + x) 
		Y := ((g.H * scale) + y) 
		index := (i * 16) 
		copy(vert[index:index+16], []float32{
			x, y, g.x, g.Y, // left-bottom (pos, pos, tex, tex)
			x, Y, g.x, g.y, // left-top
			X, y, g.X, g.Y, // right-bottom
			X, Y, g.X, g.y, // right-top
		})
		sx += g.Xadvance * scale
		
	}
}
