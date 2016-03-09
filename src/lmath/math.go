package lmath

import (
	"math"
)

var HALF_DEG2RAD float32 = 0.0087266462599716478846184538424431

var EPSILON float32  = math.Nextafter32(1, 2) - 1

func AlmostEqual(x float32, y float32, tol float32) bool {
	if x == y || (math.Abs(float64(x-y)) <= float64(tol)*math.Max(1.0, math.Max(math.Abs(float64(x)), math.Abs(float64(y))))) {
		return true
	}
	return false
}

