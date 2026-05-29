package cli

import "testing"

func TestVersionCommand(t *testing.T) {
	out, err := executeForTest("test-version", "version")
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	want := "dari-coffee test-version\n"
	if out != want {
		t.Fatalf("version output = %q, want %q", out, want)
	}
}

func TestVersionFlag(t *testing.T) {
	out, err := executeForTest("test-version", "--version")
	if err != nil {
		t.Fatalf("version flag failed: %v", err)
	}

	want := "dari-coffee test-version\n"
	if out != want {
		t.Fatalf("version output = %q, want %q", out, want)
	}
}
