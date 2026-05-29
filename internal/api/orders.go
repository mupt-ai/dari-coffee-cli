package api

import (
	"context"
	"net/http"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func (c *Client) CreateOrder(ctx context.Context, req models.CreateOrderRequest) (models.CreateOrderResponse, error) {
	var out models.CreateOrderResponse
	if err := c.doJSONBody(ctx, http.MethodPost, "/v1/orders", req, &out); err != nil {
		return models.CreateOrderResponse{}, err
	}
	return out, nil
}
