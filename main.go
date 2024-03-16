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
	if rl.IsKeyPressed(rl.KeyR) {
		gs.ground = ground.New()
		gs.player = player.New()
	}
	if gs.player.Dead {
		return
	}

	gs.ground.Update()
	gs.player.Update(gs.ground)
}

func (gs *gameState) draw() {
	rl.ClearBackground(rl.SkyBlue)

	rl.DrawFPS(10, 10)

	if gs.player.Dead {
		gs.drawGameOver()
		return
	}

	gs.ground.Draw()
	gs.player.Draw()
}

func (gs *gameState) drawGameOver() {
	w, h := rl.GetScreenWidth(), rl.GetScreenHeight()
	deadText := "You are dead"
	textWidth := rl.MeasureText(deadText, 60)
	deadTextX := int32(w/2) - textWidth/2
	deadTextY := int32(h / 2)
	rl.DrawText(deadText, deadTextX, deadTextY, 60, rl.Red)

	restartText := "Press R to restart"
	textWidth = rl.MeasureText(restartText, 20)
	restartTextX := int32(w/2) - textWidth/2
	restartTextY := deadTextY + 60
	rl.DrawText(restartText, restartTextX, restartTextY, 20, rl.Black)
}
