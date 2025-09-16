package ulidgen

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID generates a new ULID string
func GenerateULID() string {
	source := rand.NewSource(time.Now().UnixNano())
	entropy := rand.New(source)

	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}

// GenerateULIDWithTime generates ULID with specific time
func GenerateULIDWithTime(t time.Time) string {
	source := rand.NewSource(time.Now().UnixNano())
	entropy := rand.New(source)

	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

// GenerateULIDSafe generates ULID with error handling
func GenerateULIDSafe() (string, error) {
	source := rand.NewSource(time.Now().UnixNano())
	entropy := rand.New(source)

	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
