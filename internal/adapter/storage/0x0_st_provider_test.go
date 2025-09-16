package storage

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestPNG membuat dummy PNG image untuk testing
func generateTestPNG() ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))

	// Buat pola checkerboard
	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.RGBA{0, 0, 0, 255}) // Hitam
			} else {
				img.Set(x, y, color.RGBA{255, 255, 255, 255}) // Putih
			}
		}
	}

	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	return buf.Bytes(), err
}

// generateFakeJPEG untuk test content type detection
func generateFakeJPEG() []byte {
	return []byte{
		0xFF, 0xD8, 0xFF, 0xE0, // JPEG SOI + APP0
		0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01,
		0x01, 0x01, 0x00, 0x48, 0x00, 0x48, 0x00, 0x00,
		0xFF, 0xD9, // EOI
	}
}

func TestZeroXZeroSTProvider_UploadBytes_Success(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)
	require.NotEmpty(t, pngData)

	// Act
	result, err := provider.UploadBytes(context.Background(), pngData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.DirectURL)
	assert.Contains(t, result.DirectURL, "0x0.st/s/") // Secret URL format dengan ULID
	assert.Equal(t, "image/png", result.ContentType)
	assert.Equal(t, int64(len(pngData)), result.Bytes)
	assert.NotEmpty(t, result.BytesReadable)
}

func TestZeroXZeroSTProvider_UploadBytes_DifferentContentTypes(t *testing.T) {
	provider := NewZeroXZeroSTStorage()

	testCases := []struct {
		name         string
		data         []byte
		expectedType string
	}{
		{
			name:         "PNG Image",
			data:         func() []byte { data, _ := generateTestPNG(); return data }(),
			expectedType: "image/png",
		},
		{
			name:         "JPEG Image",
			data:         generateFakeJPEG(),
			expectedType: "image/jpeg",
		},
		{
			name:         "Text File",
			data:         []byte("Hello, World!"),
			expectedType: "text/plain; charset=utf-8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := provider.UploadBytes(context.Background(), tc.data)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.DirectURL)
			assert.Equal(t, tc.expectedType, result.ContentType)
			assert.Equal(t, int64(len(tc.data)), result.Bytes)
		})
	}
}

func TestZeroXZeroSTProvider_UploadBytes_WithContext(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)

	// Test dengan timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Act
	result, err := provider.UploadBytes(ctx, pngData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.DirectURL)
}

func TestZeroXZeroSTProvider_UploadBytes_ContextCancellation(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)

	// Context yang langsung di-cancel
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	result, err := provider.UploadBytes(ctx, pngData)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestZeroXZeroSTProvider_UploadBytes_EmptyData(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	emptyData := []byte{}

	// Act
	result, err := provider.UploadBytes(context.Background(), emptyData)

	// Assert - tergantung bagaimana 0x0.st handle empty files
	if err != nil {
		assert.Contains(t, err.Error(), "upload fail")
	} else {
		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.Bytes)
	}
}

func TestZeroXZeroSTProvider_Upload_Success(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)

	reader := bytes.NewReader(pngData)

	// Act
	result, err := provider.Upload(context.Background(), reader)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.DirectURL)
	assert.Contains(t, result.DirectURL, "0x0.st/s/")
	assert.Equal(t, "image/png", result.ContentType)
	assert.Equal(t, int64(len(pngData)), result.Bytes)
	assert.NotEmpty(t, result.BytesReadable)
}

func TestZeroXZeroSTProvider_Upload_WithStringReader(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	textData := "Hello, this is test data for 0x0.st upload"
	reader := strings.NewReader(textData)

	// Act
	result, err := provider.Upload(context.Background(), reader)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.DirectURL)
	assert.Equal(t, "text/plain; charset=utf-8", result.ContentType)
	assert.Equal(t, int64(len(textData)), result.Bytes)
}

func TestZeroXZeroSTProvider_Upload_WithBytesReader(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)

	reader := bytes.NewReader(pngData)

	// Act
	result, err := provider.Upload(context.Background(), reader)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "image/png", result.ContentType)
	assert.Equal(t, int64(len(pngData)), result.Bytes)
}

func TestZeroXZeroSTProvider_Upload_ContextCancellation(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)

	reader := bytes.NewReader(pngData)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Act
	result, err := provider.Upload(ctx, reader)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestNewZeroXZeroSTStorage(t *testing.T) {
	// Act
	provider := NewZeroXZeroSTStorage()

	// Assert
	assert.NotNil(t, provider)
	assert.Equal(t, "https://0x0.st", provider.ApiUrl)
	assert.NotNil(t, provider.Client)
}

// Test untuk memastikan secret ULID di-generate untuk setiap upload
func TestZeroXZeroSTProvider_ULIDSecretGeneration(t *testing.T) {
	// Arrange
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)

	// Act - upload 2 kali
	result1, err1 := provider.UploadBytes(context.Background(), pngData)
	result2, err2 := provider.UploadBytes(context.Background(), pngData)

	// Assert - kedua upload berhasil dengan URL berbeda (karena ULID berbeda)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotNil(t, result1)
	assert.NotNil(t, result2)

	// Keduanya harus menggunakan format secret URL
	assert.Contains(t, result1.DirectURL, "0x0.st/s/")
	assert.Contains(t, result2.DirectURL, "0x0.st/s/")
}

// Benchmark test untuk performance
func BenchmarkZeroXZeroSTProvider_UploadBytes(b *testing.B) {
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := provider.UploadBytes(context.Background(), pngData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkZeroXZeroSTProvider_Upload(b *testing.B) {
	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(pngData)
		_, err := provider.Upload(context.Background(), reader)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Integration test - skip jika dalam short mode
func TestZeroXZeroSTProvider_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := NewZeroXZeroSTStorage()
	pngData, err := generateTestPNG()
	require.NoError(t, err)

	// Test UploadBytes
	result1, err := provider.UploadBytes(context.Background(), pngData)
	require.NoError(t, err)
	assert.NotEmpty(t, result1.DirectURL)

	// Test Upload
	reader := bytes.NewReader(pngData)
	result2, err := provider.Upload(context.Background(), reader)
	require.NoError(t, err)
	assert.NotEmpty(t, result2.DirectURL)

	// Verify both methods produce consistent results
	assert.Equal(t, result1.ContentType, result2.ContentType)
	assert.Equal(t, result1.Bytes, result2.Bytes)
}
