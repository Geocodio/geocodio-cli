package output

import (
	"fmt"
	"github.com/fatih/color"
	"io"
)

func Success(w io.Writer, message string) {
	fmt.Fprintln(w, color.GreenString("✅ Success:") + " " + message)
}
