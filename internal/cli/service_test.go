package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mupt-ai/dari-coffee-cli/internal/models"
)

func TestServiceCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/service" {
			t.Fatalf("path = %q, want /v1/service", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(models.ServiceStatus{
			ServiceEnabled:    false,
			CheckoutAvailable: false,
			Reason:            "service_off",
			OpenTime:          "10:30",
			CloseTime:         "17:30",
		}); err != nil {
			t.Fatalf("write service response: %v", err)
		}
	}))
	defer server.Close()

	out, err := executeForTestWithAPI(t, "test-version", server.URL, "service")
	if err != nil {
		t.Fatalf("service command failed: %v", err)
	}

	for _, want := range []string{
		"not accepting checkout requests",
		"Reason: service_off",
		"Hours: 10:30-17:30",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("service output does not contain %q:\n%s", want, out)
		}
	}
}
