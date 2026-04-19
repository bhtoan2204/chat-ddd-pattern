package tokendigest

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type hmacSHA256Digester struct {
	secret []byte
}

func NewHMACSHA256Digester(secret string) (Digester, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, stackErr.Error(fmt.Errorf("token digest secret is empty"))
	}

	return &hmacSHA256Digester{
		secret: []byte(secret),
	}, nil
}

func (d *hmacSHA256Digester) Digest(ctx context.Context, value string) (string, error) {
	mac := hmac.New(sha256.New, d.secret)
	if _, err := mac.Write([]byte(value)); err != nil {
		return "", stackErr.Error(fmt.Errorf("compute token digest failed: %w", err))
	}
	return base64.RawStdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (d *hmacSHA256Digester) Verify(ctx context.Context, value string, digest string) (bool, error) {
	computed, err := d.Digest(ctx, value)
	if err != nil {
		return false, stackErr.Error(err)
	}

	return hmac.Equal([]byte(computed), []byte(strings.TrimSpace(digest))), nil
}
