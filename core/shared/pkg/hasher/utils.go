package hasher

import (
	"crypto/rand"
	"fmt"
	stackerr "go-socket/core/shared/pkg/stackErr"
)

func genSalt(keyLen uint32) ([]byte, error) {
	salt := make([]byte, keyLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, stackerr.Error(fmt.Errorf("failed to generate salt: %w", err))
	}
	return salt, nil
}
