package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/manuelpepe/gol/gol"
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

type GameOfLife struct {
	cols, rows int

	ticks   int
	running bool

	delayOptions []float64
	delaySetting int

	grid []bool

	aliveImage *ebiten.Image
	deadImage  *ebiten.Image

	md metadata
}

type metadata struct {
	width  int
	height int

	cellWidthPrecise  float64
	cellHeightPrecise float64

	cellWidth  int
	cellHeight int

	usableWidth  int
	usableHeight int

	padding float64
}

const PADDING float64 = 1

func newmetadata(width, height, x, y int) metadata {
	usableWidth := int(float64(width) - float64(x)*PADDING)
	usableHeight := int(float64(height) - float64(y)*PADDING)
	return metadata{
		width:  width,
		height: height,

		usableWidth:  usableWidth,
		usableHeight: usableHeight,

		cellWidth:  usableWidth / x,
		cellHeight: usableHeight / y,

		cellWidthPrecise:  float64(usableWidth) / float64(x),
		cellHeightPrecise: float64(usableHeight) / float64(y),

		padding: PADDING,
	}
}

func NewGameOfLife(width, height, x, y int) *GameOfLife {
	md := newmetadata(width, height, x, y)

	aliveImage := ebiten.NewImage(md.cellWidth, md.cellHeight)
	aliveImage.Fill(color.White)
	deadImage := ebiten.NewImage(md.cellWidth, md.cellHeight)
	deadImage.Fill(color.Gray16{0x9999})

	return &GameOfLife{
		cols: x,
		rows: y,

		running: false,
		ticks:   0,

		delayOptions: []float64{2, 1, 0.5, 0.25, 0.1, 0.02},
		delaySetting: 0,

		grid: make([]bool, x*y),

		aliveImage: aliveImage,
		deadImage:  deadImage,

		md: md,
	}
}

// Max TPS as specified by ebiten
const MAX_TPS = 60

// Delay between updates
const DELAY_SEC = 1

func (g *GameOfLife) Update() error {
	interval := int(MAX_TPS * g.delayOptions[g.delaySetting])
	g.ticks = (g.ticks + 1) % interval
	if g.running && g.ticks == interval-1 {
		g.grid = gol.NextGrid(g.cols, g.rows, g.grid)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ix, ok := g.getCursorPositionInGrid()
		if ok {
			g.grid[ix] = true
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		ix, ok := g.getCursorPositionInGrid()
		if ok {
			g.grid[ix] = false
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.running = !g.running
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.nextSpeedModifier()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.grid = make([]bool, g.cols*g.rows)
	}
	return nil
}

func (g *GameOfLife) getCursorPositionInGrid() (int, bool) {
	x, y := ebiten.CursorPosition()
	cx := int(float64(x) / (g.md.cellWidthPrecise + g.md.padding))
	cy := int(float64(y) / (g.md.cellHeightPrecise + g.md.padding))
	pos := cy*g.cols + cx
	if pos >= 0 && pos < g.cols*g.rows {
		return pos, true
	}
	return pos, false
}

func (g *GameOfLife) nextSpeedModifier() {
	g.delaySetting = (g.delaySetting + 1) % len(g.delayOptions)
}

func (g *GameOfLife) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for y := range g.rows {
		for x := range g.cols {
			p := y*g.rows + x
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Translate(
				(g.md.cellWidthPrecise+g.md.padding)*float64(x),
				(g.md.cellHeightPrecise+g.md.padding)*float64(y),
			)

			if g.grid[p] {
				screen.DrawImage(g.aliveImage, &opts)
			} else {
				screen.DrawImage(g.deadImage, &opts)
			}
		}
	}

	str := fmt.Sprintf("Speed: %.2fs", g.delayOptions[g.delaySetting])
	textFace := &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}
	w, h := text.Measure(str, textFace, 0)
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(w, h)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(screen, str, textFace, textOp)
}

func (g *GameOfLife) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	W, H := 1024, 720
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGameOfLife(W, H, 100, 100)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
