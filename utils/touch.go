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
	isPinch, isPan   bool
}

type Pinch struct {
	ID1, ID2 ebiten.TouchID

	OriginDistance float64
	Distance       float64

	CenterX, CenterY int
}

func (p Pinch) IsInward() bool {
	return p.OriginDistance > p.Distance
}

func (p Pinch) IsOutward() bool {
	return p.OriginDistance < p.Distance
}

type TwoFingerPan struct {
	ID1, ID2 ebiten.TouchID

	LastX, LastY     int
	OriginX, OriginY int

	isHorizontal bool
}

func (p TwoFingerPan) IsHorizontal() bool {
	return p.isHorizontal
}

func (p TwoFingerPan) IsVertical() bool {
	return !p.isHorizontal
}

type Tap struct {
	X, Y int
}

type TouchTracker struct {
	touchIDs []ebiten.TouchID
	touches  map[ebiten.TouchID]*touch
	pinch    *Pinch
	pan      *TwoFingerPan
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

	// Handle released touches in this frame
	for id, t := range tt.touches {
		if inpututil.IsTouchJustReleased(id) {
			// clear pinch if part of it was released
			if tt.pinch != nil && (id == tt.pinch.ID1 || id == tt.pinch.ID2) {
				tt.pinch = nil
			}

			// clear pan if part of it was released
			if tt.pan != nil && (id == tt.pan.ID1 || id == tt.pan.ID2) {
				tt.pan = nil
			}

			// If this one has not been touched long (30 frames can be assumed
			// to be 500ms), or moved far, then record tap.
			diff := distance2d(t.originX, t.originY, t.currX, t.currY)
			if !t.isPinch && !t.isPan && (t.duration <= 30 || diff < 2) {
				tt.taps = append(tt.taps, Tap{
					X: t.currX,
					Y: t.currY,
				})
			}

			delete(tt.touches, id)
		}
	}

	// Store new touches in this frame
	tt.touchIDs = inpututil.AppendJustPressedTouchIDs(tt.touchIDs[:0])
	for _, id := range tt.touchIDs {
		x, y := ebiten.TouchPosition(id)
		tt.touches[id] = &touch{
			originX: x, originY: y,
			currX: x, currY: y,
		}
	}

	// Store all touchIDs (new and old) in this frame
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
		if tt.pan == nil && math.Abs(originDiff-currDiff) > 10 {
			if tt.pinch == nil {
				t1.isPinch = true
				t2.isPinch = true
				tt.pinch = &Pinch{
					ID1:            id1,
					ID2:            id2,
					OriginDistance: originDiff,
					Distance:       currDiff,
					CenterX:        (t1.currX + t2.currX) / 2,
					CenterY:        (t1.currY + t2.currY) / 2,
				}
			} else {
				tt.pinch.Distance = currDiff
			}
		}

		// If the distance between the fingers did not change significantly, this is
		// potentially a new two-finger horizontal pan. We need to check that one finger
		// moved horizontally by an arbitraty margin
		id, id2 := tt.touchIDs[0], tt.touchIDs[1]
		t, t2 := tt.touches[id], tt.touches[1]
		diffX := distance(t.originX, t.currX)
		diffY := distance(t.originY, t.currY)
		if tt.pinch == nil {
			if tt.pan == nil && (math.Abs(diffX) > 10 || math.Abs(diffY) > 10) {
				t.isPan = true
				t2.isPan = true
				tt.pan = &TwoFingerPan{
					ID1:          id,
					ID2:          id2,
					OriginX:      t.originX,
					LastX:        t.currX,
					OriginY:      t.originY,
					LastY:        t.currY,
					isHorizontal: math.Abs(diffX) > 10,
				}
			} else if tt.pan != nil {
				if tt.pan.IsHorizontal() {
					tt.pan.LastX = t.currX
				} else {
					tt.pan.LastY = t.currY
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

func (tt *TouchTracker) TappedThree() (Tap, Tap, Tap, bool) {
	if len(tt.taps) == 3 {
		return tt.taps[0], tt.taps[1], tt.taps[2], true
	}
	return Tap{}, Tap{}, Tap{}, false
}

func (tt *TouchTracker) TappedTwo() (Tap, Tap, bool) {
	if len(tt.taps) == 2 {
		return tt.taps[0], tt.taps[1], true
	}
	return Tap{}, Tap{}, false
}

func (tt *TouchTracker) Tapped() (Tap, bool) {
	if len(tt.taps) == 1 {
		return tt.taps[0], true
	}
	return Tap{}, false
}

func (tt *TouchTracker) Pan() *TwoFingerPan {
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
