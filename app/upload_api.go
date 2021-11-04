package app

import (
	"backend-image-server/app/repositories/images"
	"backend-image-server/pkg/httpext"
	"backend-image-server/pkg/phash"
	"backend-image-server/pkg/resize"
	"backend-image-server/pkg/utils"
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type UploadResponse struct {
	ImageToken string `json:"imageToken"`
	PHash      string `json:"pHash"`
}

// Upload godoc
// @Summary Upload jpg image
// @Description upload jpg image into MySQL database (max allowed size - 50 mb)
// @Param image formData file true "image to upload"
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} UploadResponse
// @Failure 400 {object} httpext.ErrorResponse
// @Failure 500 {object} httpext.ErrorResponse
// @Router /upload [post]
func UploadImage(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	err := r.ParseMultipartForm(50 << 20) // 50 mbs max
	if err != nil {
		logrus.Errorf("Can't parse multipart form: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Can't parse multipart form: %s", err),
		}, http.StatusInternalServerError)
		return
	}
	file, handler, err := r.FormFile("image")
	if err != nil {
		logrus.Errorf("Can't parse multipart form: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Can't retrieve image from form data: %s", err),
		}, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileBytesBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("Can't  read bytes from files: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}, http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(fileBytesBuffer)
	if contentType != "image/jpeg" {
		logrus.Errorf("Wrong MIME type: %s", contentType)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Wrong file MIME type (jpeg ONLY)",
		}, http.StatusBadRequest)
		return
	}

	logrus.Infof("File Size: %d", handler.Size)

	token, err := images.SaveNewImage(ctx, fileBytesBuffer)
	if err != nil {
		logrus.Errorf("Can't save image into database: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Can't save image into database",
		}, http.StatusInternalServerError)
		return
	}

	img, err := jpeg.Decode(bytes.NewReader(fileBytesBuffer))
	if err != nil {
		logrus.Errorf("Can't decode image to jpeg: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Can't decode image to jpeg",
		}, http.StatusInternalServerError)
		return
	}
	hashP := phash.PHash(img)

	// find similar images by p-hash

	httpext.JSON(w, UploadResponse{
		ImageToken: token,
		PHash:      fmt.Sprintf("%x", hashP),
	})
}

// GetImage godoc
// @Summary get image by its id
// @Description return image by its id
// @Param id path string true "image id"
// @Param scale query number false "scale coeff"
// @Produce jpeg
// @Success 200 {string} image/png
// @Failure 404 {object} httpext.ErrorResponse
// @Failure 500 {object} httpext.ErrorResponse
// @Router /get/{id} [get]
func GetImage(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	imageToken := chi.URLParam(r, "id")
	scaleStr := r.URL.Query().Get("scale")

	scale := 1.0
	if scaleStr != "" {
		var err error
		scale, err = strconv.ParseFloat(scaleStr, 64)
		if err != nil {
			logrus.Errorf("Can't parse scale parameter: %s", scaleStr)
		}
	}

	imageData, err := images.GetImageByToken(ctx, imageToken)
	if err == images.ErrImageNotFound {
		logrus.Errorf("Not found image with token: %s", imageToken)
		httpext.AbortWithoutContent(w, http.StatusNotFound)
		return
	}
	if err != nil {
		logrus.Errorf("Can't get image from database: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Can't get image from database",
		}, http.StatusInternalServerError)
		return
	}

	if !utils.AlmostEqual(scale, 1.0) {
		logrus.Infof("Resizing image with scale: %f", scale)
		img, err := jpeg.Decode(bytes.NewReader(imageData))
		if err != nil {
			logrus.Errorf("Can't decode image to jpeg: %s", err)
			httpext.AbortJSON(w, httpext.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Can't decode image to jpeg",
			}, http.StatusInternalServerError)
			return
		}

		newHeight := int(math.Ceil(float64(img.Bounds().Dy()) * scale))
		newWidth := int(math.Ceil(float64(img.Bounds().Dx()) * scale))

		m := resize.Resize(img, newHeight, newWidth)
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		jpeg.Encode(writer, m, nil)

		imageData = buf.Bytes()
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imageData)))
	_, err = w.Write(imageData)
	if err != nil {
		logrus.Errorf("Can't write image data to response: %s", err)
	}
}
