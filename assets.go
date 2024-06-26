package main

import (
	"embed"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets
var assetsFs embed.FS

type assets struct {
	player       playerAssets
	ground       groundAssets
	platform     platformAssets
	collectibles collectibleAssets

	music rl.Music

	allTextures []rl.Texture2D
}

type playerAssets struct {
	shapes struct {
		body        rl.Texture2D
		runningFace rl.Texture2D
		midAirFace  rl.Texture2D
	}
}

type groundAssets struct {
	center rl.Texture2D
}

type platformAssets struct {
	left   rl.Texture2D
	center rl.Texture2D
	right  rl.Texture2D
}

type collectibleAssets struct {
	ingredient rl.Texture2D
}

func loadAssets() assets {
	a := assets{}

	a.player.shapes.body = a.loadTexture("assets/player/shapes/yellow_body_squircle.png")
	a.player.shapes.runningFace = a.loadTexture("assets/player/shapes/face_a.png")
	a.player.shapes.midAirFace = a.loadTexture("assets/player/shapes/face_g.png")

	a.ground.center = a.loadTexture("assets/ground/tile_center.png")

	a.platform.left = a.loadTexture("assets/platform/tile_half_left.png")
	a.platform.center = a.loadTexture("assets/platform/tile_half_center.png")
	a.platform.right = a.loadTexture("assets/platform/tile_half_right.png")

	a.collectibles.ingredient = a.loadTexture("assets/collectibles/tile_coin.png")

	a.music = loadMusicStream("assets/music/I Got a Stick Arr Bryan Teoh.mp3")

	return a
}

func (a *assets) loadTexture(path string) rl.Texture2D {
	t := loadTextureFromImage(path)
	a.allTextures = append(a.allTextures, t)
	return t
}

func (a *assets) unload() {
	for _, t := range a.allTextures {
		rl.UnloadTexture(t)
	}
	rl.UnloadMusicStream(a.music)
}

func loadTextureFromImage(imgPath string) rl.Texture2D {
	file, err := assetsFs.ReadFile(imgPath)
	if err != nil {
		panic(err)
	}
	fileExtention := imgPath[strings.LastIndexByte(imgPath, '.'):]
	img := rl.LoadImageFromMemory(fileExtention, file, int32(len(file)))
	defer rl.UnloadImage(img)
	texture := rl.LoadTextureFromImage(img)

	return texture
}

func loadMusicStream(musicPath string) rl.Music {
	file, err := assetsFs.ReadFile(musicPath)
	if err != nil {
		panic(err)
	}
	fileExtention := musicPath[strings.LastIndexByte(musicPath, '.'):]
	music := rl.LoadMusicStreamFromMemory(fileExtention, file, int32(len(file)))

	return music
}
