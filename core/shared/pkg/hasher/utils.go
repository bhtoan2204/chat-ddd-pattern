package hasher

import (
	"crypto/rand"
	"fmt"
)

func genSalt(keyLen uint32) ([]byte, error) {
	salt := make([]byte, keyLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}
