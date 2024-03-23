package rlutils

import rl "github.com/gen2brain/raylib-go/raylib"

func VerticalCollision(r, in rl.Rectangle) bool {
	// 					in.X ....................... in.X+in.Width
	//  				|                           |
	//  				|                           |
	//  				|                           |
	//  				|                           |
	// 					r.X ....................... r.X+r.Width
	// or
	//         r.X ....................... r.X+r.Width
	// or
	// 	                               r.X ....................... r.X+r.Width

	// left side of r is inside in
	if in.X <= r.X && r.X <= in.X+in.Width {
		return true
	}

	// right side of r is inside in
	if in.X <= r.X+r.Width && r.X+r.Width <= in.X+in.Width {
		return true
	}

	return false
}

func SetToZero(r *rl.Rectangle) {
	r.X = 0
	r.Y = 0
	r.Width = 0
	r.Height = 0
}
