package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/manuelpepe/gol/utils"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

type Gesture struct {
	w, h int

	touch *utils.TouchTracker
}

func NewGestureDemo(width, height int) *Gesture {

	return &Gesture{
		w: width,
		h: height,

		touch: utils.NewTouchTracker(),
	}
}

// Max TPS as specified by ebiten
const MAX_TPS = 60

// Delay between updates
const DELAY_SEC = 1

func (g *Gesture) Update() error {
	g.touch.Update()
	return nil
}

func (g *Gesture) Draw(screen *ebiten.Image) {
	status := ""
	if g.touch.IsTouchingThree() {
		status = "touching three"
	} else if g.touch.IsTouchingTwo() {
		status = "touching two"
		if pan := g.touch.Pan(); pan != nil {
			status += " - pan"
			deltaX := pan.OriginX - pan.PrevX
			if deltaX < -10 {
				status += fmt.Sprintf(" - swipe right - delta: %.2f", deltaX)
			} else if deltaX > 10 {
				status += fmt.Sprintf(" - swipe right - delta: %.2f", deltaX)
			}

			vector.DrawFilledCircle(screen, float32(pan.OriginX), float32(g.h)/2, 5, color.RGBA{255, 0, 0, 1}, true)
			vector.DrawFilledCircle(screen, float32(pan.PrevX), float32(g.h)/2, 5, color.RGBA{0, 255, 0, 1}, true)
			vector.StrokeLine(screen, float32(pan.OriginX), float32(g.h)/2, float32(pan.PrevX), float32(g.h)/2, 1, color.White, true)
		}

		if pinch := g.touch.Pinch(); pinch != nil {
			status += "- pinch"
		}
	} else if g.touch.IsTouching() {
		x, y, _ := g.touch.GetFirstTouchPosition()
		status = "touching one"
		vector.DrawFilledCircle(screen, float32(x), float32(y), 5, color.RGBA{0, 0, 255, 1}, true)
	}

	str := fmt.Sprintf("STATUS: %s", status)
	textFace := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}
	_, h := text.Measure(str, textFace, 0)
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(g.w)/2, h)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, str, textFace, textOp)
}

func (g *Gesture) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	W, H := 300, 500
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGestureDemo(W, H)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
