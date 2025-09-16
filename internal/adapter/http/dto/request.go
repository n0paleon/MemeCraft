package dto

type CreateMemeRequest struct {
	Overlay    string            `json:"overlay"`
	ResizeMode string            `json:"resize_mode"`
	Text       map[string]string `json:"text"`
}
