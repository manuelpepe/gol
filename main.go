package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/manuelpepe/gol/gol"
)

type GameOfLife struct {
	cols, rows int
	ticks      int
	grid       []bool
}

func NewGameOfLife(x, y int) *GameOfLife {
	return &GameOfLife{
		cols:  x,
		rows:  y,
		ticks: 0,
		grid:  make([]bool, x*y),
	}
}

// Max TPS as specified by ebiten
const MAX_TPS = 60

func (g *GameOfLife) Update() error {
	interval := MAX_TPS * 2
	g.ticks = (g.ticks + 1) % interval
	if g.ticks == interval-1 {
		g.grid = gol.NextGrid(g.cols, g.rows, g.grid)
	}
	return nil
}

func (g *GameOfLife) Draw(screen *ebiten.Image) {
	str := ""
	for y := range g.rows {
		for x := range g.cols {
			if g.grid[y*g.rows+x] {
				str += "[x] "
			} else {
				str += "[ ] "
			}
		}
		str += "\n"
	}
	ebitenutil.DebugPrint(screen, str)
}

func (g *GameOfLife) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGameOfLife(10, 10)

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
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
