package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatThousands(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{0, "0"},
		{7, "7"},
		{999, "999"},
		{1000, "1,000"},
		{10000, "10,000"},
		{250000, "250,000"},
		{1234567, "1,234,567"},
		{-4200, "-4,200"},
	}

	for _, tt := range tests {
		if got := formatThousands(tt.in); got != tt.want {
			t.Errorf("formatThousands(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCreditsPerCalculation(t *testing.T) {
	if got := creditsPerCalculation("straightline"); got != 1 {
		t.Errorf("straightline = %d, want 1", got)
	}
	if got := creditsPerCalculation("driving"); got != 2 {
		t.Errorf("driving = %d, want 2", got)
	}
	if got := creditsPerCalculation(""); got != 1 {
		t.Errorf("empty mode = %d, want 1 (API defaults to straightline)", got)
	}
}

func TestPrintCostNoticeStaysQuietBelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	printCostNotice(&buf, 100, 100, "driving")

	if buf.Len() != 0 {
		t.Errorf("expected no notice for 10,000 calculations, got %q", buf.String())
	}
}

func TestPrintCostNoticeStraightline(t *testing.T) {
	var buf bytes.Buffer
	printCostNotice(&buf, 500, 500, "straightline")

	out := buf.String()
	for _, want := range []string{
		"500 origins x 500 destinations",
		"250,000 distance calculations",
		"straightline mode",
		"1 credit per calculation",
		"about 250,000 credits",
		"--mode driving costs 2 credits per calculation",
		usageLimitURL,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("notice missing %q; got:\n%s", want, out)
		}
	}
}

func TestPrintCostNoticeDriving(t *testing.T) {
	var buf bytes.Buffer
	printCostNotice(&buf, 500, 500, "driving")

	out := buf.String()
	if !strings.Contains(out, "2 credits per calculation, about 500,000 credits") {
		t.Errorf("driving notice should double the credits; got:\n%s", out)
	}
	if strings.Contains(out, "--mode driving costs") {
		t.Errorf("driving notice should not suggest driving mode; got:\n%s", out)
	}
}
