package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pechorka/bread-game-jam/pkg/random"
	"github.com/pechorka/bread-game-jam/pkg/rlutils"
)

// TODO: spawn ingredients on platforms
// TODO: enemy. More frequent attacks on higher score (more ingredients)
// TODO: add sounds for death and collecting ingredients
// TODO: add animated assets for player and enemy and ingredients

const (
	holeSpawnRate float32 = 3
	holeSpeed     float32 = 400
	maxHoleCount  int     = 10

	holeMinWidth = 200
	holeMaxWidth = 400
)

const (
	platformSpawnRate float32 = 2 * holeSpawnRate
	platformSpeed     float32 = holeSpeed
	maxPlatformCount  int     = 10

	platformMinWidth = 200
	platformMaxWidth = 400
)

const (
	playerMaxJumpHeightScreenPercent float32 = 0.2
	playerLeftMarginScreenPercent    float32 = 0.2

	playerInitialVerticalSpeed float32 = holeSpeed - 100
)

const (
	maxCollectibleCount  = 40
	collectibleSpeed     = platformSpeed
	collectibleSpawnRate = holeSpawnRate / 2
)

type gameState struct {
	screen gameScreen

	assets assets

	commonState  commonState
	ground       ground
	platforms    platforms
	collectibles collectibles
	player       player

	musicPlaying bool
	paused       bool
}

