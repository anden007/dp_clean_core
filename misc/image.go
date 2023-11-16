package misc

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

// 缩放图片imgType：jpg/png
func ResizeImage(imgData []byte, width, height uint) (result []byte, err error) {
	Img, _, err := image.Decode(bytes.NewReader(imgData))
	imgInfo := Img.Bounds()
	if imgInfo.Max.X < 750 {
		result = imgData
	} else {
		targetImg := resize.Resize(width, 0, Img, resize.Lanczos2)
		buf := bytes.NewBuffer(result)
		err = jpeg.Encode(buf, targetImg, &jpeg.Options{Quality: 80})
		result = buf.Bytes()
	}
	return
}

// ImageZoom 按宽度缩放图片
func ImageZoom(url string, width uint, m image.Image) (image.Image, error) {

	thImg := resize.Resize(width, 0, m, resize.Lanczos3)
	return thImg, nil
}

// 缩放图片imgType：jpg/png
func ResizeImageByData(imgData []byte, width, height uint) (result []byte, err error) {

	Img, _, err := image.Decode(bytes.NewReader(imgData))
	targetImg := resize.Resize(width, 0, Img, resize.Lanczos2)
	buf := bytes.NewBuffer(result)
	err = jpeg.Encode(buf, targetImg, &jpeg.Options{Quality: 80})
	result = buf.Bytes()

	return
}
