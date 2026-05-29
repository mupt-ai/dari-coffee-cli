package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestMenuCommand(t *testing.T) {
	server := newMenuTestServer(t)
	defer server.Close()

	out, err := executeForTestWithAPI(t, "test-version", server.URL, "menu")
	if err != nil {
		t.Fatalf("menu command failed: %v", err)
	}

	for _, want := range []string{
		"Dari Coffee CLI Menu",
		"Starbucks",
		"Hours: 07:00-18:00 (orders until 17:00)",
		"Iced Caramel Macchiato",
		"Sizes: grande $6.25",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("menu output does not contain %q:\n%s", want, out)
		}
	}
}

func TestMenuCommandJSON(t *testing.T) {
	server := newMenuTestServer(t)
	defer server.Close()

	out, err := executeForTestWithAPI(t, "test-version", server.URL, "menu", "--json")
	if err != nil {
		t.Fatalf("menu --json command failed: %v", err)
	}

	var got models.Menu
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("parse menu JSON output: %v\n%s", err, out)
	}
	if got.Shops[0].Slug != "starbucks" {
		t.Fatalf("shop slug = %q, want starbucks", got.Shops[0].Slug)
	}
	if got.Shops[0].Drinks[0].Slug != "iced-caramel-macchiato" {
		t.Fatalf("drink slug = %q, want iced-caramel-macchiato", got.Shops[0].Drinks[0].Slug)
	}
	if got.Shops[0].Drinks[0].Sizes[0].Name != "grande" {
		t.Fatalf("size name = %q, want grande", got.Shops[0].Drinks[0].Sizes[0].Name)
	}
}

func newMenuTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/menu" {
			t.Fatalf("path = %q, want /v1/menu", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(models.Menu{
			Title:       "Dari Coffee CLI Menu",
			Description: "Three shops. Three drinks each. No menu sprawl.",
			Shops: []models.Shop{
				{
					Slug: "starbucks",
					Name: "Starbucks",
					Hours: models.ShopHours{
						OpenTime:      "07:00",
						CloseTime:     "18:00",
						LastOrderTime: "17:00",
					},
					Drinks: []models.Drink{
						{
							Slug:        "iced-caramel-macchiato",
							Name:        "Iced Caramel Macchiato",
							Description: "Milk, espresso, vanilla, and caramel.",
							Sizes: []models.SizePrice{
								{Name: "grande", Currency: "usd", PriceCents: 625},
							},
						},
					},
				},
			},
		}); err != nil {
			t.Fatalf("write menu response: %v", err)
		}
	}))
	return server
}

func TestFormatMenuWithoutOpenShops(t *testing.T) {
	out := formatMenu(models.Menu{
		Title:       "Dari Coffee CLI Menu",
		Description: "Order coffee, delivered by Dari.",
	})

	if !strings.Contains(out, "No shops are currently taking Dari orders.") {
		t.Fatalf("menu output missing closed-shop message:\n%s", out)
	}
}

func TestFormatCents(t *testing.T) {
	tests := map[int]string{
		0:    "$0.00",
		5:    "$0.05",
		525:  "$5.25",
		-525: "-$5.25",
	}

	for cents, want := range tests {
		if got := formatUSD(cents); got != want {
			t.Fatalf("formatUSD(%d) = %q, want %q", cents, got, want)
		}
	}
}
