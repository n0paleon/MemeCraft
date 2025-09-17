package meme

import (
	"MemeCraft/internal/domain"
	"MemeCraft/internal/port"
	"MemeCraft/internal/preset"
	"MemeCraft/internal/service/imageutil"
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Generator struct {
	registry       *preset.Registry
	storageAdapter port.StorageProvider
}

func (g *Generator) Generate(config *Config) (*domain.Meme, error) {
	p, ok := g.registry.Get(config.PresetId)
	if !ok {
		return nil, errors.New("preset not found")
	}

	if len(config.Overlay) == 0 {
		return nil, errors.New("no overlays specified")
	}

	overlay, _, err := image.Decode(bytes.NewReader(config.Overlay))
	if err != nil {
		log.Errorf("Error decoding overlay: %v", err)
		log.Infof("overlay length: %d", len(config.Overlay))
		return nil, err
	}

	resizeMode := p.ResizeMode
	if config.ResizeMode != "" {
		resizeMode = imageutil.ResizeMode(config.ResizeMode)
	}
	overlay, err = imageutil.ResizeWithMode(overlay, resizeMode, p.Overlay.Width, p.Overlay.Height)
	if err != nil {
		log.Errorf("Error resizing overlay: %v", err)
		return nil, errors.New("failed to resize overlay")
	}

	overlayedImage := imageutil.Overlay(p.BaseImageDecoded, overlay, p.Overlay.X, p.Overlay.Y, "back", p.Overlay.Rotate)

	textbox := convertTextBoxPreset(p.TextBoxes)
	drawnTextImage := overlayedImage
	for k, v := range config.Text {
		for _, tb := range textbox {
			if tb.Name == k {
				var box []imageutil.TextBox
				box = append(box, tb)
				text := map[string]string{k: v}
				drawnTextImage, err = imageutil.DrawTextBoxes(drawnTextImage, text, box)
				if err != nil {
					log.Errorf("Error while drawing text: %v", err)
					return nil, errors.New("failed to draw text")
				}
			}
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, drawnTextImage, nil); err != nil {
		log.Errorf("failed to encode image: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uploadResult, err := g.storageAdapter.UploadBytes(ctx, buf.Bytes())
	if err != nil {
		log.Errorf("failed to upload image: %v", err)
		return nil, errors.New("failed to upload image")
	}

	return &domain.Meme{
		ImageUrl:    uploadResult.DirectURL,
		ContentType: uploadResult.ContentType,
		Size:        uploadResult.BytesReadable,
	}, nil
}

func convertTextBoxPreset(sets []preset.TextBox) []imageutil.TextBox {
	textboxes := make([]imageutil.TextBox, len(sets))

	for i, p := range sets {
		textboxes[i] = imageutil.TextBox{
			Name:        p.Name,
			X:           p.X,
			Y:           p.Y,
			Width:       p.Width,
			Height:      p.Height,
			Font:        p.Font,
			Size:        p.Size,
			LineSpacing: p.LineSpacing,
			Color:       p.Color,
			Align:       p.Align,
			Padding:     p.Padding,
			MaxChars:    p.MaxChars,
			Normalize:   p.Normalize,
		}
	}

	return textboxes
}

func (g *Generator) GetPresetById(presetId string) (*preset.PresetSummary, error) {
	p, ok := g.registry.GetSummaryById(presetId)
	if !ok {
		return nil, errors.New("preset not found")
	}
	return p, nil
}

func (g *Generator) GetAllPreset() []*preset.PresetSummary {
	return g.registry.GetAll()
}

func NewGenerator(reg *preset.Registry, storageAdapter port.StorageProvider) *Generator {
	log.Infof("Storage adapter => %s", storageAdapter.GetStorageName())
	return &Generator{
		registry:       reg,
		storageAdapter: storageAdapter,
	}
}
