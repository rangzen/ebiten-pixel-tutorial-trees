package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1024
	screenHeight = 768
	spriteSize   = 32
)

//go:embed trees.png
var treesPng []byte
var treesImg *ebiten.Image
var treeFrames []image.Rectangle

func init() {
	var err error
	treesDecoded, _, err := image.Decode(bytes.NewReader(treesPng))
	if err != nil {
		log.Fatal(err)
	}
	treesImg = ebiten.NewImageFromImage(treesDecoded)

	for x := treesImg.Bounds().Min.X; x < treesImg.Bounds().Max.X; x += 32 {
		for y := treesImg.Bounds().Min.Y; y < treesImg.Bounds().Max.Y; y += 32 {
			treeFrames = append(treeFrames, image.Rect(x, y, x+32, y+32))
		}
	}
}

type Vec struct {
	X, Y float64
}

type Game struct {
	camPos       Vec
	camSpeed     float64
	camZoom      float64
	camZoomSpeed float64
	trees        []int
	matrices     []ebiten.GeoM
	last         time.Time
}

func (g *Game) Update() error {
	dt := time.Since(g.last).Seconds()
	g.last = time.Now()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) { // && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.trees = append(g.trees, rand.Intn(len(treeFrames)))
		unproject := g.cam()
		unproject.Invert()
		mx, my := ebiten.CursorPosition()
		ux, uy := unproject.Apply(float64(mx), float64(my))
		mat := ebiten.GeoM{}
		mat.Translate(ux, uy)
		mat.Translate(-spriteSize/2, -spriteSize/2)
		g.matrices = append(g.matrices, mat)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.camPos.X -= g.camSpeed * dt * 1 / g.camZoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.camPos.X += g.camSpeed * dt * 1 / g.camZoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.camPos.Y += g.camSpeed * dt * 1 / g.camZoom
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.camPos.Y -= g.camSpeed * dt * 1 / g.camZoom
	}
	_, wy := ebiten.Wheel()
	g.camZoom *= math.Pow(g.camZoomSpeed, wy)
	return nil
}

func (g Game) cam() ebiten.GeoM {
	cam := ebiten.GeoM{}
	cam.Translate(-g.camPos.X, -g.camPos.Y)
	cam.Scale(g.camZoom, g.camZoom)
	cam.Translate(g.camPos.X, g.camPos.Y)
	cam.Translate(screenWidth/2-g.camPos.X, screenHeight/2-g.camPos.Y)
	return cam
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Forestgreen)
	for i, index := range g.trees {
		op := &ebiten.DrawImageOptions{}
		op.GeoM = g.matrices[i]
		op.GeoM.Concat(g.cam())
		screen.DrawImage(treesImg.SubImage(treeFrames[index]).(*ebiten.Image), op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f for %d trees", ebiten.CurrentTPS(), len(g.trees)))
}

func (g *Game) Layout(_, _ int) (sw, sh int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Transpose Pixel tutorial to Ebiten")
	game := Game{
		camPos: Vec{
			X: 0,
			Y: 0,
		},
		camSpeed:     500,
		camZoom:      1,
		camZoomSpeed: 1.2,
		trees:        []int{},
		matrices:     []ebiten.GeoM{},
		last:         time.Now(),
	}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
