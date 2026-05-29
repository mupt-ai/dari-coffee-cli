package cli

import (
	"bytes"
	"strings"
	"testing"
)

func executeForTest(version string, args ...string) (string, error) {
	cmd := newRootCommand(version)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func executeForTestWithAPI(t *testing.T, version string, baseURL string, args ...string) (string, error) {
	t.Helper()
	t.Setenv(devBaseURLEnv, baseURL)
	return executeForTest(version, args...)
}

func TestDefaultAPIBaseURL(t *testing.T) {
	t.Setenv(devBaseURLEnv, "")
	if got := defaultAPIBaseURL(); got != defaultServiceBaseURL {
		t.Fatalf("defaultAPIBaseURL() = %q, want %q", got, defaultServiceBaseURL)
	}
}

func TestDefaultAPIBaseURLUsesDevOverride(t *testing.T) {
	t.Setenv(devBaseURLEnv, "http://127.0.0.1:8080")
	if got, want := defaultAPIBaseURL(), "http://127.0.0.1:8080"; got != want {
		t.Fatalf("defaultAPIBaseURL() = %q, want %q", got, want)
	}
}

func TestHelpDoesNotExposeAPIURL(t *testing.T) {
	out, err := executeForTest("test-version", "--help")
	if err != nil {
		t.Fatalf("help command failed: %v", err)
	}
	if strings.Contains(out, "--api-url") {
		t.Fatalf("help output should not expose --api-url:\n%s", out)
	}
}
