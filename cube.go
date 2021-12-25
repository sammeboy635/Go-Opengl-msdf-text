package main

import (
	"fmt"
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

func New_Cube(x float32, y float32, z float32) []float32 {
	X := x + size
	Y := y + size

	cube := []float32{
		x, Y, z,
		x, y, z,
		X, Y, z,
		X, y, z,
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
