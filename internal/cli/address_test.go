package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestCheckAddressCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/address/check" {
			t.Fatalf("path = %q, want /v1/address/check", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(models.AddressCheck{
			Address:          "1 Market St",
			FormattedAddress: "1 Market St, San Francisco, CA",
			Status:           models.AddressCheckStatusEligible,
			CheckoutEligible: true,
		}); err != nil {
			t.Fatalf("write address response: %v", err)
		}
	}))
	defer server.Close()

	out, err := executeForTestWithAPI(t, "test-version", server.URL, "check-address", "1 Market St")
	if err != nil {
		t.Fatalf("check-address command failed: %v", err)
	}
	for _, want := range []string{
		"eligible: address is inside the default delivery zone",
		"Checkout authorizes payment only",
		"Address resolved to: 1 Market St, San Francisco, CA",
		"Is this correct?",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output does not contain %q:\n%s", want, out)
		}
	}
}

func TestCheckAddressCommandJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(models.AddressCheck{
			Address:          "1 Market St",
			FormattedAddress: "1 Market St, San Francisco, CA",
			Status:           models.AddressCheckStatusEligible,
			CheckoutEligible: true,
		}); err != nil {
			t.Fatalf("write address response: %v", err)
		}
	}))
	defer server.Close()

	out, err := executeForTestWithAPI(t, "test-version", server.URL, "check-address", "--json", "1 Market St")
	if err != nil {
		t.Fatalf("check-address command failed: %v", err)
	}

	var check models.AddressCheck
	if err := json.NewDecoder(bytes.NewBufferString(out)).Decode(&check); err != nil {
		t.Fatalf("decode json output: %v\n%s", err, out)
	}
	if got, want := check.Status, models.AddressCheckStatusEligible; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func TestCheckAddressCommandReturnsActionableError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(models.AddressCheck{
			Address:          "1 Market",
			Status:           models.AddressCheckStatusAddressUnresolved,
			CheckoutEligible: false,
		}); err != nil {
			t.Fatalf("write address response: %v", err)
		}
	}))
	defer server.Close()

	out, err := executeForTestWithAPI(t, "test-version", server.URL, "check-address", "1 Market")
	if err == nil {
		t.Fatal("check-address command succeeded with unresolved address")
	}
	for _, want := range []string{
		"address_unresolved",
		"complete and specific delivery address",
	} {
		if !strings.Contains(out, want) && !strings.Contains(err.Error(), want) {
			t.Fatalf("output/error does not contain %q:\nout=%s\nerr=%v", want, out, err)
		}
	}
}
