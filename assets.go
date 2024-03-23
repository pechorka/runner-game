package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type assets struct {
	player playerAssets

	allTextures []rl.Texture2D
}

type playerAssets struct {
	shapes struct {
		body        rl.Texture2D
		runningFace rl.Texture2D
		midAirFace  rl.Texture2D
	}
}

func loadAssets() assets {
	a := assets{}

	a.player.shapes.body = a.loadTexture("assets/player/shapes/yellow_body_squircle.png")
	a.player.shapes.runningFace = a.loadTexture("assets/player/shapes/face_a.png")
	a.player.shapes.midAirFace = a.loadTexture("assets/player/shapes/face_g.png")

	return a
}

func (a *assets) loadTexture(path string) rl.Texture2D {
	t := rl.LoadTexture(path)
	a.allTextures = append(a.allTextures, t)
	return t
}

func (a *assets) unload() {
	for _, t := range a.allTextures {
		rl.UnloadTexture(t)
	}
}
