package app

import (
	"backend-image-server/pkg/httpext"
	"backend-image-server/pkg/phash"
	"backend-image-server/pkg/utils"
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type CompareResponse struct {
	IsSimilar       bool   `json:"isSimilar"`
	PHash1          string `json:"pHash1"`
	PHash2          string `json:"pHash2"`
	HammingDistance int64  `json:"hammingDestance"`
}

// Compare godoc
// @Summary compare two jpg image
// @Description compare two jpg image (max allowed size - 50 mb)
// @Param image1 formData file true "image1 to compare"
// @Param image2 formData file true "image2 to compare"
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} CompareResponse
// @Failure 500 {object} httpext.ErrorResponse
// @Router /compare [post]
func CompareImage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(50 << 20) // 50 mbs max
	if err != nil {
		logrus.Errorf("Can't parse multipart form: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Can't parse multipart form: %s", err),
		}, http.StatusInternalServerError)
		return
	}
	file1, _, err := r.FormFile("image1")
	if err != nil {
		logrus.Errorf("Can't parse multipart form: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Can't retrieve image from form data: %s", err),
		}, http.StatusInternalServerError)
		return
	}
	defer file1.Close()

	file2, _, err := r.FormFile("image2")
	if err != nil {
		logrus.Errorf("Can't parse multipart form: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Can't retrieve image from form data: %s", err),
		}, http.StatusInternalServerError)
		return
	}
	defer file2.Close()

	fileBytesBuffer1, err := ioutil.ReadAll(file1)
	if err != nil {
		logrus.Errorf("Can't  read bytes from files: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}, http.StatusInternalServerError)
		return
	}

	fileBytesBuffer2, err := ioutil.ReadAll(file2)
	if err != nil {
		logrus.Errorf("Can't  read bytes from files: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}, http.StatusInternalServerError)
		return
	}

	contentType1 := http.DetectContentType(fileBytesBuffer1)
	contentType2 := http.DetectContentType(fileBytesBuffer2)

	if contentType1 != "image/jpeg" || contentType2 != "image/jpeg" {
		logrus.Errorf("Wrong MIME type: %s, %s", contentType1, contentType2)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Wrong file MIME type (jpeg ONLY)",
		}, http.StatusBadRequest)
		return
	}

	img1, err := jpeg.Decode(bytes.NewReader(fileBytesBuffer1))
	if err != nil {
		logrus.Errorf("Can't decode image to jpeg: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Can't decode image to jpeg",
		}, http.StatusInternalServerError)
		return
	}

	img2, err := jpeg.Decode(bytes.NewReader(fileBytesBuffer2))
	if err != nil {
		logrus.Errorf("Can't decode image to jpeg: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Can't decode image to jpeg",
		}, http.StatusInternalServerError)
		return
	}

	width1 := img1.Bounds().Dx()
	height1 := img1.Bounds().Dy()
	aspectRatio1 := float64(height1) / float64(width1)
	logrus.Infof("width1: %d", width1)
	logrus.Infof("height1: %d", height1)
	logrus.Infof("aspectRatio1: %f", aspectRatio1)

	width2 := img2.Bounds().Dx()
	height2 := img2.Bounds().Dy()
	aspectRatio2 := float64(height2) / float64(width2)
	logrus.Infof("width2: %d", width2)
	logrus.Infof("height2: %d", height2)
	logrus.Infof("aspectRatio2: %f", aspectRatio2)

	hashP1 := phash.PHash(img1)
	hashPString1 := fmt.Sprintf("%x", hashP1)

	hashP2 := phash.PHash(img2)
	hashPString2 := fmt.Sprintf("%x", hashP2)

	if !utils.EqualWithPrecision(aspectRatio1, aspectRatio2, 0.01) {
		httpext.JSON(w, CompareResponse{
			IsSimilar:       false,
			PHash1:          hashPString1,
			PHash2:          hashPString2,
			HammingDistance: -1,
		})
		return
	}

	hammingDistance := phash.Hamming(hashP1, hashP2)

	const similarityTreashold = 3

	if hammingDistance <= similarityTreashold {
		httpext.JSON(w, CompareResponse{
			IsSimilar:       true,
			PHash1:          hashPString1,
			PHash2:          hashPString2,
			HammingDistance: int64(hammingDistance),
		})
		return
	}

	httpext.JSON(w, CompareResponse{
		IsSimilar:       false,
		PHash1:          hashPString1,
		PHash2:          hashPString2,
		HammingDistance: int64(hammingDistance),
	})

}
