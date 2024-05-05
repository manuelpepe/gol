package utils

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func distance(a, b int) float64 {
	return math.Abs(float64(a - b))
}

// distance2d between points a and b.
func distance2d(xa, ya, xb, yb int) float64 {
	x := distance(xa, xb)
	y := distance(ya, yb)
	return math.Sqrt(x*x + y*y)
}

type touch struct {
	originX, originY int
	currX, currY     int
	duration         int
	wasPinch, isPan  bool
}

type Pinch struct {
	ID1, ID2 ebiten.TouchID
	OriginH  float64
	PrevH    float64
}

type TwoFingerHPan struct {
	ID1, ID2 ebiten.TouchID

	PrevX   int
	OriginX int
}

type Tap struct {
	X, Y int
}

type TouchTracker struct {
	touchIDs []ebiten.TouchID
	touches  map[ebiten.TouchID]*touch
	pinch    *Pinch
	pan      *TwoFingerHPan
	taps     []Tap
}

func NewTouchTracker() *TouchTracker {
	return &TouchTracker{
		touchIDs: make([]ebiten.TouchID, 0),
		taps:     make([]Tap, 0),
		touches:  make(map[ebiten.TouchID]*touch),
	}
}

func (tt *TouchTracker) Update() {
	// Clear the previous frame's taps.
	tt.taps = tt.taps[:0]
	// What touches have just ended?
	for id, t := range tt.touches {
		if inpututil.IsTouchJustReleased(id) {
			if tt.pinch != nil && (id == tt.pinch.ID1 || id == tt.pinch.ID2) {
				tt.pinch = nil
			}
			if tt.pan != nil && (id == tt.pan.ID1 || id == tt.pan.ID2) {
				tt.pan = nil
			}

			// If this one has not been touched long (30 frames can be assumed
			// to be 500ms), or moved far, then it's a tap.
			diff := distance2d(t.originX, t.originY, t.currX, t.currY)
			if !t.wasPinch && !t.isPan && (t.duration <= 30 || diff < 2) {
				tt.taps = append(tt.taps, Tap{
					X: t.currX,
					Y: t.currY,
				})
			}

			delete(tt.touches, id)
		}
	}

	// What touches are new in this frame?
	tt.touchIDs = inpututil.AppendJustPressedTouchIDs(tt.touchIDs[:0])
	for _, id := range tt.touchIDs {
		x, y := ebiten.TouchPosition(id)
		tt.touches[id] = &touch{
			originX: x, originY: y,
			currX: x, currY: y,
		}
	}

	tt.touchIDs = ebiten.AppendTouchIDs(tt.touchIDs[:0])

	// Update the current position and durations of any touches that have
	// neither begun nor ended in this frame.
	for _, id := range tt.touchIDs {
		t := tt.touches[id]
		t.duration = inpututil.TouchPressDuration(id)
		t.currX, t.currY = ebiten.TouchPosition(id)
	}

	// Interpret the raw touch data that's been collected into tt.touches into
	// gestures like two-finger pinch or two-finger pan.
	if len(tt.touches) == 2 {
		// Potentially the user is making a pinch gesture with two fingers.
		// If the diff between their origins is different to the diff between
		// their currents and if these two are not already a pinch, then this is
		// a new pinch!
		id1, id2 := tt.touchIDs[0], tt.touchIDs[1]
		t1, t2 := tt.touches[id1], tt.touches[id2]
		originDiff := distance2d(t1.originX, t1.originY, t2.originX, t2.originY)
		currDiff := distance2d(t1.currX, t1.currY, t2.currX, t2.currY)
		if tt.pinch == nil && tt.pan == nil && math.Abs(originDiff-currDiff) > 10 {
			t1.wasPinch = true
			t2.wasPinch = true
			tt.pinch = &Pinch{
				ID1:     id1,
				ID2:     id2,
				OriginH: originDiff,
				PrevH:   originDiff,
			}
		}

		// If the distance between the fingers did not change significantly, this is
		// potentially a new two-finger horizontal pan. We need to check that one finger
		// moved horizontally by an arbitraty margin
		id, id2 := tt.touchIDs[0], tt.touchIDs[1]
		t, t2 := tt.touches[id], tt.touches[1]
		if !t.wasPinch && tt.pan == nil && tt.pinch == nil {
			diff := distance(t.originX, t.currX)
			if math.Abs(diff) > 10 {
				t.isPan = true
				t2.isPan = true
				tt.pan = &TwoFingerHPan{
					ID1:     id,
					ID2:     id2,
					OriginX: t.originX,
					PrevX:   t.currX,
				}
			}
		}

	}
}

func (tt *TouchTracker) IsTouchingThree() bool {
	return len(tt.touchIDs) == 3
}

func (tt *TouchTracker) IsTouchingTwo() bool {
	return len(tt.touchIDs) == 2
}

func (tt *TouchTracker) IsTouching() bool {
	return len(tt.touchIDs) > 0
}

func (tt *TouchTracker) Pan() *TwoFingerHPan {
	return tt.pan
}

func (tt *TouchTracker) Pinch() *Pinch {
	return tt.pinch
}

// Return X, Y coordinates of the first touch recorded, if any.
func (tt *TouchTracker) GetFirstTouchPosition() (int, int, bool) {
	if len(tt.touchIDs) > 0 {
		x, y := ebiten.TouchPosition(tt.touchIDs[0])
		return x, y, true
	}
	return -1, -1, false
}
