package ui

import (
	"bytes"
	"testing"
)

func TestNewSpinner_NonTTY(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner(&buf, "test message")

	if s.running {
		t.Error("NewSpinner on non-TTY should not be running initially")
	}
	if s.program != nil {
		t.Error("NewSpinner on non-TTY should not create a program")
	}
}

func TestSpinner_StartStopNonTTY(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner(&buf, "test message")

	s.Start()
	if s.running {
		t.Error("Start on non-TTY spinner should not set running to true")
	}

	s.Stop()
}

func TestSpinner_UpdateMessageNonTTY(t *testing.T) {
	var buf bytes.Buffer
	s := NewSpinner(&buf, "test message")

	s.UpdateMessage("new message")
}

func TestWithSpinner_NonTTY(t *testing.T) {
	var buf bytes.Buffer

	result, err := WithSpinner(&buf, "test", func() (string, error) {
		return "success", nil
	})

	if err != nil {
		t.Fatalf("WithSpinner() error = %v", err)
	}
	if result != "success" {
		t.Errorf("WithSpinner() = %q, want %q", result, "success")
	}
}
