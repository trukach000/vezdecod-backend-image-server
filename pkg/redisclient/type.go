package redisclient

type ImgData struct {
	W           int64   `json:"w"`
	H           int64   `json:"h"`
	AspectRatio float64 `json:"aspectRatio"`
	Token       string  `json:"token"`
}
