package repository

import "testing"

func TestStorageExternalRefPlaceholderRoundTrip(t *testing.T) {
	stored := toStorageExternalRef(" stripe ", " txn-1 ", "")
	if stored == nil {
		t.Fatalf("expected placeholder external ref")
	}
	if *stored == "" {
		t.Fatalf("expected non-empty placeholder")
	}

	got := fromStorageExternalRef(*stored)
	if got != "" {
		t.Fatalf("expected placeholder to map back to empty external ref, got %q", got)
	}
}

func TestStorageExternalRefKeepsRealValue(t *testing.T) {
	stored := toStorageExternalRef("stripe", "txn-1", " ext-1 ")
	if stored == nil {
		t.Fatalf("expected stored external ref")
	}
	if *stored != "ext-1" {
		t.Fatalf("unexpected stored value: %q", *stored)
	}

	got := fromStorageExternalRef(*stored)
	if got != "ext-1" {
		t.Fatalf("unexpected restored value: %q", got)
	}
}
