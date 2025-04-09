package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 480
	screenHeight = 320

	initialPlatformX = 50
	platformWidth    = 40
	platformHeight   = 100
)

type Game struct {
	state        GameState
	stickLength  float64
	stickRotated bool
	stickAngle   float64
	platforms    []Platform

	charX      float64
	charMoving bool

	bgOffset     float64
	scrollTarget float64
	scrolling    bool
}

type GameState int

const (
	StateIdle GameState = iota
	StateStretching
	StateRotating
	StateWalking
	StateFalling
)

type Platform struct {
	X int
	W int
}

func NewGame() *Game {
	startPlatform := Platform{X: initialPlatformX, W: platformWidth}
	nextPlatform := Platform{
		X: initialPlatformX + 100 + rand.Intn(100),
		W: 40 + rand.Intn(30),
	}

	return &Game{
		state:       StateIdle,
		platforms:   []Platform{startPlatform, nextPlatform},
		stickLength: 0,
		charX:       float64(startPlatform.X + 10),
	}
}

func (g *Game) Update() error {
	switch g.state {
	case StateIdle:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.state = StateStretching
		}
	case StateStretching:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.stickLength += 2
		} else {
			g.state = StateRotating
		}
	case StateRotating:
		if g.stickAngle < 90 {
			g.stickAngle += 4 // tốc độ xoay
		} else {
			g.stickAngle = 90
			g.state = StateWalking // chuyển sang bước tiếp theo sau khi gậy đã xoay xong
		}
	case StateWalking:
		g.charMoving = true
		if g.charX < float64(g.platforms[0].X+g.platforms[0].W)+g.stickLength {
			g.charX += 2 // tốc độ đi
		} else {
			g.charMoving = false

			stickEndX := float64(g.platforms[0].X+g.platforms[0].W) + g.stickLength
			targetPlatform := g.platforms[1]
			if stickEndX >= float64(targetPlatform.X) && stickEndX <= float64(targetPlatform.X+targetPlatform.W) {
				g.charMoving = false
				g.scrolling = true
				g.state = StateIdle
				g.scrollTarget = float64(targetPlatform.X - g.platforms[0].X)
				// g.resetLevel()
			} else {
				g.state = StateFalling
			}
		}
	case StateFalling:
		// g.resetLevel()
		g.state = StateIdle
	}

	if g.scrolling {
		g.bgOffset += 4
		if g.bgOffset >= g.scrollTarget {
			delta := g.scrollTarget
			// Dịch tất cả platform
			for i := range g.platforms {
				g.platforms[i].X -= int(delta)
			}
			// Dịch nhân vật
			g.charX -= delta

			g.bgOffset = 0
			g.scrolling = false
			g.resetLevel()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{240, 240, 240, 255})
	// Vẽ các platform
	for _, p := range g.platforms {
		ebitenutil.DrawRect(screen, float64(p.X)-g.bgOffset, float64(screenHeight-platformHeight), float64(p.W), float64(platformHeight), color.Black)
	}

	// Vẽ gậy
	startX := float64(g.platforms[0].X+g.platforms[0].W) - g.bgOffset
	startY := float64(screenHeight - platformHeight)

	switch g.state {
	case StateStretching:
		ebitenutil.DrawRect(screen, startX, startY-g.stickLength, 2, g.stickLength, color.RGBA{150, 75, 0, 255})
	case StateRotating, StateWalking:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-1, -g.stickLength)
		op.GeoM.Rotate(g.stickAngle * (3.14159265 / 180))
		op.GeoM.Translate(startX, startY)

		stickImage := ebiten.NewImage(2, int(g.stickLength))
		stickImage.Fill(color.RGBA{150, 75, 0, 255})
		screen.DrawImage(stickImage, op)
	}

	// Vẽ nhân vật (placeholder)
	ebitenutil.DrawRect(screen, g.charX-g.bgOffset, float64(screenHeight-platformHeight-20), 10, 20, color.RGBA{255, 0, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) resetLevel() {
	g.platforms[0] = g.platforms[1]
	next := Platform{
		X: g.platforms[1].X + 100 + rand.Intn(100),
		W: 40 + rand.Intn(30),
	}
	g.platforms[1] = next

	// Reset gậy và nhân vật
	g.stickLength = 0
	g.stickAngle = 0
}

func main() {
	ebiten.SetWindowTitle("Stick Hero")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
