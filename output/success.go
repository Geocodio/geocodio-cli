package output

import (
	"fmt"
	"github.com/fatih/color"
)

func Success(message string) {
	fmt.Println(color.GreenString("✅ Success:") + " " + message)
}
