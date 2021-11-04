package resize

import (
	"backend-image-server/pkg/utils"
	"image"
	"image/draw"
	"math"

	"github.com/sirupsen/logrus"
)

func Resize(img image.Image, scale float64) image.Image {

	newHeight := int(math.Ceil(float64(img.Bounds().Dy()) * scale))
	newWidth := int(math.Ceil(float64(img.Bounds().Dx()) * scale))

	if utils.AlmostEqual(scale, 1.0) {
		return img
	}

	if img.Bounds().Dx() <= 0 || img.Bounds().Dy() <= 0 {
		return img
	}

	logrus.Infof("Old height: %d", img.Bounds().Dy())
	logrus.Infof("Old width: %d", img.Bounds().Dx())

	logrus.Infof("New height: %d", newHeight)
	logrus.Infof("New width: %d", newWidth)

	var rgbaImg *image.RGBA
	if rgba, ok := img.(*image.RGBA); ok {
		rgbaImg = rgba
		logrus.Info("RGBA already")
	} else {
		newImg := image.NewRGBA(img.Bounds())
		draw.Draw(newImg, img.Bounds(), img, img.Bounds().Min, draw.Src)
		rgbaImg = newImg
	}

	var dst *image.RGBA

	dst = resampleHorizontal(rgbaImg, newWidth)
	dst = resampleVertical(dst, newHeight)

	return dst
}
