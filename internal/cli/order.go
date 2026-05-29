package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mupt-ai/dari-coffee-cli/internal/api"
	"github.com/mupt-ai/dari-coffee-cli/internal/models"
	"github.com/spf13/cobra"
)

func newOrderCommand() *cobra.Command {
	var req models.CreateOrderRequest
	var itemsJSON string
	cmd := &cobra.Command{
		Use:   "order",
		Short: "Create a coffee delivery order",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := parseOrderItemsJSON(itemsJSON)
			if err != nil {
				return err
			}
			req.Items = items
			req = models.NormalizeCreateOrderRequest(req)
			if err := models.ValidateCreateOrderRequest(req); err != nil {
				return err
			}

			baseURL := apiBaseURL()
			client, err := api.New(baseURL)
			if err != nil {
				return err
			}
			out, err := client.CreateOrder(cmd.Context(), req)
			if err != nil {
				return formatCreateOrderError(baseURL, err)
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), formatCreateOrderResponse(out))
			return err
		},
	}
	cmd.Flags().StringVar(&req.Name, "name", "", "Customer name")
	cmd.Flags().StringVar(&req.Email, "email", "", "Customer email")
	cmd.Flags().StringVar(&req.Phone, "phone", "", "Customer phone number")
	cmd.Flags().StringVar(&req.Company, "company", "", "Customer company")
	cmd.Flags().StringVar(&req.Address, "address", "", "Complete delivery address")
	cmd.Flags().StringVar(&req.DeliveryNotes, "delivery-notes", "", "Delivery notes for the courier")
	cmd.Flags().StringVar(&itemsJSON, "items-json", "", "JSON array of order items from the menu command")
	cmd.Flags().StringVar(&req.CustomerNote, "note", "", "Additional customer note")
	return cmd
}

func parseOrderItemsJSON(raw string) ([]models.CreateOrderItemRequest, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("--items-json is required")
	}
	var items []models.CreateOrderItemRequest
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&items); err != nil {
		return nil, fmt.Errorf("parse --items-json: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, errors.New("parse --items-json: must contain exactly one JSON array")
	}
	return items, nil
}

func formatCreateOrderError(baseURL string, err error) error {
	var apiErr api.Error
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case models.ErrorCodeAddressUnresolved:
			return fmt.Errorf("create order at %s: address_unresolved: provide a complete and specific delivery address, including street number, street name, city, state, and ZIP if available", baseURL)
		case models.ErrorCodeGeocoderUnavailable:
			return fmt.Errorf("create order at %s: geocoder_unavailable: address validation is temporarily unavailable; try again later", baseURL)
		case models.ErrorCodeServiceOff:
			return fmt.Errorf("create order at %s: service_off: Dari Coffee is not accepting checkout requests right now", baseURL)
		case models.ErrorCodeOutsideHours:
			return fmt.Errorf("create order at %s: outside_hours: Dari Coffee is outside service hours", baseURL)
		case models.ErrorCodeCheckoutUnavailable:
			return fmt.Errorf("create order at %s: checkout_unavailable: Dari Coffee is not accepting checkout requests right now", baseURL)
		case models.ErrorCodeInvalidOrder:
			return fmt.Errorf("create order at %s: invalid_order: %s", baseURL, apiErr.Message)
		case models.ErrorCodeMenuItemUnavailable:
			return fmt.Errorf("create order at %s: menu_item_unavailable: choose a currently available shop, drink, and size from `dari-coffee menu`", baseURL)
		case models.ErrorCodePaymentUnavailable:
			return fmt.Errorf("create order at %s: payment_unavailable: Stripe Checkout is temporarily unavailable; try again later", baseURL)
		}
	}
	return fmt.Errorf("create order at %s: %w", baseURL, err)
}

func formatCreateOrderResponse(out models.CreateOrderResponse) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Order created: %s\n", out.Order.ID)
	fmt.Fprintf(&b, "Status: %s\n", out.Order.Status)
	fmt.Fprintf(&b, "Total: %s %s\n", formatUSD(out.Order.TotalCents), out.Order.Currency)
	if out.AddressCheck.FormattedAddress != "" {
		fmt.Fprintf(&b, "Address resolved to: %s\n", out.AddressCheck.FormattedAddress)
	}
	if out.Order.ZoneStatus == string(ordersZoneOutside) {
		fmt.Fprintln(&b, "This address is outside the default delivery zone and may be denied after review.")
	}
	if out.CheckoutURL != "" {
		fmt.Fprintf(&b, "Checkout: %s\n", out.CheckoutURL)
		fmt.Fprintln(&b, "Authorize payment in Stripe. Dari captures only after accepting the order.")
	} else {
		fmt.Fprintln(&b, "Payment is not wired in this build; no charge was made.")
	}
	return b.String()
}

const (
	ordersZoneOutside = "outside"
)
