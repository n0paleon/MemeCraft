package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
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

func DrawTextBoxes(base image.Image, text map[string]string, boxes []TextBox) (image.Image, error) {
	bounds := base.Bounds()
	dc := gg.NewContext(bounds.Dx(), bounds.Dy())
	dc.DrawImage(base, 0, 0)

	for _, box := range boxes {
		userText, ok := text[box.Name]
		if !ok || userText == "" {
			continue
		}

		switch strings.ToLower(box.Normalize) {
		case "tolower":
			userText = strings.ToLower(userText)
		case "toupper":
			userText = strings.ToUpper(userText)
			// "normal" atau lainnya -> no normalization
		}

		if box.MaxChars > 0 && len(userText) > box.MaxChars {
			userText = userText[:box.MaxChars]
		}

		fontSize := box.Size
		if fontSize == 0 {
			fontSize = 24
		}
		if err := dc.LoadFontFace(filepath.Clean(box.Font), fontSize); err != nil {
			return nil, fmt.Errorf("failed to load font %s: %w", box.Font, err)
		}

		if box.Color != "" {
			if c, err := hexToColor(box.Color); err == nil {
				dc.SetColor(c)
			} else {
				dc.SetRGB(0, 0, 0)
			}
		} else {
			dc.SetRGB(0, 0, 0)
		}

		padding := float64(box.Padding)
		if padding < 0 {
			padding = 0
		}

		effX := box.X + padding
		effY := box.Y + padding
		effW := float64(box.Width) - 2*padding
		effH := float64(box.Height) - 2*padding

		if effW <= 0 || effH <= 0 {
			continue
		}

		// text h-alignment
		align := gg.AlignLeft
		switch box.Align {
		case "center":
			align = gg.AlignCenter
		case "right":
			align = gg.AlignRight
		}

		centerX := effX + effW/2
		centerY := effY + effH/2

		dc.Push()
		dc.DrawRectangle(effX, effY, effW, effH)
		dc.Clip()

		dc.DrawStringWrapped(
			userText,
			centerX,
			centerY,
			0.5,             // anchor X
			0.5,             // anchor Y
			effW,            // wrapping width
			box.LineSpacing, // line spacing
			align,
		)

		dc.Pop()
	}

	return dc.Image(), nil
}

// helper: hex string (#RRGGBB atau #RRGGBBAA) → color.Color
func hexToColor(hex string) (color.Color, error) {
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b, a uint64
	var err error

	switch len(hex) {
	case 6: // RRGGBB
		r, err = strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return nil, err
		}
		g, err = strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return nil, err
		}
		b, err = strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return nil, err
		}
		a = 255
	case 8: // RRGGBBAA
		r, _ = strconv.ParseUint(hex[0:2], 16, 8)
		g, _ = strconv.ParseUint(hex[2:4], 16, 8)
		b, _ = strconv.ParseUint(hex[4:6], 16, 8)
		a, _ = strconv.ParseUint(hex[6:8], 16, 8)
	default:
		return nil, fmt.Errorf("invalid hex color: %s", hex)
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
}
