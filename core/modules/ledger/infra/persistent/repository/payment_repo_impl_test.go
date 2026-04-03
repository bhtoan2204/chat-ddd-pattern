package repository

import "testing"

func TestToStoredExternalRefUsesProvisionalValueWhenEmpty(t *testing.T) {
	first := toStoredExternalRef("stripe", "tx-001", "")
	second := toStoredExternalRef("stripe", "tx-002", "")

	if first == nil || *first == "" {
		t.Fatalf("expected provisional external ref for first transaction")
	}
	if second == nil || *second == "" {
		t.Fatalf("expected provisional external ref for second transaction")
	}
	if *first == *second {
		t.Fatalf("expected different provisional external refs, got %q", *first)
	}
}

func TestFromStoredExternalRefHidesProvisionalValue(t *testing.T) {
	stored := toStoredExternalRef("stripe", "tx-001", "")
	got := fromStoredExternalRef("stripe", "tx-001", stored)

	if got != "" {
		t.Fatalf("expected provisional external ref to be hidden, got %q", got)
	}
}

func TestFromStoredExternalRefPreservesRealValue(t *testing.T) {
	stored := toStoredExternalRef("stripe", "tx-001", "cs_test_123")
	got := fromStoredExternalRef("stripe", "tx-001", stored)

	if got != "cs_test_123" {
		t.Fatalf("expected real external ref to round-trip, got %q", got)
	}
}
