package ground

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pechorka/bread-game-jam/pkg/rlutils"
)

const (
	initialHoleCount = 10
	holeMinWidth     = 100
	holeMaxWidth     = 200
)

const (
	screenHeightPercentage float32 = 0.2
	spawnHoleTime          float32 = 1
	holeSpeed              float32 = 10
)

type Ground struct {
	Rect  rl.Rectangle
	Holes []*hole

	holesTimer float32
}

type hole struct {
	rect   rl.Rectangle
	active bool
}

func New() *Ground {
	holes := make([]*hole, 0, initialHoleCount)
	for i := 0; i < initialHoleCount; i++ {
		holes = append(holes, &hole{})
	}
	return &Ground{
		Holes: holes,
	}
}

func (g *Ground) Update() {
	g.updateTimers()

	w, h := rlutils.GetScreenDimensions()

	groundHeight := h * screenHeightPercentage
	g.Rect = rl.Rectangle{
		X:      0,
		Y:      h - groundHeight,
		Width:  w,
		Height: groundHeight,
	}

	g.updateHoles()
}

func (g *Ground) AboveHole(p rl.Rectangle) bool {
	for _, hole := range g.Holes {
		if hole.active && rlutils.VerticalCollision(p, hole.rect) {
			return true
		}
	}
	return false
}

func (g *Ground) updateHoles() {
	for _, hole := range g.Holes {
		if hole.active {
			hole.rect.X -= holeSpeed
		}

		if hole.rect.X+hole.rect.Width < g.Rect.X {
			hole.active = false
		}

		if !hole.active && g.holesTimer > spawnHoleTime {
			g.holesTimer = 0
			hole.active = true
			hole.rect = spawnHole(g.Rect)
		}
	}
}

func spawnHole(groundRect rl.Rectangle) rl.Rectangle {
	holeWidth := float32(rl.GetRandomValue(holeMinWidth, holeMaxWidth))
	return rl.Rectangle{
		X:      groundRect.Width + holeWidth,
		Y:      groundRect.Y,
		Width:  holeWidth,
		Height: groundRect.Height,
	}
}

func (g *Ground) updateTimers() {
	dt := rl.GetFrameTime()
	g.holesTimer += dt
}

func (g *Ground) Draw() {
	rl.DrawRectangleRec(g.Rect, rl.Brown)
	for _, hole := range g.Holes {
		rl.DrawRectangleRec(hole.rect, rl.Black)
	}
}
