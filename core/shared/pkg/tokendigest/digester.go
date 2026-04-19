package tokendigest

import "context"

//go:generate mockgen -package=tokendigest -destination=digester_mock.go -source=digester.go
type Digester interface {
	Digest(ctx context.Context, value string) (string, error)
	Verify(ctx context.Context, value string, digest string) (bool, error)
}
