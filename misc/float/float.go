package float

import "math"

const (
	floatThreshold = 1e-12
)

func Equal(a, b float64) bool {
	return math.Abs(a-b) < floatThreshold
}
