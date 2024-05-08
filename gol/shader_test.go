package gol

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type game struct {
	m    *testing.M
	code int
}

func (g *game) Update() error {
	g.code = g.m.Run()
	return ebiten.Termination
}

func (*game) Draw(*ebiten.Image) {
}

func (*game) Layout(int, int) (int, int) {
	return 320, 240
}

func MainWithRunLoop(m *testing.M) {
	// Run an Ebiten process so that (*Image).At is available.
	g := &game{
		m:    m,
		code: 1,
	}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
	os.Exit(g.code)
}

func TestMain(m *testing.M) {
	MainWithRunLoop(m)
}

type XY struct {
	X, Y int
}

func logPixels(x, y int, img *ebiten.Image) {
	b := make([]byte, 4*x*y)
	img.ReadPixels(b)
	log.Printf("%+v\n", b)
}

//go:embed shader.kage
var shaderProgramT []byte

func TestShaderEncodeDecode(t *testing.T) {
	const w, h = 10, 10
	s, err := ebiten.NewShader([]byte(`
		//kage:unit pixels

		package gol

		func Fragment(pos vec4, srcPos vec2, color vec4) vec4 {
			return imageSrc0At(srcPos.xy)
		}
	`))
	if err != nil {
		t.Fatal(err)
	}

	gridData := ebiten.NewImage(w, h)

	live := []XY{
		{2, 2},
		{2, 3},
		{3, 2},
		{3, 3},
	}

	for x := range w {
		for y := range h {
			if slices.Contains(live, XY{x, y}) {
				gridData.Set(x, y, color.RGBA{R: 0xff})
			}
		}
	}

	dst := ebiten.NewImage(w, h)
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = gridData
	dst.DrawRectShader(w, h, s, opts)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			got := dst.At(x, y).(color.RGBA)
			var want color.RGBA
			if slices.Contains(live, XY{x, y}) {
				want = color.RGBA{R: 0xff}
			} else {
				want = color.RGBA{}
			}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", x, y, got, want)
			}
		}
	}

}

func TestShaderCountNeighs(t *testing.T) {
	s, err := ebiten.NewShader([]byte(`
		//kage:unit pixels

		package gol

		func isAlive(pos vec2) float {
			c := imageSrc0At(pos)
			if c.a > 0 {
				return 1
			}
			return 0
		}

		func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
			alive := 0.0
			pos := srcPos.xy

			deltas := [...]float{-1, 0, 1}

			for x := 0; x < len(deltas); x+=1 {
				for y := 0; y < len(deltas); y+=1 {
					if deltas[x] == 0 && deltas[y] == 0 {
						continue
					}

					alive += isAlive(pos + vec2(deltas[x], deltas[y]))
				}
			}

			return vec4(alive/255, 0, 0, 0)
		}
		
	`))
	if err != nil {
		t.Fatal(err)
	}

	const w, h = 5, 5

	gridData := ebiten.NewImage(w, h)
	gridData.Set(2, 2, color.RGBA{A: 0xff})

	dst := ebiten.NewImage(w, h)
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = gridData
	dst.DrawRectShader(w, h, s, opts)

	neighs := []XY{
		{1, 1},
		{1, 2},
		{1, 3},
		{2, 1},
		{2, 3},
		{3, 1},
		{3, 2},
		{3, 3},
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			got := dst.At(x, y).(color.RGBA)
			var want color.RGBA
			if slices.Contains(neighs, XY{x, y}) {
				want = color.RGBA{R: 1}
			} else {
				want = color.RGBA{}
			}
			if got != want {
				t.Errorf("dst.At(%d, %d): got: %v, want: %v", x, y, got, want)
			}
		}
	}

}

var dumpList []bool

type t struct {
	w, h int
}

var ts = []t{
	{12, 12},
	{100, 100},
	{1000, 1000},
	{3000, 3000},
}

func BenchmarkNextGridShaderV2(b *testing.B) {
	for _, tt := range ts {
		l := make([]bool, tt.w*tt.h)
		b.Run(fmt.Sprintf("window-(H:%d-W:%d)", tt.h, tt.w), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l = NextGridShaderV2(tt.w, tt.h, l)
			}
		})
		dumpList = l
	}
}

func BenchmarkNextGridShader(b *testing.B) {
	for _, tt := range ts {
		l := make([]bool, tt.w*tt.h)
		b.Run(fmt.Sprintf("window-(H:%d-W:%d)", tt.h, tt.w), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l = NextGridShader(tt.w, tt.h, l)
			}
		})
		dumpList = l
	}
}

func BenchmarkNextGrid(b *testing.B) {
	for _, tt := range ts {
		l := make([]bool, tt.w*tt.h)
		b.Run(fmt.Sprintf("window-(H:%d-W:%d)", tt.h, tt.w), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l = NextGrid(tt.w, tt.h, l)
			}
		})
		dumpList = l
	}
}
