package resizer

import (
	"bytes"
	"image"
	"image/png"
	"io"

	"github.com/anthonynsimon/bild/transform"
)

// Resize file
func Resize(f io.ReadCloser, w, h int) (io.ReadSeeker, error) {
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	resized := transform.Resize(img, w, h, transform.Linear)

	buff := bytes.NewBuffer([]byte{})
	if err := png.Encode(buff, resized); err != nil {
		return nil, err
	}

	return bytes.NewReader(buff.Bytes()), nil
}
