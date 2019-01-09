package util

import (
	"fmt"
	"os"

	"github.com/jurekbarth/pup/client/internal/colors"
)

// Fatal error.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "\n     %s %s\n\n", colors.Red("Error:"), err)
	os.Exit(1)
}
