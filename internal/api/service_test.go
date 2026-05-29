package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v1/service" {
			t.Fatalf("path = %q, want /v1/service", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(models.ServiceStatus{
			ServiceEnabled:    true,
			CheckoutAvailable: true,
			OpenTime:          "10:30",
			CloseTime:         "17:30",
		}); err != nil {
			t.Fatalf("write service response: %v", err)
		}
	}))
	defer server.Close()

	client, err := New(server.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	status, err := client.Service(context.Background())
	if err != nil {
		t.Fatalf("Service: %v", err)
	}
	if !status.CheckoutAvailable {
		t.Fatalf("CheckoutAvailable = false, want true")
	}
	if got, want := status.OpenTime, "10:30"; got != want {
		t.Fatalf("OpenTime = %q, want %q", got, want)
	}
}
