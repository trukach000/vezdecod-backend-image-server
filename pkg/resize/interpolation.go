package resize

import "math"

func LanczosInterpolation(x float64) float64 {
	x = math.Abs(x)
	if x == 0 {
		return 1.0
	} else if x < 3.0 {
		return (3.0 * math.Sin(math.Pi*x) * math.Sin(math.Pi*(x/3.0))) / (math.Pi * math.Pi * x * x)
	}
	return 0.0
}

func cubicInt(x float64) float64 {

	x = math.Abs(x)

	if x <= 1 {
		return 1.0 + x*x*(1.5*x-2.5)
	}

	if x <= 2 {
		return 2.0 + x*(x*(2.5-0.5*x)-4.0)
	}

	return 0
}
