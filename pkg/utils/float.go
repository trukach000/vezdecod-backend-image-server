package utils

import "math"

const float64EqualityThreshold = 1e-9

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func EqualWithPrecision(a, b, precision float64) bool {
	return math.Abs(a-b) <= precision
}
