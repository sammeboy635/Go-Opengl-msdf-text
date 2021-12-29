package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	X, Y, Z float32
	Rx, Ry  float32
	t       mgl.Vec3
}

func newPlayer() Camera {
	var p Camera
	p.X = 1
	p.Y = 2
	p.Z = 5
	return p
}
func (p *Camera) newLoc(_x float32, _y float32, _z float32) {
	p.X = _x
	p.Y = _y
	p.Z = _z
}
