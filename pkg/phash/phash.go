package phash

import (
	"backend-image-server/pkg/resize"
	"image"

	"github.com/sirupsen/logrus"
)

const (
	reducedSize = 32
)

func PHash(img image.Image) []byte {
	smallImage := resize.Resize(img, reducedSize, reducedSize)
	grayPixels := ToGrayscalePixels(smallImage)
	dctPixels := DCTPixels(grayPixels, 32, 32)

	// reducing to 8x8 dct array
	dctReducesPixes := make([][]float64, 8)
	for i := 0; i < 8; i++ {
		dctReducesPixes[i] = make([]float64, 8)
		for j := 0; j < 8; j++ {
			dctReducesPixes[i][j] = dctPixels[i][j]
		}
	}

	flattenDct := Flatten(dctReducesPixes, 8, 8)

	sum := 0.0
	for i := 0; i < 64; i++ {
		sum += (flattenDct[i])
	}
	dctAvg := (float64(sum)) / (float64(32 * 32))

	logrus.Infof("dctAvg: %+v", dctAvg)

	dctBytes := make([]byte, 8)
	// imgSet := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if dctReducesPixes[i][j] > dctAvg {
				dctBytes[i] |= (1 << j)
				//imgSet.Set(i, j, color.Gray{Y: uint8(255)})
			}
		}
	}

	return dctBytes
}
