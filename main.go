package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/manuelpepe/gol/gol"
)

type GameOfLife struct {
	cols, rows int

	ticks int

	grid []bool

	aliveImage *ebiten.Image
	deadImage  *ebiten.Image

	md metadata
}

type metadata struct {
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
		cols:       x,
		rows:       y,
		ticks:      0,
		grid:       make([]bool, x*y),
		aliveImage: aliveImage,
		deadImage:  deadImage,
		md:         md,
	}
}

// Max TPS as specified by ebiten
const MAX_TPS = 60

// Delay between updates
const DELAY_SEC = 1

func (g *GameOfLife) Update() error {
	interval := MAX_TPS * DELAY_SEC
	g.ticks = (g.ticks + 1) % interval
	if g.ticks == interval-1 {
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
