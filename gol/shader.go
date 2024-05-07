package gol

import (
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shader.kage
var shaderProgram []byte

var shaderInstance *ebiten.Shader

func init() {
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil {
		log.Fatal(err)
	}
	shaderInstance = shader
}

func NextGridShader(X, Y int, grid []bool) []bool {
	outImage := ebiten.NewImage(X, Y)
	gridImage := ebiten.NewImage(X, Y)

	// encode current grid
	for y := range Y {
		for x := range X {
			if grid[y*Y+x] {
				gridImage.Set(x, y, color.RGBA{A: 255})
			}
		}
	}

	// calculate with shader
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = gridImage
	outImage.DrawRectShader(X, Y, shaderInstance, opts)

	// decode shader output
	data := make([]byte, 4*Y*X)
	outImage.ReadPixels(data)
	nextGrid := make([]bool, X*Y)
	gix := 0
	for i := 0; i < 4*X*Y; i += 4 {
		alphaIx := i + 3 // value is encoded in alpha channel
		nextGrid[gix] = data[alphaIx] > 0
		gix++
	}

	return nextGrid
}
