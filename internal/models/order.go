package models

import (
	"fmt"
	"strings"

	"github.com/mupt-ai/dari-coffee-cli/internal/openapi"
)

type CreateOrderRequest = openapi.CreateOrderRequest
type CreateOrderItemRequest = openapi.CreateOrderItemRequest
type CreateOrderResponse = openapi.CreateOrderResponse
type OrderSummary = openapi.OrderSummary

func NormalizeCreateOrderRequest(req CreateOrderRequest) CreateOrderRequest {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Company = strings.TrimSpace(req.Company)
	req.Address = strings.TrimSpace(req.Address)
	req.DeliveryNotes = strings.TrimSpace(req.DeliveryNotes)
	for i := range req.Items {
		req.Items[i].ShopSlug = strings.TrimSpace(req.Items[i].ShopSlug)
		req.Items[i].DrinkSlug = strings.TrimSpace(req.Items[i].DrinkSlug)
		req.Items[i].Size = strings.TrimSpace(req.Items[i].Size)
		req.Items[i].Modifications = strings.TrimSpace(req.Items[i].Modifications)
	}
	req.CustomerNote = strings.TrimSpace(req.CustomerNote)
	return req
}

func ValidateCreateOrderRequest(req CreateOrderRequest) error {
	required := []struct {
		field string
		value string
	}{
		{field: "name", value: req.Name},
		{field: "email", value: req.Email},
		{field: "company", value: req.Company},
		{field: "address", value: req.Address},
	}
	for _, field := range required {
		if field.value == "" {
			return fmt.Errorf("%s is required", field.field)
		}
	}
	if len(req.Items) == 0 {
		return fmt.Errorf("items is required")
	}
	for i, item := range req.Items {
		if err := ValidateCreateOrderItemRequest(item); err != nil {
			return fmt.Errorf("items[%d]: %w", i, err)
		}
	}
	if err := validateSingleShop(req.Items); err != nil {
		return err
	}
	return nil
}

func ValidateCreateOrderItemRequest(item CreateOrderItemRequest) error {
	required := []struct {
		field string
		value string
	}{
		{field: "shop_slug", value: item.ShopSlug},
		{field: "drink_slug", value: item.DrinkSlug},
		{field: "size", value: item.Size},
	}
	for _, field := range required {
		if field.value == "" {
			return fmt.Errorf("%s is required", field.field)
		}
	}
	if item.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	return nil
}

func validateSingleShop(items []CreateOrderItemRequest) error {
	shopSlug := items[0].ShopSlug
	for i, item := range items[1:] {
		if item.ShopSlug != shopSlug {
			return fmt.Errorf("items[%d]: all items must use the same shop_slug", i+1)
		}
	}
	return nil
}
