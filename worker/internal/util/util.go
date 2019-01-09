package util

import (
	"fmt"
	"os"
)

// Fatal error.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "\n     %s %s\n\n", "Error:", err)
	os.Exit(1)
}
