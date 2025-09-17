package storage

import (
	"MemeCraft/internal/port"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type CatboxMoeProvider struct {
	ApiUrl string
	Client *resty.Client
}

func (s *CatboxMoeProvider) Upload(ctx context.Context, data io.Reader) (*port.UploadResult, error) {
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	return s.UploadBytes(ctx, dataBytes)
}

func (s *CatboxMoeProvider) UploadBytes(ctx context.Context, data []byte) (*port.UploadResult, error) {
	contentType := http.DetectContentType(data)
	filename := GenerateFileName(contentType)

	resp, err := s.Client.R().
		SetContext(ctx).
		SetFileReader("fileToUpload", filename, bytes.NewReader(data)).
		SetFormData(map[string]string{
			"reqtype": "fileupload",
		}).
		Post(s.ApiUrl)

	if err != nil {
		log.Printf("upload file error: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status code %d", resp.StatusCode())
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

func (s *CatboxMoeProvider) GetStorageName() string {
	return "catbox.moe"
}

func NewCatboxMoeStorage() *CatboxMoeProvider {
	return &CatboxMoeProvider{
		ApiUrl: "https://catbox.moe/user/api.php",
		Client: resty.New(),
	}
}
