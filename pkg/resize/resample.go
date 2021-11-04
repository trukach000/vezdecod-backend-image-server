package resize

import (
	"backend-image-server/pkg/utils"
	"image"
	"math"
)

func resampleHorizontal(src *image.RGBA, width int) *image.RGBA {
	srcWidth, srcHeight := src.Bounds().Dx(), src.Bounds().Dy()
	srcStride := src.Stride

	delta := float64(srcWidth) / float64(width)
	scale := math.Max(delta, 1.0)

	dst := image.NewRGBA(image.Rect(0, 0, width, srcHeight))
	dstStride := dst.Stride

	filterRadius := math.Ceil(scale * 3.0)

	utils.Line(srcHeight, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < width; x++ {
				ix := (float64(x)+0.5)*delta - 0.5
				istart, iend := int(ix-filterRadius+0.5), int(ix+filterRadius)

				if istart < 0 {
					istart = 0
				}
				if iend >= srcWidth {
					iend = srcWidth - 1
				}

				var r, g, b, a float64
				var sum float64
				for kx := istart; kx <= iend; kx++ {

					srcPos := y*srcStride + kx*4
					normPos := (float64(kx) - ix) / scale
					fValue := LanczosInterpolation(normPos)

					r += float64(src.Pix[srcPos+0]) * fValue
					g += float64(src.Pix[srcPos+1]) * fValue
					b += float64(src.Pix[srcPos+2]) * fValue
					a += float64(src.Pix[srcPos+3]) * fValue
					sum += fValue
				}

				dstPos := y*dstStride + x*4
				dst.Pix[dstPos+0] = uint8(Keep((r/sum)+0.5, 0, 255))
				dst.Pix[dstPos+1] = uint8(Keep((g/sum)+0.5, 0, 255))
				dst.Pix[dstPos+2] = uint8(Keep((b/sum)+0.5, 0, 255))
				dst.Pix[dstPos+3] = uint8(Keep((a/sum)+0.5, 0, 255))
			}
		}
	})

	return dst
}

func resampleVertical(src *image.RGBA, height int) *image.RGBA {
	srcWidth, srcHeight := src.Bounds().Dx(), src.Bounds().Dy()
	srcStride := src.Stride

	delta := float64(srcHeight) / float64(height)
	scale := math.Max(delta, 1.0)

	dst := image.NewRGBA(image.Rect(0, 0, srcWidth, height))
	dstStride := dst.Stride

	filterRadius := math.Ceil(scale * 3.0)

	utils.Line(height, func(start, end int) {
		for y := start; y < end; y++ {
			iy := (float64(y)+0.5)*delta - 0.5

			istart, iend := int(iy-filterRadius+0.5), int(iy+filterRadius)

			if istart < 0 {
				istart = 0
			}
			if iend >= srcHeight {
				iend = srcHeight - 1
			}

			for x := 0; x < srcWidth; x++ {
				var r, g, b, a float64
				var sum float64
				for ky := istart; ky <= iend; ky++ {

					srcPos := ky*srcStride + x*4
					normPos := (float64(ky) - iy) / scale
					fValue := LanczosInterpolation(normPos)

					r += float64(src.Pix[srcPos+0]) * fValue
					g += float64(src.Pix[srcPos+1]) * fValue
					b += float64(src.Pix[srcPos+2]) * fValue
					a += float64(src.Pix[srcPos+3]) * fValue
					sum += fValue
				}

				dstPos := y*dstStride + x*4
				dst.Pix[dstPos+0] = uint8(Keep((r/sum)+0.5, 0, 255))
				dst.Pix[dstPos+1] = uint8(Keep((g/sum)+0.5, 0, 255))
				dst.Pix[dstPos+2] = uint8(Keep((b/sum)+0.5, 0, 255))
				dst.Pix[dstPos+3] = uint8(Keep((a/sum)+0.5, 0, 255))
			}
		}
	})

	return dst
}
