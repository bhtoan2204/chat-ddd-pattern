package handler

import (
	"net/http"
	"testing"
)

func TestWebhookSignatureUsesStripeHeaderForStripeProvider(t *testing.T) {
	header := http.Header{}
	header.Set("Stripe-Signature", "stripe-signature")
	header.Set("X-Signature", "legacy-signature")

	got := webhookSignature("stripe", header)
	if got != "stripe-signature" {
		t.Fatalf("expected stripe signature header, got %q", got)
	}
}

func TestWebhookSignatureUsesLegacyHeaderForNonStripeProvider(t *testing.T) {
	header := http.Header{}
	header.Set("Stripe-Signature", "stripe-signature")
	header.Set("X-Signature", "legacy-signature")

	got := webhookSignature("mock", header)
	if got != "legacy-signature" {
		t.Fatalf("expected legacy signature header, got %q", got)
	}
}

func TestWebhookSignatureFallsBackToStripeHeader(t *testing.T) {
	header := http.Header{}
	header.Set("Stripe-Signature", "stripe-signature")

	got := webhookSignature("mock", header)
	if got != "stripe-signature" {
		t.Fatalf("expected stripe signature fallback, got %q", got)
	}
}
