package domain

type Meme struct {
	ImageUrl    string `json:"image_url"`
	ContentType string `json:"content_type"`
	Size        string `json:"size"`
}
