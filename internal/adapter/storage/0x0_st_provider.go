package storage

import (
	"MemeCraft/internal/port"
	"MemeCraft/pkg/ulidgen"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ZeroXZeroSTProvider struct {
	ApiUrl string
	Client *resty.Client
}

func (s *ZeroXZeroSTProvider) Upload(ctx context.Context, data io.Reader) (*port.UploadResult, error) {
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	return s.UploadBytes(ctx, dataBytes)
}

func (s *ZeroXZeroSTProvider) UploadBytes(ctx context.Context, data []byte) (*port.UploadResult, error) {
	contentType := http.DetectContentType(data)
	filename := GenerateFileName(contentType)

	resp, err := s.Client.R().
		SetContext(ctx).
		SetFileReader("file", filename, bytes.NewReader(data)).
		SetFormData(map[string]string{
			"secret":  ulidgen.GenerateULID(),
			"expires": "24", // 24 hour
		}).
		Post(s.ApiUrl)
	if err != nil {
		log.Printf("upload file error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("upload fail, status code: %v", resp.StatusCode())
	}

	directURL := strings.TrimSpace(string(resp.Body()))
	contentLength := int64(len(data))

	return &port.UploadResult{
		DirectURL:     directURL,
		ContentType:   contentType,
		Bytes:         contentLength,
		BytesReadable: ByteCountSI(contentLength),
	}, nil
}

func (s *ZeroXZeroSTProvider) GetStorageName() string {
	return "0x0.st"
}

func NewZeroXZeroSTStorage() *ZeroXZeroSTProvider {
	return &ZeroXZeroSTProvider{
		ApiUrl: "https://0x0.st",
		Client: resty.New(),
	}
}
