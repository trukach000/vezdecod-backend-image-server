package main

import (
	"backend-image-server/pkg/phash"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	grayGeneration()
	swapColorGeneration()
	rgbNoiseGeneration()
	log.Printf("ALL IMAGES GENERATED")

}

func grayGeneration() {
	existingImageFile, err := os.Open("./p1_1.jpg")
	if err != nil {
		panic(err)
	}
	defer existingImageFile.Close()

	img, err := jpeg.Decode(existingImageFile)
	if err != nil {
		panic(err)
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, img.Bounds(), img, img.Bounds().Min, draw.Src)
	rgbaImg := newImg

	graysaledPixes := phash.ToGrayscalePixels(rgbaImg)
	resImg := image.NewRGBA(img.Bounds())
	for i := 0; i < img.Bounds().Dx(); i++ {
		for j := 0; j < img.Bounds().Dy(); j++ {
			resImg.Set(i, j, color.Gray{Y: uint8(graysaledPixes[j][i])})
		}
	}

	fres, err := os.Create("p1_2.jpg")
	if err != nil {
		panic(err)
	}
	defer fres.Close()
	if err = jpeg.Encode(fres, resImg, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	log.Printf("GRAY IMAGE GENERATED")
}

func swapColorGeneration() {
	existingImageFile, err := os.Open("./p2_1.jpg")
	if err != nil {
		panic(err)
	}
	defer existingImageFile.Close()

	img, err := jpeg.Decode(existingImageFile)
	if err != nil {
		panic(err)
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, img.Bounds(), img, img.Bounds().Min, draw.Src)
	rgbaImg := newImg

	resImg := image.NewRGBA(img.Bounds())
	for i := 0; i < img.Bounds().Dx(); i++ {
		for j := 0; j < img.Bounds().Dy(); j++ {
			r, g, b, a := rgbaImg.At(i, j).RGBA()
			resImg.Set(i, j, color.RGBA{
				R: uint8(g),
				G: uint8(b),
				B: uint8(r),
				A: uint8(a),
			})
		}
	}

	fres, err := os.Create("p2_2.jpg")
	if err != nil {
		panic(err)
	}
	defer fres.Close()
	if err = jpeg.Encode(fres, resImg, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	log.Printf("SWAP COLOR IMAGE GENERATED")
}

func rgbNoiseGeneration() {
	existingImageFile, err := os.Open("./p3_1.jpg")
	if err != nil {
		panic(err)
	}
	defer existingImageFile.Close()

	img, err := jpeg.Decode(existingImageFile)
	if err != nil {
		panic(err)
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, img.Bounds(), img, img.Bounds().Min, draw.Src)
	rgbaImg := newImg

	noisePower := 30

	resImg := image.NewRGBA(img.Bounds())
	for i := 0; i < img.Bounds().Dx(); i++ {
		for j := 0; j < img.Bounds().Dy(); j++ {
			r, g, b, a := rgbaImg.At(i, j).RGBA()
			rd := uint32(rand.Intn(noisePower+noisePower+1) - noisePower)
			gd := uint32(rand.Intn(noisePower+noisePower+1) - noisePower)
			bd := uint32(rand.Intn(noisePower+noisePower+1) - noisePower)
			resImg.Set(i, j, color.RGBA{
				R: uint8(r + rd),
				G: uint8(g + gd),
				B: uint8(b + bd),
				A: uint8(a),
			})
		}
	}

	fres, err := os.Create("p3_2.jpg")
	if err != nil {
		panic(err)
	}
	defer fres.Close()
	if err = jpeg.Encode(fres, resImg, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	log.Printf("RGB NOISED IMAGE GENERATED")
}
