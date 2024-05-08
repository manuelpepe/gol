package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	touchutils "github.com/manuelpepe/ebiten-touchutils"
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

	scale                  float64
	camX, camY             float64
	panOriginX, panOriginY float64

	ticks   int
	running bool

	delayOptions []float64
	delaySetting int

	grid      []bool
	useShader bool

	aliveImage *ebiten.Image
	deadImage  *ebiten.Image

	md metadata

	touch *touchutils.TouchTracker
}

type metadata struct {
	width  int
	height int

	cellWidth  int
	cellHeight int

	padding int
}

const PADDING int = 1

func newmetadata(width, height, x, y int) metadata {
	return metadata{
		width:  width,
		height: height,

		cellWidth:  10,
		cellHeight: 10,

		padding: PADDING,
	}
}

func (md metadata) FullCellSize() (int, int) {
	return md.FullCellWidth(), md.FullCellHeight()
}

func (md metadata) FullCellWidth() int {
	return md.cellWidth + md.padding
}

func (md metadata) FullCellHeight() int {
	return md.cellHeight + md.padding
}

func (md metadata) ScaledCellWidth(scale float64) float64 {
	return float64(md.cellWidth)*scale + float64(md.padding)
}

func (md metadata) ScaledCellHeight(scale float64) float64 {
	return float64(md.cellHeight)*scale + float64(md.padding)
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

		scale: 1,
		camX:  0,
		camY:  0,

		running: false,
		ticks:   0,

		delayOptions: []float64{0.1, 0.05, 0.02, 0.2},
		delaySetting: 0,

		grid:      make([]bool, x*y),
		useShader: true,

		aliveImage: aliveImage,
		deadImage:  deadImage,

		md: md,

		touch: touchutils.NewTouchTracker(),
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
		if g.useShader {
			g.grid = gol.NextGridShader(g.cols, g.rows, g.grid)
		} else {
			g.grid = gol.NextGrid(g.cols, g.rows, g.grid)
		}
	}
	g.handleInputs()
	return nil
}

func (g *GameOfLife) handleInputs() {
	// handle map panning
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		x, y := ebiten.CursorPosition()
		g.panOriginX = float64(x)
		g.panOriginY = float64(y)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		g.panOriginX = 0
		g.panOriginY = 0
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		x, y := ebiten.CursorPosition()
		offsetX := (g.panOriginX - float64(x)) * 0.01
		offsetY := (g.panOriginY - float64(y)) * 0.01
		g.camX += offsetX
		g.camY += offsetY
		g.camX = min(float64(g.md.FullCellWidth()*g.cols-g.md.width)/float64(g.md.FullCellWidth()), max(0, g.camX))
		g.camY = min(float64(g.md.FullCellHeight()*g.rows-g.md.height)/float64(g.md.FullCellHeight()), max(0, g.camY))
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

	if inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) {
		g.scale = min(g.scale+0.1, 1.5)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyNumpadSubtract) || inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.scale = max(g.scale-0.1, 0.5)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.useShader = !g.useShader
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
	g.handleTouches()
}

func (g *GameOfLife) handleTouches() {
	g.touch.Update()

	if _, _, _, ok := g.touch.TappedThree(); ok {
		g.grid = make([]bool, g.cols*g.rows)
	} else if _, _, ok := g.touch.TappedTwo(); ok {
		g.running = !g.running
	} else if pan, ok := g.touch.TwoFingerPan(); ok {
		deltaX := pan.OriginX - pan.LastX
		if deltaX < -10 {
			// swipe right
			g.nextSpeedModifier()
		} else if deltaX > 10 {
			// swipe left
			g.running = false
		}
	} else if g.touch.IsTouching() {
		x, y, ok := g.touch.GetFirstTouchPosition()
		pos, ok2 := g.positionToCell(x, y)
		if ok && ok2 {
			g.grid[pos] = true
		}
	}
}

func (g *GameOfLife) getCursorPositionInGrid() (int, bool) {
	x, y := ebiten.CursorPosition()
	return g.positionToCell(x, y)
}

func (g *GameOfLife) positionToCell(x, y int) (int, bool) {
	minX := int(math.Round(g.camX * float64(g.md.ScaledCellWidth(g.scale))))
	minY := int(math.Round(g.camY * float64(g.md.ScaledCellHeight(g.scale))))
	realX := float64(x+minX) / g.md.ScaledCellWidth(g.scale)
	realY := float64(y+minY) / g.md.ScaledCellHeight(g.scale)
	pos := int(realY)*g.cols + int(realX)
	return pos, pos >= 0 && pos < g.cols*g.rows
}

func (g *GameOfLife) nextSpeedModifier() {
	g.delaySetting = (g.delaySetting + 1) % len(g.delayOptions)
}

func (g *GameOfLife) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// screen boundries
	minX := int(math.Floor(g.camX))
	minY := int(math.Floor(g.camY))
	maxX := int(float64(g.md.width)/g.md.ScaledCellWidth(g.scale) + float64(minX))
	maxY := int(float64(g.md.height)/g.md.ScaledCellHeight(g.scale) + float64(minY))

	for y := range g.rows {
		if y < minY {
			continue // don't draw off screen
		}

		if y > maxY {
			break // stop drawing off screen
		}

		for x := range g.cols {
			if x < minX || x > maxX {
				continue // don't draw off screen
			}

			p := y*g.rows + x
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Scale(g.scale, g.scale)
			opts.GeoM.Translate(
				float64(g.md.ScaledCellWidth(g.scale))*(float64(x)-g.camX),
				float64(g.md.ScaledCellWidth(g.scale))*(float64(y)-g.camY),
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

	v, _ := g.getCursorPositionInGrid()
	x, y := ebiten.CursorPosition()

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"TPS: %.2f\nFPS: %.2f\n\n\n\n\np: %d - x: %d y: %d\ncamx: %.2f camy: %.2f\nscale: %.2f\nscaled size: %.2f\nshader: %v",
			ebiten.ActualTPS(), ebiten.ActualFPS(),
			v, x, y,
			g.camX, g.camY,
			g.scale,
			g.md.ScaledCellWidth(g.scale),
			g.useShader),
	)
}

func (g *GameOfLife) Layout(_, _ int) (_, _ int) {
	panic("unused in favor of LayoutF")
}

func (_ *GameOfLife) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	canvasWidth := math.Ceil(logicWinWidth * scale)
	canvasHeight := math.Ceil(logicWinHeight * scale)
	return canvasWidth, canvasHeight
}

func main() {
	W, H := 1024, 720
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Hello, World!")
	game := NewGameOfLife(W, H, 5000, 5000)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
