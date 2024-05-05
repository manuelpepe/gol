package utils

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type touch struct {
	originX, originY int
	currX, currY     int
	duration         int
}

type tap struct {
	X, Y int
}

type TouchTracker struct {
	touchIDs []ebiten.TouchID
	touches  map[ebiten.TouchID]*touch
	taps     []tap
}

func NewTouchTracker() *TouchTracker {
	return &TouchTracker{
		touchIDs: make([]ebiten.TouchID, 0),
	}
}

func (tt *TouchTracker) Update() {
	tt.touchIDs = ebiten.AppendTouchIDs(tt.touchIDs[:0])
}

func (tt *TouchTracker) IsTripleTap() bool {
	return len(tt.touchIDs) == 3
}

func (tt *TouchTracker) IsDoubleTap() bool {
	return len(tt.touchIDs) == 2
}

func (tt *TouchTracker) IsTouching() bool {
	return len(tt.touchIDs) > 0
}

// Return X, Y coordinates of the first touch recorded, if any.
func (tt *TouchTracker) GetFirstTouchPosition() (int, int, bool) {
	if len(tt.touchIDs) > 0 {
		x, y := ebiten.TouchPosition(tt.touchIDs[0])
		slog.Info(fmt.Sprintf("touch id: %+v  x=%d  y=%d", tt.touchIDs[0], x, y))
		return x, y, true
	}
	return -1, -1, false
}
