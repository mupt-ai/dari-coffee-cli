package api

import (
	"context"
	"net/http"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func (c *Client) Service(ctx context.Context) (models.ServiceStatus, error) {
	var out models.ServiceStatus
	if err := c.doJSON(ctx, http.MethodGet, "/v1/service", &out); err != nil {
		return models.ServiceStatus{}, err
	}
	return out, nil
}
