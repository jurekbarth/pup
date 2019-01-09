package app

import (
	"os"

	"github.com/jurekbarth/pup/client/internal/root"
)

// Run the command.
func Run(version string) error {
	root.Cmd.Version(version)
	_, err := root.Cmd.Parse(os.Args[1:])
	return err
}
