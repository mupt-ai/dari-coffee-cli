package models

import (
	"strings"
	"testing"
)

func TestValidateCreateOrderRequestRequiresSingleShop(t *testing.T) {
	req := CreateOrderRequest{
		Name:    "Ada Lovelace",
		Email:   "ada@example.com",
		Company: "Analytical Engines Inc.",
		Address: "1 Market St, San Francisco, CA 94105",
		Items: []CreateOrderItemRequest{
			{ShopSlug: "starbucks", DrinkSlug: "americano", Size: "grande", Quantity: 1},
			{ShopSlug: "philz", DrinkSlug: "mint-mojito", Size: "medium", Quantity: 1},
		},
	}

	err := ValidateCreateOrderRequest(req)
	if err == nil {
		t.Fatal("ValidateCreateOrderRequest succeeded")
	}
	if !strings.Contains(err.Error(), "same shop_slug") {
		t.Fatalf("error = %v, want same shop_slug message", err)
	}
}
