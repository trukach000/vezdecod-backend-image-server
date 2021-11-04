package phash

func Flatten(arr2d [][]float64, x int, y int) []float64 {
	flattenArr := make([]float64, x*y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattenArr[y*i+j] = arr2d[i][j]
		}
	}
	return flattenArr
}
