package resizer

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

type ImageType string

const (
	ImageTypePNG ImageType = "image/png"
	ImageTypeJPG ImageType = "image/jpeg"
)

// Resize file
func Resize(f io.ReadCloser, w, h uint, imageType string) (io.ReadSeeker, error) {
	switch ImageType(imageType) {
	case ImageTypeJPG:
		img, err := jpeg.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("decode jpeg: %w", err)
		}

		w, h = newWidthHeight(w, h)
		resized := resize.Resize(w, h, img, resize.Lanczos3)

		buff := bytes.NewBuffer([]byte{})
		if err := jpeg.Encode(buff, resized, &jpeg.Options{Quality: 100}); err != nil {
			return nil, fmt.Errorf("encode jpeg: %w", err)
		}

		return bytes.NewReader(buff.Bytes()), nil

	case ImageTypePNG:
		img, _, err := image.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("decode png: %w", err)
		}

		w, h = newWidthHeight(w, h)
		resized := resize.Resize(w, h, img, resize.Lanczos3)

		buff := bytes.NewBuffer([]byte{})
		if err := png.Encode(buff, resized); err != nil {
			return nil, fmt.Errorf("encode png: %w", err)
		}

		return bytes.NewReader(buff.Bytes()), nil
	}

	return nil, errors.New("image format must be *.png or *.jpg")
}

// TODO: handle case when origin sizes less then new ones
func newWidthHeight(w, h uint) (uint, uint) {
	if w >= h {
		return w, 0
	}
	return 0, h
}
