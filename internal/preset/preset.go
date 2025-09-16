package preset

import (
	"MemeCraft/internal/service/imageutil"
	"image"
)

type TextBox struct {
	Name        string  `json:"name"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Font        string  `json:"font,omitempty"`
	Size        float64 `json:"size,omitempty"`
	Color       string  `json:"color,omitempty"`
	Align       string  `json:"align,omitempty"` // "left", "center", "right"
	LineSpacing float64 `json:"line_spacing,omitempty"`
	Padding     int     `json:"padding,omitempty"` // jarak dari tepi
	MaxChars    int     `json:"max_chars,omitempty"`
	Normalize   string  `json:"normalize,omitempty"` // "normal", "toupper", "tolower"
}

type Overlay struct {
	X      int     `json:"x"`
	Y      int     `json:"y"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Rotate float64 `json:"rotate"`
}

type Preset struct {
	Name             string               `json:"name"`
	ID               string               `json:"id"`
	BaseImage        string               `json:"base_image"`
	ExampleImage     string               `json:"example_image"`
	ResizeMode       imageutil.ResizeMode `json:"resize_mode"`
	BaseImageDecoded image.Image          `json:"-"`
	Overlay          Overlay              `json:"overlay"`
	TextBoxes        []TextBox            `json:"text_boxes"`
}
