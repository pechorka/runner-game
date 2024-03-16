package ground

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/pechorka/bread-game-jam/pkg/rlutils"
)

const (
	screenHeightPercentage float32 = 0.2
)

type Ground struct {
	Rect rl.Rectangle
}

func New() *Ground {
	return &Ground{}
}

func (g *Ground) Update() {
	w, h := rlutils.GetScreenDimensions()

	groundHeight := h * screenHeightPercentage
	g.Rect = rl.Rectangle{
		X:      0,
		Y:      h - groundHeight,
		Width:  w,
		Height: groundHeight,
	}
}

func (g *Ground) Draw() {
	rl.DrawRectangleRec(g.Rect, rl.Brown)
}
