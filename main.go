package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	const factor = 100
	rl.InitWindow(16*factor, 9*factor, "Bread game jam")
	rl.SetTargetFPS(60)

	rl.InitAudioDevice()

	assets := loadAssets()
	gs := newGameState(assets)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		gs.update()
		gs.draw()

		rl.EndDrawing()
	}

	rl.CloseAudioDevice()
	assets.unload()
	rl.CloseWindow()
}
