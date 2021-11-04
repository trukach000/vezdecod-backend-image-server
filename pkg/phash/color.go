package phash

import (
	"image"
)

// func ToGrayscale(img *image.RGBA) *image.RGBA {

// 	b := img.Bounds()
// 	imgSet := image.NewRGBA(b)
// 	for y := 0; y < b.Max.Y; y++ {
// 		for x := 0; x < b.Max.X; x++ {
// 			imgSet.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
// 		}
// 	}

// 	return imgSet
// }

func ToGrayscalePixels(img *image.RGBA) [][]float64 {

	b := img.Bounds()
	w, h := b.Max.X-b.Min.X, b.Max.Y-b.Min.Y
	pixels := make([][]float64, h)

	for i := range pixels {
		pixels[i] = make([]float64, w)
		for j := range pixels[i] {
			color := img.At(j, i)
			r, g, b, _ := color.RGBA()
			s := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
			pixels[i][j] = s
		}
	}

	return pixels
}
