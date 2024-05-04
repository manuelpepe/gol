package gol

import (
	"fmt"
	"testing"
)

func TestCountNeighsBlock(t *testing.T) {
	X := 10
	Y := 10
	onPos := []int{10, 11, 20, 21}
	ts := []struct {
		p   int
		exp int
	}{
		{p: 10, exp: 3},
		{p: 11, exp: 3},
		{p: 20, exp: 3},
		{p: 21, exp: 3},
		{p: 0, exp: 2},
		{p: 1, exp: 2},
		{p: 2, exp: 1},
		{p: 12, exp: 2},
		{p: 32, exp: 1},
		{p: 31, exp: 2},
		{p: 30, exp: 2},
		{p: 33, exp: 0},
		{p: 24, exp: 0},
		{p: 23, exp: 0},
	}

	for ix, tt := range ts {
		t.Run(fmt.Sprintf("%d - pos: %d", ix, tt.p), func(t *testing.T) {
			grid := make([]bool, X*Y)
			for _, v := range onPos {
				grid[v] = true
			}
			n := CalcNeighs(tt.p, X, Y, grid)
			if n != tt.exp {
				t.Errorf("expected %d to equal %d on pos %d", n, tt.exp, tt.p)
			}
		})
	}
}
