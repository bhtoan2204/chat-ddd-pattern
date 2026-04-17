package stripe

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
	"wechat-clone/core/modules/payment/domain/entity"
	"wechat-clone/core/modules/payment/providers"

	stripe "github.com/stripe/stripe-go/v75"
	stripeclient "github.com/stripe/stripe-go/v75/client"
	stripewebhook "github.com/stripe/stripe-go/v75/webhook"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestCreatePaymentUsesStripeGoCheckoutSession(t *testing.T) {
	var receivedForm url.Values
	var receivedVersion string

	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/checkout/sessions" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}

		receivedForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse request body: %v", err)
		}
		receivedVersion = r.Header.Get("Stripe-Version")

		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(`{"id":"cs_test_123","object":"checkout.session","url":"https://checkout.stripe.com/c/pay/cs_test_123"}`)),
		}, nil
	})}

	provider := &Provider{
		secretKey:  "sk_test_123",
		successURL: "https://merchant.example/success",
		cancelURL:  "https://merchant.example/cancel",
		httpClient: httpClient,
		apiBaseURL: "https://api.stripe.test",
	}
	provider.client = stripeclient.New(provider.secretKey, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
			HTTPClient: httpClient,
			URL:        stripe.String(provider.apiBaseURL),
		}),
		Connect: stripe.GetBackendWithConfig(stripe.ConnectBackend, &stripe.BackendConfig{
			HTTPClient: httpClient,
		}),
		Uploads: stripe.GetBackendWithConfig(stripe.UploadsBackend, &stripe.BackendConfig{
			HTTPClient: httpClient,
		}),
	})

	response, err := provider.CreatePayment(context.Background(), providers.CreatePaymentRequest{
		TransactionID: "tx_123",
		Amount:        5000,
		Currency:      "USD",
		Metadata: map[string]string{
			"product_name":           "Wallet top up",
			"customer_email":         "buyer@example.com",
			"destination_account":    "acct_dest_123",
			"on_behalf_of":           "acct_platform_123",
			"application_fee_amount": "321",
			"statement_descriptor":   "TOPUP",
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Provider != ProviderName {
		t.Fatalf("expected provider %s, got %s", ProviderName, response.Provider)
	}
	if response.TransactionID != "tx_123" {
		t.Fatalf("expected transaction_id tx_123, got %s", response.TransactionID)
	}
	if response.ExternalRef != "cs_test_123" {
		t.Fatalf("expected external_ref cs_test_123, got %s", response.ExternalRef)
	}
	if response.CheckoutURL != "https://checkout.stripe.com/c/pay/cs_test_123" {
		t.Fatalf("unexpected checkout url %s", response.CheckoutURL)
	}
	if response.Status != entity.PaymentStatusPending {
		t.Fatalf("expected pending status, got %s", response.Status)
	}

	if receivedVersion != apiVersion {
		t.Fatalf("expected Stripe-Version %s, got %s", apiVersion, receivedVersion)
	}
	if receivedForm.Get("mode") != "payment" {
		t.Fatalf("expected mode payment, got %s", receivedForm.Get("mode"))
	}
	if receivedForm.Get("client_reference_id") != "tx_123" {
		t.Fatalf("expected client_reference_id tx_123, got %s", receivedForm.Get("client_reference_id"))
	}
	if receivedForm.Get("success_url") != "https://merchant.example/success" {
		t.Fatalf("unexpected success_url %s", receivedForm.Get("success_url"))
	}
	if receivedForm.Get("cancel_url") != "https://merchant.example/cancel" {
		t.Fatalf("unexpected cancel_url %s", receivedForm.Get("cancel_url"))
	}
	if receivedForm.Get("customer_email") != "buyer@example.com" {
		t.Fatalf("unexpected customer_email %s", receivedForm.Get("customer_email"))
	}
	if receivedForm.Get("line_items[0][quantity]") != "1" {
		t.Fatalf("unexpected quantity %s", receivedForm.Get("line_items[0][quantity]"))
	}
	if receivedForm.Get("line_items[0][price_data][currency]") != "usd" {
		t.Fatalf("unexpected currency %s", receivedForm.Get("line_items[0][price_data][currency]"))
	}
	if receivedForm.Get("line_items[0][price_data][unit_amount]") != "5000" {
		t.Fatalf("unexpected unit amount %s", receivedForm.Get("line_items[0][price_data][unit_amount]"))
	}
	if receivedForm.Get("line_items[0][price_data][product_data][name]") != "Wallet top up" {
		t.Fatalf("unexpected product name %s", receivedForm.Get("line_items[0][price_data][product_data][name]"))
	}
	if receivedForm.Get("metadata[transaction_id]") != "tx_123" {
		t.Fatalf("unexpected session transaction metadata %s", receivedForm.Get("metadata[transaction_id]"))
	}
	if receivedForm.Get("payment_intent_data[metadata][transaction_id]") != "tx_123" {
		t.Fatalf("unexpected payment intent transaction metadata %s", receivedForm.Get("payment_intent_data[metadata][transaction_id]"))
	}
	if receivedForm.Get("payment_intent_data[transfer_data][destination]") != "acct_dest_123" {
		t.Fatalf("unexpected destination %s", receivedForm.Get("payment_intent_data[transfer_data][destination]"))
	}
	if receivedForm.Get("payment_intent_data[on_behalf_of]") != "acct_platform_123" {
		t.Fatalf("unexpected on_behalf_of %s", receivedForm.Get("payment_intent_data[on_behalf_of]"))
	}
	if receivedForm.Get("payment_intent_data[application_fee_amount]") != "321" {
		t.Fatalf("unexpected application_fee_amount %s", receivedForm.Get("payment_intent_data[application_fee_amount]"))
	}
	if receivedForm.Get("payment_intent_data[statement_descriptor_suffix]") != "TOPUP" {
		t.Fatalf("unexpected statement_descriptor_suffix %s", receivedForm.Get("payment_intent_data[statement_descriptor_suffix]"))
	}
}

func TestVerifyWebhookAndParseCheckoutSessionEvent(t *testing.T) {
	provider := &Provider{webhookSecret: "whsec_test"}

	payload := []byte(`{
		"id":"evt_test_123",
		"object":"event",
		"api_version":"2026-02-25.clover",
		"type":"checkout.session.completed",
		"data":{
			"object":{
				"id":"cs_test_123",
				"object":"checkout.session",
				"client_reference_id":"tx_123",
				"payment_status":"paid",
				"amount_total":5000,
				"currency":"usd"
			}
		}
	}`)

	now := time.Now().UTC()
	signature := hex.EncodeToString(stripewebhook.ComputeSignature(now, payload, provider.webhookSecret))
	header := fmt.Sprintf("t=%d,v1=%s", now.Unix(), signature)

	event, err := provider.VerifyWebhook(context.Background(), payload, header)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.EventID != "evt_test_123" {
		t.Fatalf("expected event id evt_test_123, got %s", event.EventID)
	}
	if event.EventType != "checkout.session.completed" {
		t.Fatalf("unexpected event type %s", event.EventType)
	}
	if event.Attributes["api_version"] != "2026-02-25.clover" {
		t.Fatalf("unexpected api version %s", event.Attributes["api_version"])
	}

	result, err := provider.ParseEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no parse error, got %v", err)
	}
	if result.TransactionID != "tx_123" {
		t.Fatalf("expected transaction_id tx_123, got %s", result.TransactionID)
	}
	if result.Status != entity.PaymentStatusSuccess {
		t.Fatalf("expected success status, got %s", result.Status)
	}
	if result.Amount != 5000 {
		t.Fatalf("expected amount 5000, got %d", result.Amount)
	}
	if result.Currency != "usd" {
		t.Fatalf("expected currency usd, got %s", result.Currency)
	}
	if result.ExternalRef != "cs_test_123" {
		t.Fatalf("expected external_ref cs_test_123, got %s", result.ExternalRef)
	}
}

func TestParseChargeSucceededEvent(t *testing.T) {
	provider := &Provider{}

	event := &providers.WebhookEvent{
		Provider:  ProviderName,
		EventID:   "evt_test_charge_123",
		EventType: "charge.succeeded",
		Attributes: map[string]string{
			"object": `{
				"id":"ch_test_123",
				"object":"charge",
				"status":"succeeded",
				"paid":true,
				"amount":5000,
				"currency":"usd",
				"metadata":{
					"transaction_id":"tx_123"
				}
			}`,
		},
	}

	result, err := provider.ParseEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no parse error, got %v", err)
	}
	if result.TransactionID != "tx_123" {
		t.Fatalf("expected transaction_id tx_123, got %s", result.TransactionID)
	}
	if result.Status != entity.PaymentStatusSuccess {
		t.Fatalf("expected success status, got %s", result.Status)
	}
	if result.Amount != 5000 {
		t.Fatalf("expected amount 5000, got %d", result.Amount)
	}
	if result.Currency != "usd" {
		t.Fatalf("expected currency usd, got %s", result.Currency)
	}
	if result.ExternalRef != "ch_test_123" {
		t.Fatalf("expected external_ref ch_test_123, got %s", result.ExternalRef)
	}
}

