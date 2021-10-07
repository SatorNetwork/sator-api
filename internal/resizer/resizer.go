package resizer

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"github.com/nfnt/resize"
)

type ImageType string

const (
	ImageTypePNG ImageType = ".png"
	ImageTypeJPG ImageType = ".jpg"
)

// Resize file
func Resize(f io.ReadCloser, w, h uint, imageType ImageType) (io.ReadSeeker, error) {
	switch imageType {
	case ImageTypeJPG:
		img, err := jpeg.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		resized := resize.Thumbnail(w, h, img, resize.Bilinear)

		buff := bytes.NewBuffer([]byte{})
		if err := png.Encode(buff, resized); err != nil {
			return nil, err
		}

		return bytes.NewReader(buff.Bytes()), nil

	case ImageTypePNG:
		img, _, err := image.Decode(f)
		if err != nil {
			return nil, err
		}
		resized := resize.Thumbnail(w, h, img, resize.Bilinear)

		buff := bytes.NewBuffer([]byte{})
		if err := png.Encode(buff, resized); err != nil {
			return nil, err
		}

		return bytes.NewReader(buff.Bytes()), nil
	}

	return nil, errors.New("image format must be *.png or *.jpg")
}
