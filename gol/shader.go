package gol

import (
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed shader.kage
var shaderProgram []byte

//go:embed shaderv2.kage
var shaderProgramV2 []byte

var shaderInstance *ebiten.Shader
var shaderV2Instance *ebiten.Shader

func init() {
	shader, err := ebiten.NewShader(shaderProgram)
	if err != nil {
		log.Fatal(err)
	}
	shaderInstance = shader

	shaderv2, err := ebiten.NewShader(shaderProgramV2)
	if err != nil {
		log.Fatal(err)
	}
	shaderV2Instance = shaderv2
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

var shaderOutputBuffer []byte
var shaderOutputBufferX, shaderOutputBufferY int

// in this version grid is modified
func NextGridShaderV2(X, Y int, grid []bool) []bool {
	if X%4 != 0 {
		panic("grid length must be divisible by 4 to use shader v2")
	}

	if shaderOutputBuffer == nil || X != shaderOutputBufferX || Y != shaderOutputBufferY {
		shaderOutputBuffer = make([]byte, X*Y)
		shaderOutputBufferX = X
		shaderOutputBufferY = Y
	}

	outImage := ebiten.NewImage(X/4, Y)
	gridImage := ebiten.NewImage(X/4, Y)

	// encode current grid
	for y := range Y {
		xPacketIx := 0
		xPacketColor := color.RGBA{}
		for x := range X {
			if grid[y*Y+x] {
				switch x % 4 {
				case 0:
					xPacketColor.R = 255
				case 1:
					xPacketColor.G = 255
				case 2:
					xPacketColor.B = 255
				case 3:
					xPacketColor.A = 255
				}
			}
			if x%4 == 3 {
				gridImage.Set(xPacketIx, y, xPacketColor)
				xPacketIx += 1
				xPacketColor = color.RGBA{}
			}
		}
	}

	// calculate with shader
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = gridImage
	outImage.DrawRectShader(X/4, Y, shaderV2Instance, opts)

	// decode shader output
	outImage.ReadPixels(shaderOutputBuffer)
	for ix, b := range shaderOutputBuffer {
		grid[ix] = b > 0
	}

	return grid
}
