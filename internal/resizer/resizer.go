package resizer

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime"

	"github.com/nfnt/resize"
)

type ImageType string

const (
	ImageTypePNG  ImageType = "image/png"
	ImageTypeJPEG ImageType = "image/jpeg"
)

// Resize file
func Resize(f io.ReadCloser, w, h uint) (io.ReadSeeker, error) {
	imageType := guessImageMimeTypes(f)

	switch ImageType(imageType) {
	case ImageTypeJPEG:
		img, err := jpeg.Decode(f)
		if err != nil {
			return nil, err
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

	return nil, errors.New("image format must be PNG or JPG/JPEG")
}

// Guess image format from gif/jpeg/png/webp
func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// Guess image mime types from gif/jpeg/png/webp
func guessImageMimeTypes(r io.Reader) string {
	format, _ := guessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}
