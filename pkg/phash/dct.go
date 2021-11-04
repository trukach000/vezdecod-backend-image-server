package phash

import (
	"math"
	"math/cmplx"
	"sync"

	"github.com/mjibson/go-dsp/fft"
)

func DCT(x []float64) []float64 {
	n := len(x)
	n2 := n / 2

	evenVector := make([]float64, n)
	for i := 0; i < n2; i++ {
		evenVector[i], evenVector[n-1-i] = x[2*i], x[2*i+1]
	}

	array := fft.FFTReal(evenVector)
	theta := math.Pi / (2.0 * float64(n))

	for j := 1; j < n2; j++ {
		w := cmplx.Exp(complex(0, -float64(j)*theta))
		wCont := -complex(imag(w), real(w))
		array[j] = array[j] * w
		array[n-j] = array[n-j] * wCont
	}
	array[n2] = array[n2] * complex(math.Cos(theta*float64(n2)), 0.0)

	dctK := make([]float64, n)
	for i := range dctK {
		dctK[i] = real(array[i])
	}

	return dctK
}

func DCTPixels(input [][]float64, w int, h int) [][]float64 {
	output := make([][]float64, h)
	for i := range output {
		output[i] = make([]float64, w)
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < h; i++ {
		wg.Add(1)
		go func(i int) {
			cols := DCT(input[i])
			output[i] = cols
			wg.Done()
		}(i)
	}

	wg.Wait()
	for i := 0; i < w; i++ {
		wg.Add(1)
		in := make([]float64, h)
		go func(i int) {
			for j := 0; j < h; j++ {
				in[j] = output[j][i]
			}
			rows := DCT(in)
			for j := 0; j < len(rows); j++ {
				output[j][i] = rows[j]
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	return output
}
