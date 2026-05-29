package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestCreateOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/orders" {
			t.Fatalf("path = %q, want /v1/orders", r.URL.Path)
		}
		var req models.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if got, want := len(req.Items), 2; got != want {
			t.Fatalf("item count = %d, want %d", got, want)
		}
		if got, want := req.Items[0].ShopSlug, "starbucks"; got != want {
			t.Fatalf("shop slug = %q, want %q", got, want)
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(models.CreateOrderResponse{
			Order: models.OrderSummary{
				ID:            "drk_test",
				Status:        "pending_payment",
				ZoneStatus:    "inside",
				TotalCents:    625,
				SubtotalCents: 625,
				Currency:      "usd",
			},
			CheckoutSessionID: "cs_test_order",
			CheckoutURL:       "https://checkout.stripe.com/c/pay/cs_test_order",
		}); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	order, err := client.CreateOrder(context.Background(), models.CreateOrderRequest{
		Name:    "Ada",
		Email:   "ada@example.com",
		Company: "Analytical Engines Inc.",
		Address: "1 Market St, San Francisco, CA 94105",
		Items: []models.CreateOrderItemRequest{
			{ShopSlug: "starbucks", DrinkSlug: "americano", Size: "grande", Quantity: 1},
			{ShopSlug: "starbucks", DrinkSlug: "mocha-frappuccino", Size: "grande", Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if got, want := order.Order.ID, "drk_test"; got != want {
		t.Fatalf("order ID = %q, want %q", got, want)
	}
	if got, want := order.CheckoutURL, "https://checkout.stripe.com/c/pay/cs_test_order"; got != want {
		t.Fatalf("checkout URL = %q, want %q", got, want)
	}
}

func TestAPIErrorIncludesServerMessage(t *testing.T) {
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

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, err = client.CreateOrder(context.Background(), models.CreateOrderRequest{})
	if err == nil {
		t.Fatal("CreateOrder succeeded")
	}
	var apiErr Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want api.Error", err)
	}
	if got, want := apiErr.Message, "address_unresolved"; got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
	if got, want := apiErr.Code, models.ErrorCodeAddressUnresolved; got != want {
		t.Fatalf("code = %q, want %q", got, want)
	}
}
