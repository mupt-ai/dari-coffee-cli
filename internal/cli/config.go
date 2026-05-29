package cli

import "os"

const (
	defaultServiceBaseURL = "https://coffee.dari.dev"
	devBaseURLEnv         = "DARI_COFFEE_DEV_BASE_URL"
)

func defaultAPIBaseURL() string {
	if baseURL := os.Getenv(devBaseURLEnv); baseURL != "" {
		return baseURL
	}
	return defaultServiceBaseURL
}
