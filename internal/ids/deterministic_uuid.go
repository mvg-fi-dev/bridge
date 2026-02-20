package ids

import "github.com/google/uuid"

// DeterministicUUID returns a stable UUID derived from input string.
// Useful for trace_id idempotency when the API requires UUID format.
func DeterministicUUID(name string) string {
	// NameSpaceOID is stable; SHA1 is fine for namespacing trace IDs.
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(name)).String()
}
