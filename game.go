package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pechorka/bread-game-jam/pkg/random"
	"github.com/pechorka/bread-game-jam/pkg/rlutils"
)

const (
	groundHeightPercent float32 = 0.2
)

const (
	holeSpawnRate float32 = 1
	holeSpeed     float32 = 10
	maxHoleCount  int     = 10

	holeMinWidthScreenPercent float32 = 0.05
	holeMaxWidthScreenPercent float32 = 0.1
)

const (
	platformSpawnRate float32 = 2
	platformSpeed     float32 = 10
	maxPlatformCount  int     = 10

	platformHeightScreenPercent    float32 = 0.02
	platformWidthScreenPercentFrom float32 = 0.2
	platformWidthScreenPercentTo   float32 = 0.3
)

const (
	playerMaxJumpHeightScreenPercent float32 = 0.2
	playerHeightScreenPercent        float32 = 0.1
	playerLeftMarginScreenPercent    float32 = 0.2

	playerInitialVerticalSpeed float32 = 3
)

type gameState struct {
	screen gameScreen

	commonState commonState
	ground      ground
	platforms   platforms
	player      player
}

func newGameState() *gameState {
	return &gameState{
		screen:    gameScreenGame,
		ground:    newGround(),
		platforms: newPlatforms(),
		player:    newPlayer(),
	}
}

type gameScreen string

const (
	gameScreenMenu     gameScreen = "menu"
	gameScreenGame     gameScreen = "game"
	gameScreenGameOver gameScreen = "game_over"
	gameScreenWin      gameScreen = "win"
)

func (gs *gameState) update() {
	if rl.IsKeyPressed(rl.KeyR) {
		*gs = *newGameState()
		return
	}

	gs.commonState.update()
	gs.ground.update(gs.commonState)
	gs.platforms.updatePlatforms(gs.commonState, gs.ground.border)
	gs.player.update(gs.commonState, gs.ground.border, gs.platforms.borders, gs.ground.holes.borders)
	if gs.player.dead {
		gs.screen = gameScreenGameOver
	}
}

func (gs *gameState) draw() {
	rl.ClearBackground(rl.SkyBlue)

	rl.DrawFPS(10, 10)

	switch gs.screen {
	case gameScreenMenu:
		panic("not implemented")
	case gameScreenGame:
		gs.ground.draw()
		gs.platforms.draw()
		gs.player.draw()
	case gameScreenGameOver:
		gs.drawGameOver()
	case gameScreenWin:
		panic("not implemented")
	}
}

type commonState struct {
	screenWidth  float32
	screenHeight float32
	dt           float32
}

func (cs *commonState) update() {
	cs.screenWidth, cs.screenHeight = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	cs.dt = rl.GetFrameTime()
}

type ground struct {
	border rl.Rectangle
	holes  holes
}

func newGround() ground {
	return ground{
		holes: newHoles(),
	}
}

func (g *ground) update(cs commonState) {
	g.updateBorders(cs)
	g.holes.updateHoles(cs, g.border)
}

func (g *ground) updateBorders(cs commonState) {
	groundHeight := cs.screenHeight * groundHeightPercent
	g.border = rl.Rectangle{
		X:      0,
		Y:      cs.screenHeight - groundHeight,
		Width:  cs.screenWidth,
		Height: groundHeight,
	}
}

func (g *ground) draw() {
	rl.DrawRectangleRec(g.border, rl.Brown)
	g.holes.draw()
}

type holes struct {
	spawnTimer float32
	borders    []rl.Rectangle
}

func newHoles() holes {
	return holes{
		borders: make([]rl.Rectangle, maxHoleCount),
	}
}

func (h *holes) updateHoles(cs commonState, groundBorders rl.Rectangle) {
	h.spawnTimer += cs.dt

	for i := range h.borders {
		holeBorder := &h.borders[i]
		holeBorder.X -= holeSpeed

		holeNotVisible := holeBorder.X+holeBorder.Width < 0

		if holeNotVisible && h.spawnTimer >= holeSpawnRate {
			h.spawnTimer = 0

			holeMinWidth := cs.screenWidth * holeMinWidthScreenPercent
			holeMaxWidth := cs.screenWidth * holeMaxWidthScreenPercent
			holeWidth := random.Float32(holeMinWidth, holeMaxWidth)

			holeBorder.X = cs.screenWidth + holeWidth
			holeBorder.Y = groundBorders.Y
			holeBorder.Width = holeWidth
			holeBorder.Height = groundBorders.Height
		}
	}
}

