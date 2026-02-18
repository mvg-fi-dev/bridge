package ids

import (
	"strings"

	"github.com/google/uuid"
)

func NewPublicID(prefix string) string {
	// Short, URL-safe-ish ID. Not cryptographic; fine for public order ids.
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	return prefix + "_" + id[:12]
}

func NewUUID() string {
	return uuid.New().String()
}
