package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const MaxImageSize = 2 * 1024 * 1024 // 2MB

var allowedTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/jpg":  true,
}

func DownloadImageAsBytes(url string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.RawBody().Close()

	contentType := resp.Header().Get("Content-Type")
	if !allowedTypes[contentType] {
		return nil, errors.New("unsupported content type: " + contentType)
	}

	var buf bytes.Buffer
	limitedReader := io.LimitReader(resp.RawBody(), MaxImageSize+1)

	n, err := io.Copy(&buf, limitedReader)
	if err != nil {
		return nil, err
	}
	if n > MaxImageSize {
		return nil, errors.New("file too large (actual body)")
	}

	detected := http.DetectContentType(buf.Bytes())
	if !allowedTypes[detected] {
		return nil, errors.New("unsupported content type (detected): " + detected)
	}

	return buf.Bytes(), nil
}