func TestParseChargeRefundedEvent(t *testing.T) {
	provider := &Provider{}

	event := &providers.WebhookEvent{
		Provider:  ProviderName,
		EventID:   "evt_test_refund_123",
		EventType: "charge.refunded",
		Attributes: map[string]string{
			"object": `{
				"id":"ch_test_123",
				"object":"charge",
				"status":"succeeded",
				"paid":true,
				"refunded":true,
				"amount":5000,
				"amount_refunded":5000,
				"currency":"usd",
				"metadata":{
					"transaction_id":"tx_123"
				}
			}`,
		},
	}

	result, err := provider.ParseEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no parse error, got %v", err)
	}
	if result.TransactionID != "tx_123" {
		t.Fatalf("expected transaction_id tx_123, got %s", result.TransactionID)
	}
	if result.Status != entity.PaymentStatusRefunded {
		t.Fatalf("expected refunded status, got %s", result.Status)
	}
	if result.Amount != 5000 {
		t.Fatalf("expected amount 5000, got %d", result.Amount)
	}
	if result.ExternalRef != "ch_test_123" {
		t.Fatalf("expected external_ref ch_test_123, got %s", result.ExternalRef)
	}
}

func TestParseChargeDisputeFundsWithdrawnEvent(t *testing.T) {
	provider := &Provider{}

	event := &providers.WebhookEvent{
		Provider:  ProviderName,
		EventID:   "evt_test_dispute_123",
		EventType: "charge.dispute.funds_withdrawn",
		Attributes: map[string]string{
			"object": `{
				"id":"dp_test_123",
				"object":"dispute",
				"amount":5000,
				"currency":"usd",
				"metadata":{
					"transaction_id":"tx_123"
				},
				"charge":{
					"id":"ch_test_123"
				}
			}`,
		},
	}

	result, err := provider.ParseEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no parse error, got %v", err)
	}
	if result.TransactionID != "tx_123" {
		t.Fatalf("expected transaction_id tx_123, got %s", result.TransactionID)
	}
	if result.Status != entity.PaymentStatusChargeback {
		t.Fatalf("expected chargeback status, got %s", result.Status)
	}
	if result.Amount != 5000 {
		t.Fatalf("expected amount 5000, got %d", result.Amount)
	}
	if result.ExternalRef != "ch_test_123" {
		t.Fatalf("expected external_ref ch_test_123, got %s", result.ExternalRef)
	}
}

func TestVerifyWebhookRejectsInvalidSignature(t *testing.T) {
	provider := &Provider{webhookSecret: "whsec_test"}

	_, err := provider.VerifyWebhook(context.Background(), []byte(`{"id":"evt_test_123"}`), "t=1,v1=bad")
	if !errors.Is(err, providers.ErrInvalidWebhookSignature) {
		t.Fatalf("expected invalid webhook signature error, got %v", err)
	}
}

func TestParseUnsupportedStripeEventIgnored(t *testing.T) {
	provider := &Provider{}

	_, err := provider.ParseEvent(context.Background(), &providers.WebhookEvent{
		Provider:  ProviderName,
		EventID:   "evt_test_ignore_123",
		EventType: "payment_intent.created",
		Attributes: map[string]string{
			"object": `{}`,
		},
	})
	if !errors.Is(err, providers.ErrWebhookEventIgnored) {
		t.Fatalf("expected ignored webhook event error, got %v", err)
	}
}
