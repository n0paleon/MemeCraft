package preset

import (
	"encoding/json"
	"image"
	"log"
	"os"
	"path/filepath"
)

type Registry struct {
	presets map[string]*Preset
	cache   []*PresetSummary
}

type TextBoxSummary struct {
	Name     string `json:"name"`
	MaxChars int    `json:"max_chars"`
}

type OverlaySummary struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type PresetSummary struct {
	Name         string           `json:"name"`
	ID           string           `json:"id"`
	ExampleImage string           `json:"example_image"`
	Overlay      OverlaySummary   `json:"overlay"`
	TextBoxes    []TextBoxSummary `json:"text_boxes"`
}

func NewRegistry() *Registry {
	return &Registry{
		presets: make(map[string]*Preset),
	}
}

func (r *Registry) LoadFromDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		path := filepath.Join(dir, f.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var p Preset
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}

		// BaseImage dianggap relatif terhadap workdir
		imgPath := filepath.Clean(p.BaseImage)
		log.Println("loading base image =>", imgPath)

		imgFile, err := os.Open(imgPath)
		if err != nil {
			return err
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			return err
		}
		p.BaseImageDecoded = img

		for i := range p.TextBoxes {
			if p.TextBoxes[i].Font != "" {
				p.TextBoxes[i].Font = filepath.Clean(p.TextBoxes[i].Font)
				log.Println("registered font for box", p.TextBoxes[i].Name, "=>", p.TextBoxes[i].Font)
			}
		}

		r.presets[p.ID] = &p
	}

	for p := range r.presets {
		log.Println("loaded preset =>", p)
	}
	return nil
}

func (r *Registry) Get(presetId string) (*Preset, bool) {
	p, ok := r.presets[presetId]
	return p, ok
}

func (r *Registry) GetSummaryById(presetId string) (*PresetSummary, bool) {
	p, ok := r.presets[presetId]
	if !ok {
		return nil, false
	}

	textbox := make([]TextBoxSummary, len(p.TextBoxes))
	for i, tb := range p.TextBoxes {
		textbox[i] = TextBoxSummary{
			Name:     tb.Name,
			MaxChars: tb.MaxChars,
		}
	}

	return &PresetSummary{
		Name:         p.Name,
		ID:           p.ID,
		ExampleImage: p.ExampleImage,
		Overlay: OverlaySummary{
			Width:  p.Overlay.Width,
			Height: p.Overlay.Height,
		},
		TextBoxes: textbox,
	}, true
}

func (r *Registry) GetAll() []*PresetSummary {
	if r.cache != nil {
		return r.cache
	}

	summaries := make([]*PresetSummary, 0, len(r.presets))
	for _, p := range r.presets {
		overlay := OverlaySummary{
			Width:  p.Overlay.Width,
			Height: p.Overlay.Height,
		}

		tbSummaries := make([]TextBoxSummary, 0, len(p.TextBoxes))
		for _, tb := range p.TextBoxes {
			tbSummaries = append(tbSummaries, TextBoxSummary{
				Name:     tb.Name,
				MaxChars: tb.MaxChars,
			})
		}

		summary := &PresetSummary{
			Name:         p.Name,
			ID:           p.ID,
			ExampleImage: p.ExampleImage,
			Overlay:      overlay,
			TextBoxes:    tbSummaries,
		}

		summaries = append(summaries, summary)
	}

	r.cache = summaries

	return summaries
}
