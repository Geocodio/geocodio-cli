package cli

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// defaultDistanceMode matches the API's own default. Straightline costs
// 1 credit per calculation; driving costs 2.
const defaultDistanceMode = "straightline"

// costNoticeThreshold is the number of distance calculations above which the
// CLI prints an estimate before submitting, so a large run doesn't turn into a
// surprise bill.
const costNoticeThreshold = 10000

// usageLimitURL points at the guide for capping spend on an account.
const usageLimitURL = "https://www.geocod.io/guides/set-a-usage-limit"

// creditsPerCalculation returns the credit cost of a single distance
// calculation in the given routing mode.
func creditsPerCalculation(mode string) int {
	if mode == "driving" {
		return 2
	}
	return 1
}

// printCostNotice writes an estimate to w when a distance request would run
// more than costNoticeThreshold calculations. It is written before the request
// is submitted, and goes to stderr so it never pollutes --json output.
func printCostNotice(w io.Writer, origins, destinations int, mode string) {
	calculations := origins * destinations
	if calculations <= costNoticeThreshold {
		return
	}

	if mode == "" {
		mode = defaultDistanceMode
	}
	perCalculation := creditsPerCalculation(mode)
	credits := calculations * perCalculation

	plural := "s"
	if perCalculation == 1 {
		plural = ""
	}

	fmt.Fprintf(w, "Note: %s origins x %s destinations = %s distance calculations in %s mode\n",
		formatThousands(origins), formatThousands(destinations), formatThousands(calculations), mode)
	fmt.Fprintf(w, "      %d credit%s per calculation, about %s credits, plus one lookup per address that needs geocoding.\n",
		perCalculation, plural, formatThousands(credits))
	if mode != "driving" {
		fmt.Fprintln(w, "      --mode driving costs 2 credits per calculation.")
	}
	fmt.Fprintf(w, "      Set a usage limit: %s\n", usageLimitURL)
}

// formatThousands renders an integer with comma separators (1234567 →
// "1,234,567") so large calculation counts stay readable.
func formatThousands(n int) string {
	s := strconv.Itoa(n)
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}

	var b strings.Builder
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(digit)
	}
	return sign + b.String()
}
