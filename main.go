package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pechorka/bread-game-jam/internal/ground"
	"github.com/pechorka/bread-game-jam/internal/player"
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	const factor = 100
	rl.InitWindow(16*factor, 9*factor, "Bread game jam")
	rl.SetTargetFPS(60)

	gs := gameState{
		ground: ground.New(),
		player: player.New(),
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		gs.update()
		gs.draw()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

type gameState struct {
	ground *ground.Ground
	player *player.Player
}

func (gs *gameState) update() {
	gs.ground.Update()
	gs.player.Update(gs.ground.Rect)
}

func (gs *gameState) draw() {
	rl.ClearBackground(rl.Black)

	gs.ground.Draw()
	gs.player.Draw()
}
