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
	msgs := make([]string, 3)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		vector.DrawFilledCircle(screen, float32(x), float32(y), 5, color.RGBA{0, 0, 255, 1}, true)
		msgs = append(msgs, "left mouse button")
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		vector.DrawFilledCircle(screen, float32(x), float32(y), 5, color.RGBA{0, 255, 0, 1}, true)
		msgs = append(msgs, "right mouse button")
	}

	if g.touch.IsTouchingThree() {
		msgs = append(msgs, "touching three")
	} else if g.touch.IsTouchingTwo() {
		msgs = append(msgs, "touching two")
		if pan := g.touch.Pan(); pan != nil {
			msgs = append(msgs, "pan")
			deltaX := pan.OriginX - pan.PrevX
			if deltaX < -10 {
				msgs = append(msgs, fmt.Sprintf("swipe right - delta: %d", deltaX))
			} else if deltaX > 10 {
				msgs = append(msgs, fmt.Sprintf("swipe left - delta: %d", deltaX))
			}

			vector.DrawFilledCircle(screen, float32(pan.OriginX), float32(g.h)/2, 5, color.RGBA{255, 0, 0, 1}, true)
			vector.DrawFilledCircle(screen, float32(pan.PrevX), float32(g.h)/2, 5, color.RGBA{0, 255, 0, 1}, true)
			vector.StrokeLine(screen, float32(pan.OriginX), float32(g.h)/2, float32(pan.PrevX), float32(g.h)/2, 1, color.White, true)
		}

		if pinch := g.touch.Pinch(); pinch != nil {
			msgs = append(msgs, "pinch")
		}
	} else if g.touch.IsTouching() {
		x, y, _ := g.touch.GetFirstTouchPosition()
		msgs = append(msgs, "touching one")
		vector.DrawFilledCircle(screen, float32(x), float32(y), 5, color.RGBA{0, 0, 255, 1}, true)
	}

	for ix, m := range msgs {
		textFace := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   24,
		}
		w, h := text.Measure(m, textFace, 0)
		textOp := &text.DrawOptions{}
		textOp.GeoM.Translate(w/2+10, h*float64(ix+1))
		textOp.PrimaryAlign = text.AlignCenter
		textOp.SecondaryAlign = text.AlignCenter
		text.Draw(screen, m, textFace, textOp)
	}

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
