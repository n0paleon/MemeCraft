package imageutil

import (
	"image"
	"image/draw"

	"github.com/disintegration/imaging"
)

// mode: "front" (default) atau "back"
// angle: rotasi overlay dalam derajat (searah jarum jam)
func Overlay(base image.Image, overlay image.Image, x, y int, mode string, angle float64) image.Image {
	if angle != 0 {
		overlay = imaging.Rotate(overlay, angle, image.Transparent)
	}

	rgba := image.NewRGBA(base.Bounds())
	draw.Draw(rgba, base.Bounds(), base, image.Point{}, draw.Src)

	offset := image.Pt(x, y)
	overlayBounds := overlay.Bounds().Add(offset)

	if mode == "back" {
		draw.Draw(rgba, overlayBounds, overlay, image.Point{}, draw.Over)
		draw.Draw(rgba, base.Bounds(), base, image.Point{}, draw.Over)
	} else {
		draw.Draw(rgba, overlayBounds, overlay, image.Point{}, draw.Over)
	}

	return rgba
}
