package hasher

import (
	"context"
	"testing"
)

func TestHasher(t *testing.T) {
	t.Run("HashAndVerifyWithSameInstance", func(t *testing.T) {
		hasher, err := newHasher()
		if err != nil {
			t.Fatalf("failed to create hasher: %v", err)
		}

		hash, err := hasher.Hash(context.Background(), "password")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		t.Logf("hash: %s", hash)

		ok, err := hasher.Verify(context.Background(), "password", hash)
		if err != nil {
			t.Fatalf("failed to verify password: %v", err)
		}
		t.Logf("ok: %v", ok)

		if !ok {
			t.Fatalf("expected password to be valid")
		}
	})
}

func TestHasherWhenRebootsShouldReturnSameHash(t *testing.T) {
	t.Run("VerifyHashAfterRecreatingHasher", func(t *testing.T) {
		// First instance creates the hash
		hasher1, err := newHasher()
		if err != nil {
			t.Fatalf("failed to create hasher: %v", err)
		}

		hash, err := hasher1.Hash(context.Background(), "password")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		t.Logf("hash: %s", hash)

		// Second instance verifies the same hash (simulate reboot)
		hasher2, err := newHasher()
		if err != nil {
			t.Fatalf("failed to create hasher: %v", err)
		}

		ok, err := hasher2.Verify(context.Background(), "password", hash)
		if err != nil {
			t.Fatalf("failed to verify password: %v", err)
		}
		t.Logf("ok: %v", ok)

		if !ok {
			t.Fatalf("expected password to be valid after reboot")
		}
	})
}
