package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestMenu(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v1/menu" {
			t.Fatalf("path = %q, want /v1/menu", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(models.Menu{
			Title: "Dari Coffee CLI Menu",
			Shops: []models.Shop{
				{
					Slug: "starbucks",
					Name: "Starbucks",
					Hours: models.ShopHours{
						OpenTime:      "07:00",
						CloseTime:     "18:00",
						LastOrderTime: "17:00",
					},
				},
			},
		}); err != nil {
			t.Fatalf("write menu response: %v", err)
		}
	}))
	defer server.Close()

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	menu, err := client.Menu(context.Background())
	if err != nil {
		t.Fatalf("Menu: %v", err)
	}
	if got, want := menu.Shops[0].Slug, "starbucks"; got != want {
		t.Fatalf("shop slug = %q, want %q", got, want)
	}
	if got, want := menu.Shops[0].Hours.LastOrderTime, "17:00"; got != want {
		t.Fatalf("last order time = %q, want %q", got, want)
	}
}
