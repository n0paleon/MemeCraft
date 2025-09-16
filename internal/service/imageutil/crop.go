package imageutil

import (
	"github.com/disintegration/imaging"
	"image"
)

func Crop(img image.Image, x, y, w, h int) image.Image {
	return imaging.Crop(img, image.Rect(x, y, x+w, y+h))
}
