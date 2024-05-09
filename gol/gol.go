package gol

func NextGrid(X, Y int, grid []bool) []bool {
	if len(grid) != X*Y {
		panic("grid doesn't match dimensions")
	}
	next := make([]bool, X*Y)
	for p := 0; p < X*Y; p++ {
		neighs := CalcNeighs(p, X, Y, grid)

		// Any live cell with fewer than two live neighbors dies, as if by underpopulation.
		// Any live cell with two or three live neighbors lives on to the next generation.
		// Any live cell with more than three live neighbors dies, as if by overpopulation.
		// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
		next[p] = grid[p]
		if neighs < 2 {
			next[p] = false
		} else if neighs > 3 {
			next[p] = false
		} else if neighs == 3 {
			next[p] = true
		}
	}
	return next
}

var deltas = [3]int{-1, 0, 1}

// Grid operations:
//
//	p-X-1   p-X   p-X+1
//	 p-1     p     p+1
//	p+X-1   p+X   p+X+1
//
// Limits:
//
//	X = total cols
//	Y = total rows
//
//	- p % X == 0       => left
//	- p % X == X-1     => right
//	- p < Y            => up
//	- p >= X * Y - X   => down
func CalcNeighs(p, X, Y int, grid []bool) int {
	neighs := 0
	for _, yd := range deltas {
		// skip top side check at boundry
		if p < Y && yd == -1 {
			continue
		}

		// skip bottom side check at boundry
		if p >= X*Y-X && yd == 1 {
			continue
		}

		for _, xd := range deltas {

			// skip left side check at boundry
			if p%X == 0 && xd == -1 {
				continue
			}

			// skip right side check at boundry
			if p%X == X-1 && xd == 1 {
				continue
			}

			neighIx := p + xd + yd*X
			if neighIx < 0 || neighIx >= X*Y || neighIx == p {
				continue
			}

			if grid[neighIx] {
				neighs += 1
			}
		}
	}
	return neighs
}
