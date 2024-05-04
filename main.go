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
	cellWidth  int
	cellHeight int

	usableWidth  int
	usableHeight int

	padding int
}

const PADDING = 1

func newmetadata(width, height, x, y int) metadata {
	usableWidth := width - x*PADDING
	usableHeight := height - y*PADDING
	return metadata{
		usableWidth:  usableWidth,
		usableHeight: usableHeight,

		cellWidth:  usableWidth / x,
		cellHeight: usableHeight / y,
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

const DELAY = 1

func (g *GameOfLife) Update() error {
	interval := MAX_TPS * DELAY
	g.ticks = (g.ticks + 1) % interval
	if g.ticks == interval-1 {
		g.grid = gol.NextGrid(g.cols, g.rows, g.grid)
	}
	return nil
}

func (g *GameOfLife) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for y := range g.rows {
		for x := range g.cols {
			p := y*g.rows + x
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Translate(
				float64((g.md.cellWidth+PADDING)*x),
				float64((g.md.cellHeight+PADDING)*y),
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
	W, H := 640, 480
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGameOfLife(W, H, 10, 10)

	// block
	// game.grid[10] = true
	// game.grid[11] = true
	// game.grid[20] = true
	// game.grid[21] = true

	// toad
	game.grid[22] = true
	game.grid[23] = true
	game.grid[24] = true
	game.grid[31] = true
	game.grid[32] = true
	game.grid[33] = true

	// blinker
	game.grid[46] = true
	game.grid[47] = true
	game.grid[48] = true
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
