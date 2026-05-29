package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestCheckAddress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/address/check" {
			t.Fatalf("path = %q, want /v1/address/check", r.URL.Path)
		}
		var req models.AddressCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if got, want := req.Address, "1 Market St"; got != want {
			t.Fatalf("Address = %q, want %q", got, want)
		}
		if err := json.NewEncoder(w).Encode(models.AddressCheck{
			Address:          req.Address,
			FormattedAddress: "1 Market St, San Francisco, CA",
			Status:           models.AddressCheckStatusEligible,
			CheckoutEligible: true,
		}); err != nil {
			t.Fatalf("write address response: %v", err)
		}
	}))
	defer server.Close()

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	check, err := client.CheckAddress(context.Background(), "1 Market St")
	if err != nil {
		t.Fatalf("CheckAddress: %v", err)
	}
	if !check.CheckoutEligible {
		t.Fatal("CheckoutEligible = false, want true")
	}
}
