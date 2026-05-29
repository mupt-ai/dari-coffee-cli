package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestOrderCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders" {
			t.Fatalf("path = %q, want /v1/orders", r.URL.Path)
		}
		var req models.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if got, want := req.Address, "1 Market St, San Francisco, CA 94105"; got != want {
			t.Fatalf("address = %q, want %q", got, want)
		}
		if got, want := len(req.Items), 2; got != want {
			t.Fatalf("item count = %d, want %d", got, want)
		}
		if got, want := req.Items[0].Quantity, 2; got != want {
			t.Fatalf("quantity = %d, want %d", got, want)
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(models.CreateOrderResponse{
			Order: models.OrderSummary{
				ID:            "drk_test",
				Status:        "pending_payment",
				ZoneStatus:    "inside",
				SubtotalCents: 1725,
				TotalCents:    1725,
				Currency:      "usd",
			},
			AddressCheck: models.AddressCheck{
				FormattedAddress: "1 Market St, San Francisco, CA 94105",
			},
			CheckoutSessionID: "cs_test_order",
			CheckoutURL:       "https://checkout.stripe.com/c/pay/cs_test_order",
		}); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	out, err := executeForTestWithAPI(
		t,
		"test-version",
		server.URL,
		"order",
		"--name", "Ada Lovelace",
		"--email", "ada@example.com",
		"--company", "Analytical Engines Inc.",
		"--address", "1 Market St, San Francisco, CA 94105",
		"--items-json", `[
		  {"shop_slug":"starbucks","drink_slug":"iced-caramel-macchiato","size":"grande","quantity":2,"modifications":"extra ice"},
		  {"shop_slug":"starbucks","drink_slug":"americano","size":"grande","quantity":1}
		]`,
	)
	if err != nil {
		t.Fatalf("order command failed: %v", err)
	}
	for _, want := range []string{
		"Order created: drk_test",
		"Status: pending_payment",
		"Total: $17.25 usd",
		"Checkout: https://checkout.stripe.com/c/pay/cs_test_order",
		"Authorize payment in Stripe. Dari captures only after accepting the order.",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output does not contain %q:\n%s", want, out)
		}
	}
}

func TestOrderCommandRequiresFlags(t *testing.T) {
	out, err := executeForTest("test-version", "order")
	if err == nil {
		t.Fatal("order command succeeded")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Fatalf("error = %v, want required flag error; out=%s", err, out)
	}
}

func TestOrderCommandUsesStableErrorCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(models.ErrorResponse{
			Code:  models.ErrorCodeAddressUnresolved,
			Error: "address_unresolved",
		}); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	_, err := executeForTestWithAPI(
		t,
		"test-version",
		server.URL,
		"order",
		"--name", "Ada Lovelace",
		"--email", "ada@example.com",
		"--company", "Analytical Engines Inc.",
		"--address", "500 Terry Francois Blvd, San Francisco, CA 94158",
		"--items-json", `[{"shop_slug":"starbucks","drink_slug":"iced-caramel-macchiato","size":"grande","quantity":1}]`,
	)
	if err == nil {
		t.Fatal("order command succeeded")
	}
	if !strings.Contains(err.Error(), "provide a complete and specific delivery address") {
		t.Fatalf("error = %v, want address guidance", err)
	}
}

func TestOrderCommandOutsideZoneOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(models.CreateOrderResponse{
			Order: models.OrderSummary{
				ID:            "drk_outside",
				Status:        "pending_payment",
				ZoneStatus:    "outside",
				SubtotalCents: 475,
				TotalCents:    475,
				Currency:      "usd",
			},
			AddressCheck: models.AddressCheck{
				FormattedAddress: "501 Stanyan St, San Francisco, CA 94117",
			},
			CheckoutSessionID: "cs_test_outside",
			CheckoutURL:       "https://checkout.stripe.com/c/pay/cs_test_outside",
		}); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	out, err := executeForTestWithAPI(
		t,
		"test-version",
		server.URL,
		"order",
		"--name", "Ada Lovelace",
		"--email", "ada@example.com",
		"--company", "Analytical Engines Inc.",
		"--address", "501 Stanyan St, San Francisco, CA 94117",
		"--items-json", `[{"shop_slug":"starbucks","drink_slug":"americano","size":"grande","quantity":1}]`,
	)
	if err != nil {
		t.Fatalf("order command failed: %v", err)
	}
	if strings.Contains(out, "Zone: outside") {
		t.Fatalf("output should not contain raw outside zone line:\n%s", out)
	}
	if !strings.Contains(out, "outside the default delivery zone") {
		t.Fatalf("output missing outside-zone review warning:\n%s", out)
	}
}
