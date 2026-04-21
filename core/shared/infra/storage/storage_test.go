package storage

import "testing"

func TestParsePublicBaseURL(t *testing.T) {
	t.Parallel()

	u, err := parsePublicBaseURL("https://cdn.example.com")
	if err != nil {
		t.Fatalf("parsePublicBaseURL() error = %v", err)
	}
	if u == nil || u.Scheme != "https" || u.Host != "cdn.example.com" {
		t.Fatalf("unexpected url: %+v", u)
	}
}

func TestParsePublicBaseURLEmpty(t *testing.T) {
	t.Parallel()

	u, err := parsePublicBaseURL("")
	if err != nil {
		t.Fatalf("parsePublicBaseURL() error = %v", err)
	}
	if u != nil {
		t.Fatalf("expected nil url, got %+v", u)
	}
}

func TestParsePublicBaseURLKeepsHostPort(t *testing.T) {
	t.Parallel()

	publicBaseURL, err := parsePublicBaseURL("http://20.2.92.187:9001")
	if err != nil {
		t.Fatalf("parsePublicBaseURL() error = %v", err)
	}
	if publicBaseURL.Host != "20.2.92.187:9001" {
		t.Fatalf("expected host with port, got %s", publicBaseURL.Host)
	}
}
