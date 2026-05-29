package api

import "testing"

func TestNewValidatesBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
		wantErr bool
	}{
		{
			name:    "http",
			baseURL: "http://localhost:8080/",
			want:    "http://localhost:8080",
		},
		{
			name:    "https",
			baseURL: " https://coffee.dari.dev ",
			want:    "https://coffee.dari.dev",
		},
		{
			name:    "empty",
			baseURL: "",
			wantErr: true,
		},
		{
			name:    "relative",
			baseURL: "localhost:8080",
			wantErr: true,
		},
		{
			name:    "query",
			baseURL: "https://coffee.dari.dev?debug=true",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := New(tc.baseURL)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("New(%q) succeeded, want error", tc.baseURL)
				}
				return
			}
			if err != nil {
				t.Fatalf("New(%q): %v", tc.baseURL, err)
			}
			if client.baseURL != tc.want {
				t.Fatalf("baseURL = %q, want %q", client.baseURL, tc.want)
			}
		})
	}
}
