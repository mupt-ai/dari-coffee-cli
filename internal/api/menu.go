package api

import (
	"context"
	"net/http"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func (c *Client) Menu(ctx context.Context) (models.Menu, error) {
	var out models.Menu
	if err := c.doJSON(ctx, http.MethodGet, "/v1/menu", &out); err != nil {
		return models.Menu{}, err
	}
	return out, nil
}
