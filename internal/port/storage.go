package port

import (
	"context"
	"io"
)

type UploadResult struct {
	DirectURL     string `json:"direct_url"`
	ContentType   string `json:"content_type"`
	Bytes         int64  `json:"size"`
	BytesReadable string `json:"bytes_readable"`
}

type StorageProvider interface {
	Upload(ctx context.Context, data io.Reader) (*UploadResult, error)
	UploadBytes(ctx context.Context, data []byte) (*UploadResult, error)
}