func newGameState(assets assets) *gameState {
	return &gameState{
		screen:       gameScreenGame,
		assets:       assets,
		ground:       newGround(assets.ground),
		platforms:    newPlatforms(assets.platform),
		collectibles: newCollectibles(assets.collectibles),
		player:       newPlayer(assets.player),
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
	if rl.IsKeyPressed(rl.KeyR) || (gs.screen == gameScreenGameOver && rl.IsKeyPressed(rl.KeySpace)) {
		*gs = *newGameState(gs.assets)
		return
	}
	if rl.IsKeyPressed(rl.KeyP) {
		gs.paused = !gs.paused
	}

	gs.commonState.update()

	switch gs.screen {
	case gameScreenMenu:
		panic("not implemented")
	case gameScreenGame:
		if !gs.musicPlaying {
			rl.PlayMusicStream(gs.assets.music)
			gs.musicPlaying = true
		} else {
			rl.UpdateMusicStream(gs.assets.music)
		}
		if gs.paused {
			pausedText := "Paused"
			textWidth := rl.MeasureText(pausedText, 60)
			pausedTextX := int32(gs.commonState.screenWidth/2) - textWidth/2
			pausedTextY := int32(gs.commonState.screenHeight / 2)
			rl.DrawText(pausedText, pausedTextX, pausedTextY, 60, rl.Black)
			return
		}
		gs.ground.update(gs.commonState)
		gs.platforms.updatePlatforms(gs.commonState, gs.ground.border, gs.player.border.Height)
		gs.collectibles.updateCollectibles(gs.commonState, gs.ground.border)
		gs.player.update(gs.commonState, gs.ground.border, gs.platforms.borders, gs.ground.holes.borders, gs.collectibles.borders)
		if gs.player.dead {
			gs.screen = gameScreenGameOver
		}
	case gameScreenGameOver:
		rl.StopMusicStream(gs.assets.music)
		gs.musicPlaying = false
	case gameScreenWin:
		panic("not implemented")
	}
}

func (gs *gameState) draw() {
	rl.ClearBackground(rl.SkyBlue)

	rl.DrawFPS(10, 10)

	scoreText := "Score: " + strconv.Itoa(gs.player.score)
	scoreTextWidth := rl.MeasureText(scoreText, 20)
	rl.DrawText(scoreText, int32(gs.commonState.screenWidth)-scoreTextWidth-10, 10, 20, rl.Black)

	switch gs.screen {
	case gameScreenMenu:
		panic("not implemented")
	case gameScreenGame:
		gs.ground.draw()
		gs.platforms.draw()
		gs.collectibles.draw()
		gs.player.draw()
	case gameScreenGameOver:
		gs.drawGameOver()
		gs.ground.draw()
		gs.platforms.draw()
		gs.collectibles.draw()
		gs.player.draw()
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
	assets groundAssets
	holes  holes
}

func newGround(assets groundAssets) ground {
	return ground{
		assets: assets,
		holes:  newHoles(),
	}
}

func (g *ground) update(cs commonState) {
	g.updateBorders(cs)
	g.holes.updateHoles(cs, g.border)
}

func (g *ground) updateBorders(cs commonState) {
	groundHeight := float32(g.assets.center.Height)
	g.border = rl.Rectangle{
		X:      0,
		Y:      cs.screenHeight - groundHeight,
		Width:  cs.screenWidth,
		Height: groundHeight,
	}
}

func (g *ground) draw() {
	for x := int32(0); x < int32(g.border.Width); x += int32(g.assets.center.Width) {
		rl.DrawTexture(g.assets.center, x, int32(g.border.Y), rl.White)
	}

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
		holeBorder.X -= holeSpeed * cs.dt

		holeNotVisible := holeBorder.X+holeBorder.Width < 0

		if holeNotVisible && h.spawnTimer >= holeSpawnRate {
			h.spawnTimer = 0

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
		rl.DrawRectangleRec(h.borders[i], rl.Gray)
	}
}

type platforms struct {
	spawnTimer float32
	borders    []rl.Rectangle

	assets platformAssets
}

func newPlatforms(assets platformAssets) platforms {
	return platforms{
		borders: make([]rl.Rectangle, maxPlatformCount),
		assets:  assets,
	}
}

func (p *platforms) updatePlatforms(cs commonState, groundBorders rl.Rectangle, playerHeight float32) {
	p.spawnTimer += cs.dt

	for i := range p.borders {
		platformBorder := &p.borders[i]
		platformBorder.X -= platformSpeed * cs.dt

		platformNotVisible := platformBorder.X+platformBorder.Width < 0

		if platformNotVisible && p.spawnTimer >= platformSpawnRate {
			p.spawnTimer = 0

			platformWidth := random.Float32(platformMinWidth, platformMaxWidth)

			playerMaxJumpHeight := cs.screenHeight * playerMaxJumpHeightScreenPercent
			platformHeight := float32(p.assets.left.Height)

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
		borderX := int32(p.borders[i].X)
		borderY := int32(p.borders[i].Y)
		// left
		rl.DrawTexture(p.assets.left, borderX, borderY, rl.White)
		// center
		borderLen := borderX + int32(p.borders[i].Width)
		for x := borderX + p.assets.left.Width; x < borderLen-p.assets.right.Width; x += p.assets.center.Width {
			rl.DrawTexture(p.assets.center, x, borderY, rl.White)
		}
		// right
		rl.DrawTexture(p.assets.right, borderLen-p.assets.right.Width, borderY, rl.White)
	}
}

type collectibles struct {
	spawnTimer float32
	borders    []rl.Rectangle

	collectibleAssets collectibleAssets
}

func newCollectibles(assets collectibleAssets) collectibles {
	return collectibles{
		borders:           make([]rl.Rectangle, maxCollectibleCount),
		collectibleAssets: assets,
	}
}

func (c *collectibles) updateCollectibles(cs commonState, groundBorders rl.Rectangle) {
	c.spawnTimer += cs.dt

	for i := range c.borders {
		collectibleBorder := &c.borders[i]
		collectibleBorder.X -= collectibleSpeed * cs.dt

		collectibleNotVisible := collectibleBorder.X+collectibleBorder.Width < 0

		if collectibleNotVisible && c.spawnTimer >= collectibleSpawnRate {
			c.spawnTimer = 0

			collectibleWidth := float32(c.collectibleAssets.ingredient.Width)
			collectibleHeigth := float32(c.collectibleAssets.ingredient.Height)
			collectibleBorder.X = cs.screenWidth + collectibleWidth
			collectibleBorder.Y = groundBorders.Y - collectibleHeigth
			collectibleBorder.Width = collectibleWidth
			collectibleBorder.Height = collectibleHeigth
		}
	}
}

func (c *collectibles) draw() {
	for i := range c.borders {
		rl.DrawTexture(c.collectibleAssets.ingredient, int32(c.borders[i].X), int32(c.borders[i].Y), rl.White)
	}
}

type player struct {
	verticalPosition float32
	verticalSpeed    float32
	// can't rely on verticalPosition == 0 because it will > 0 while player is decending
	jumping    bool
	jumpStartY float32
	onPlatform bool
	border     rl.Rectangle
	dead       bool

	score int

	// assets
	assets playerAssets
}

func newPlayer(assets playerAssets) player {
	return player{
		verticalSpeed: playerInitialVerticalSpeed,
		assets:        assets,
	}
}

func (p *player) update(cs commonState, groundBorders rl.Rectangle, platformBorders, holeBorders, collectibleBorders []rl.Rectangle) {
	p.updateBorder(cs, groundBorders)
	p.updateVerticalPosition(cs, platformBorders)
	p.updateScore(collectibleBorders)
	p.updateDead(holeBorders)
}

func (p *player) updateBorder(cs commonState, groundBorders rl.Rectangle) {
	playerHeight := float32(p.assets.shapes.body.Height)
	playerLeftMargin := cs.screenWidth * playerLeftMarginScreenPercent
	p.border = rl.Rectangle{
		X:      playerLeftMargin,
		Y:      groundBorders.Y - playerHeight - p.verticalPosition,
		Width:  float32(p.assets.shapes.body.Width),
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
				break
			}
		}
	}

	isSpacePressed := rl.IsKeyDown(rl.KeySpace)

	if isSpacePressed &&
		((!p.jumping && p.verticalPosition == 0) || // player on the ground and space pressed -> jump
			p.onPlatform) { // player on the platform and space pressed -> jump
		p.jumping = true
		p.jumpStartY = p.verticalPosition
	}
	// player is jumping and space released -> start descending
	if !isSpacePressed && p.jumping {
		p.jumping = false
	}

	if p.jumping {
		p.verticalPosition += p.verticalSpeed * cs.dt
	} else if !p.onPlatform {
		p.verticalPosition -= p.verticalSpeed * cs.dt
	}

	maxJumpHeight := cs.screenHeight*playerMaxJumpHeightScreenPercent + p.jumpStartY
	p.verticalPosition = rl.Clamp(p.verticalPosition, 0, maxJumpHeight)

	// player reached max jump height -> start descending
	if p.jumping && p.verticalPosition >= maxJumpHeight {
		p.jumping = false
	}
}

func (p *player) updateScore(collectibleBorders []rl.Rectangle) {
	for i := range collectibleBorders {
		if rl.CheckCollisionRecs(p.border, collectibleBorders[i]) {
			p.score++
			rlutils.SetToZero(&collectibleBorders[i])
		}
	}
}

func (p *player) updateDead(holeBorders []rl.Rectangle) {
	if p.verticalPosition == 0 {
		playerMiddle := p.border.X + p.border.Width/2
		for _, holeBorder := range holeBorders {
			if rlutils.VerticalCollision(p.border, holeBorder) &&
				// player is inside the hole more than half
				(holeBorder.X <= playerMiddle && playerMiddle <= holeBorder.X+holeBorder.Width) {
				p.dead = true
				break
			}
		}
	}
}

func (p *player) draw() {
	rl.DrawTexture(p.assets.shapes.body, int32(p.border.X), int32(p.border.Y), rl.White)
	faceTexture := p.assets.shapes.runningFace
	if p.jumping {
		faceTexture = p.assets.shapes.midAirFace
	}
	faceX := (int32(p.border.X) + faceTexture.Width/2)
	rl.DrawTexture(faceTexture, int32(faceX), int32(p.border.Y), rl.White)
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
