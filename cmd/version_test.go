package cmd

import (
	"strings"
	"testing"
)

func TestGetVersionDefault(t *testing.T) {
	version = "v1.2.3"
	gitCommit = ""
	got := getVersion()
	if got != "v1.2.3" {
		t.Fatalf("getVersion() = %q, want %q", got, "v1.2.3")
	}
}

func TestGetVersionWithCommit(t *testing.T) {
	version = "v1.2.3"
	gitCommit = "abcdef1234567890"
	got := getVersion()
	if !strings.HasPrefix(got, "v1.2.3") {
		t.Fatalf("getVersion() = %q, want prefix %q", got, "v1.2.3")
	}
	if !strings.Contains(got, "abcdef1") {
		t.Fatalf("getVersion() = %q, want it to contain commit %q", got, "abcdef1")
	}
}

func TestGetVersionExactlySevenCharCommit(t *testing.T) {
	version = "v0.0.0"
	gitCommit = "1234567"
	got := getVersion()
	// The code only includes commit when len > 7, so exactly 7 should not be appended
	if got != "v0.0.0" {
		t.Fatalf("getVersion() with 7-char commit = %q, want %q", got, "v0.0.0")
	}
}