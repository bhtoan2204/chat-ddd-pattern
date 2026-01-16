package hasher

import "context"

type Hasher interface {
	Hash(ctx context.Context, value string) (string, error)
	Verify(ctx context.Context, val string, hash string) (bool, error)
}