func (h *holes) draw() {
	for i := range h.borders {
		rl.DrawRectangleRec(h.borders[i], rl.Black)
	}
}

type platforms struct {
	spawnTimer float32
	borders    []rl.Rectangle
}

func newPlatforms() platforms {
	return platforms{
		borders: make([]rl.Rectangle, maxPlatformCount),
	}
}

func (p *platforms) updatePlatforms(cs commonState, groundBorders rl.Rectangle) {
	p.spawnTimer += cs.dt

	for i := range p.borders {
		platformBorder := &p.borders[i]
		platformBorder.X -= platformSpeed

		platformNotVisible := platformBorder.X+platformBorder.Width < 0

		if platformNotVisible && p.spawnTimer >= platformSpawnRate {
			p.spawnTimer = 0

			platformMinWidth := cs.screenWidth * platformWidthScreenPercentFrom
			platformMaxWidth := cs.screenWidth * platformWidthScreenPercentTo
			platformWidth := random.Float32(platformMinWidth, platformMaxWidth)

			playerMaxJumpHeight := cs.screenHeight * playerMaxJumpHeightScreenPercent
			playerHeight := cs.screenHeight * playerHeightScreenPercent
			platformHeight := cs.screenHeight * platformHeightScreenPercent

			platformYFrom := groundBorders.Y - playerMaxJumpHeight + platformHeight
			platformYTo := groundBorders.Y - platformHeight - playerHeight
			platformY := random.Float32(platformYFrom, platformYTo)

			platformBorder.X = cs.screenWidth + platformWidth
			platformBorder.Y = platformY
			platformBorder.Width = platformWidth
			platformBorder.Height = platformHeight
		}
	}
}

func (p *platforms) draw() {
	for i := range p.borders {
		rl.DrawRectangleRec(p.borders[i], rl.Gray)
	}
}

type player struct {
	verticalPosition float32
	verticalSpeed    float32
	// can't rely on verticalPosition == 0 because it will > 0 while player is decending
	jumping    bool
	onPlatform bool
	border     rl.Rectangle
	dead       bool
}

func newPlayer() player {
	return player{
		verticalSpeed: playerInitialVerticalSpeed,
	}
}

func (p *player) update(cs commonState, groundBorders rl.Rectangle, platformBorders, holeBorders []rl.Rectangle) {
	p.updateBorder(cs, groundBorders)
	p.updateVerticalPosition(cs, platformBorders)
	if p.verticalPosition == 0 {
		for _, holeBorder := range holeBorders {
			if rlutils.VerticalCollision(p.border, holeBorder) {
				p.dead = true
				break
			}
		}
	}
}

func (p *player) updateBorder(cs commonState, groundBorders rl.Rectangle) {
	playerHeight := cs.screenHeight * playerHeightScreenPercent
	playerLeftMargin := cs.screenWidth * playerLeftMarginScreenPercent
	p.border = rl.Rectangle{
		X:      playerLeftMargin,
		Y:      groundBorders.Y - playerHeight - p.verticalPosition,
		Width:  playerHeight,
		Height: playerHeight,
	}
}

func (p *player) updateVerticalPosition(
	cs commonState,
	platformBorders []rl.Rectangle,
) {
	p.onPlatform = false
	for _, platformBorder := range platformBorders {
		if rl.CheckCollisionRecs(p.border, platformBorder) {
			// player above platform
			if p.border.Y+p.border.Height < platformBorder.Y+platformBorder.Height {
				p.onPlatform = true
			}
			p.jumping = false

			break
		}
	}

	isSpacePressed := rl.IsKeyDown(rl.KeySpace)

	if isSpacePressed &&
		((!p.jumping && p.verticalPosition == 0) || // player on the ground and space pressed -> jump
			p.onPlatform) { // player on the platform and space pressed -> jump
		p.jumping = true
	}
	// player is jumping and space released -> start descending
	if !isSpacePressed && p.jumping {
		p.jumping = false
	}

	if p.jumping {
		p.verticalPosition += p.verticalSpeed
	} else if !p.onPlatform {
		p.verticalPosition -= p.verticalSpeed
	}

	maxJumpHeight := cs.screenHeight * playerMaxJumpHeightScreenPercent
	p.verticalPosition = rl.Clamp(p.verticalPosition, 0, maxJumpHeight)

	// player reached max jump height -> start descending
	if p.jumping && p.verticalPosition >= maxJumpHeight {
		p.jumping = false
	}
}

func (p *player) draw() {
	rl.DrawRectangleRec(p.border, rl.Yellow)
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
