package player

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pechorka/bread-game-jam/pkg/rlutils"
)

const (
	initialVerticalSpeed       float32 = 3
	maxJumpHeightScreenPercent float32 = 0.2
)

type Player struct {
	verticalPosition float32
	verticalSpeed    float32
	jumping          bool
	Rect             rl.Rectangle
}

func New() *Player {
	return &Player{
		verticalSpeed: initialVerticalSpeed,
	}
}

func (p *Player) Update(ground rl.Rectangle) {
	w, h := rlutils.GetScreenDimensions()

	p.updateVerticalPosition(h)
	p.updateRect(w, h, ground)
}

func (p *Player) updateVerticalPosition(h float32) {
	isSpacePressed := rl.IsKeyDown(rl.KeySpace)
	if isSpacePressed && !p.jumping && p.verticalPosition == 0 {
		p.jumping = true
	}
	if !isSpacePressed && p.jumping {
		p.jumping = false
	}

	if p.jumping {
		p.verticalPosition += p.verticalSpeed
	} else {
		p.verticalPosition -= p.verticalSpeed
	}
	maxJumpHeight := h * maxJumpHeightScreenPercent
	p.verticalPosition = rl.Clamp(p.verticalPosition, 0, maxJumpHeight)
}

func (p *Player) updateRect(w, h float32, ground rl.Rectangle) {
	playerHeight := h * 0.1
	p.Rect = rl.Rectangle{
		X:      w * 0.2,
		Y:      h - ground.Height - playerHeight - p.verticalPosition,
		Width:  w * 0.1,
		Height: playerHeight,
	}
}

func (p *Player) Draw() {
	rl.DrawRectangleRec(p.Rect, rl.Yellow)
}
