package logging

import (
	"os"
	"strings"
	"testing"
)

func resetState() {
	lvl = INFO
}

func TestInitValidLevels(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  level
	}{
		{"trace", "TRACE", TRACE},
		{"debug", "DEBUG", DEBUG},
		{"info", "INFO", INFO},
		{"warn", "WARN", WARN},
		{"error", "ERROR", ERROR},
		{"case insensitive", "trace", TRACE},
		{"mixed case", "DeBuG", DEBUG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetState()
			if err := Init(tt.level); err != nil {
				t.Fatalf("Init(%q) unexpected error: %v", tt.level, err)
			}
			if lvl != tt.want {
				t.Fatalf("Init(%q) set lvl = %d, want %d", tt.level, lvl, tt.want)
			}
		})
	}
}

func TestInitInvalidLevel(t *testing.T) {
	resetState()
	err := Init("bogus")
	if err == nil {
		t.Fatal("Init(\"bogus\") expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid log level") {
		t.Fatalf("Init(\"bogus\") error = %q, want 'invalid log level'", err)
	}
}

func TestInitEmptyStringDefaultsToInfo(t *testing.T) {
	os.Unsetenv("LOG_LEVEL")
	resetState()
	if err := Init(""); err != nil {
		t.Fatalf("Init(\"\") unexpected error: %v", err)
	}
	if lvl != INFO {
		t.Fatalf("Init(\"\") set lvl = %d, want INFO (%d)", lvl, INFO)
	}
}

func TestInitEmptyStringReadsFromEnv(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("LOG_LEVEL")
	resetState()
	if err := Init(""); err != nil {
		t.Fatalf("Init(\"\") unexpected error: %v", err)
	}
	if lvl != DEBUG {
		t.Fatalf("Init(\"\") with LOG_LEVEL=DEBUG set lvl = %d, want DEBUG (%d)", lvl, DEBUG)
	}
}

func TestInitExplicitOverridesEnv(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("LOG_LEVEL")
	resetState()
	if err := Init("WARN"); err != nil {
		t.Fatalf("Init(\"WARN\") unexpected error: %v", err)
	}
	if lvl != WARN {
		t.Fatalf("Init(\"WARN\") with LOG_LEVEL=DEBUG set lvl = %d, want WARN (%d)", lvl, WARN)
	}
}
