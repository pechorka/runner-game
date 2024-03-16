package rlutils

import rl "github.com/gen2brain/raylib-go/raylib"

func GetScreenDimensions() (float32, float32) {
	return float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
}
