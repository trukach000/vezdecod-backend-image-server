package app

import (
	"backend-image-server/app/repositories/images"
	"backend-image-server/pkg/httpext"
	"backend-image-server/pkg/phash"
	"backend-image-server/pkg/redisclient"
	"backend-image-server/pkg/resize"
	"backend-image-server/pkg/utils"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
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

	img, err := jpeg.Decode(bytes.NewReader(fileBytesBuffer))
	if err != nil {
		logrus.Errorf("Can't decode image to jpeg: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Can't decode image to jpeg",
		}, http.StatusInternalServerError)
		return
	}

	newWidth := img.Bounds().Dx()
	newHeight := img.Bounds().Dy()
	newAspectRatio := float64(newHeight) / float64(newWidth)
	logrus.Infof("newWidth: %d", newWidth)
	logrus.Infof("newHeight: %d", newHeight)
	logrus.Infof("newAspectRatio: %f", newAspectRatio)

	hashP := phash.PHash(img)
	hashPString := fmt.Sprintf("%x", hashP)

	// find similar images by p-hash

	redisClient, err := redisclient.GetRedisFromContext(ctx)
	if err != nil {
		logrus.Errorf("Can't get redis from context: %s", err)
		httpext.AbortJSON(w, httpext.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}, http.StatusInternalServerError)
		return
	}

	token := ""

	valueRaw, err := redisClient.Get(images.GetRedisKey(hashPString, newAspectRatio)).Result()
	if err != redis.Nil {
		// there is a similar image in redis
		// check ration and in case of new img is larger change it, otherwise just return the old one

		var data redisclient.ImgData
		err := json.Unmarshal([]byte(valueRaw), &data)
		if err != nil {
			logrus.Errorf("Can't get data from redis (unmarshall error): %s", err)
			httpext.AbortJSON(w, httpext.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}, http.StatusInternalServerError)
			return
		}

		logrus.Infof("There is an image with the same pHash: %+v", data)

		if utils.AlmostEqual(newAspectRatio, data.AspectRatio) {
			logrus.Infof("Ratios are equals, check the sizes")
			logrus.Infof("Old width %d, new width: %d", data.W, newWidth)
			if newWidth > int(data.W) {
				logrus.Infof("The new one is larger, upload new one")

				err := images.ReplaceImage(
					ctx, data.Token, hashPString,
					newWidth, newHeight, newAspectRatio,
					fileBytesBuffer,
				)
				if err != nil {
					logrus.Errorf("Can't replace image into database: %s", err)
					httpext.AbortJSON(w, httpext.ErrorResponse{
						Code:    http.StatusInternalServerError,
						Message: "Can't replace image into database",
					}, http.StatusInternalServerError)
					return
				}
				token = data.Token
			} else {
				logrus.Infof("The old one is larger or equal, return the old token")
				token = data.Token
			}
		} else {
			logrus.Infof("New image has different aspect ratio: %f, save it", newAspectRatio)
			token, err = images.SaveNewImage(ctx, hashPString, newWidth, newHeight, newAspectRatio, fileBytesBuffer)
			if err != nil {
				logrus.Errorf("Can't save image into database: %s", err)
				httpext.AbortJSON(w, httpext.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Can't save image into database",
				}, http.StatusInternalServerError)
				return
			}
		}
	} else {
		token, err = images.SaveNewImage(ctx, hashPString, newWidth, newHeight, newAspectRatio, fileBytesBuffer)
		if err != nil {
			logrus.Errorf("Can't save image into database: %s", err)
			httpext.AbortJSON(w, httpext.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Can't save image into database",
			}, http.StatusInternalServerError)
			return
		}
	}

	httpext.JSON(w, UploadResponse{
		ImageToken: token,
		PHash:      hashPString,
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
