package api

import (
	"context"
	"net/http"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func (c *Client) CheckAddress(ctx context.Context, address string) (models.AddressCheck, error) {
	var out models.AddressCheck
	if err := c.doJSONBody(ctx, http.MethodPost, "/v1/address/check", models.AddressCheckRequest{
		Address: address,
	}, &out); err != nil {
		return models.AddressCheck{}, err
	}
	return out, nil
}
