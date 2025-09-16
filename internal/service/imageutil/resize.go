package imageutil

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
)

func ResizeWithLockRatio(img image.Image, width int) image.Image {
	return imaging.Resize(img, width, 0, imaging.Lanczos)
}

func ResizeWithoutLockRatio(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

func ResizeWithFill(img image.Image, width, height int) image.Image {
	return imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
}

type ResizeMode string

const (
	LockRatio ResizeMode = "lock"    // 1:1 ratio
	Stretch   ResizeMode = "stretch" // no lock aspect ratio
	Fill      ResizeMode = "fill"    // resize with fill
)

// ResizeWithMode resize with custom user mode
func ResizeWithMode(img image.Image, mode ResizeMode, width, height int) (image.Image, error) {
	switch mode {
	case LockRatio:
		return ResizeWithLockRatio(img, width), nil
	case Stretch:
		return ResizeWithoutLockRatio(img, width, height), nil
	case Fill:
		return ResizeWithFill(img, width, height), nil
	default:
		return nil, fmt.Errorf("invalid resize mode: %s", mode)
	}
}
