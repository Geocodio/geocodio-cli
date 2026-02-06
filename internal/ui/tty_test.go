package ui

import (
	"os"
	"testing"
)

func TestIsTTY(t *testing.T) {
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Skip("cannot open /dev/null")
	}
	defer devNull.Close()

	if IsTTY(devNull) {
		t.Error("IsTTY(/dev/null) = true, want false")
	}
}

func TestColorEnabled_NoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	os.Unsetenv("FORCE_COLOR")

	if ColorEnabled() {
		t.Error("ColorEnabled() = true with NO_COLOR set, want false")
	}
}

func TestColorEnabled_ForceColorEnv(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	t.Setenv("FORCE_COLOR", "1")

	if !ColorEnabled() {
		t.Error("ColorEnabled() = false with FORCE_COLOR set, want true")
	}
}

func TestColorEnabled_NoColorTakesPrecedence(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "1")

	if ColorEnabled() {
		t.Error("ColorEnabled() = true with both NO_COLOR and FORCE_COLOR set, want false (NO_COLOR takes precedence)")
	}
}
