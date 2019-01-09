package root

import (
	"os"
	"runtime"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/jurekbarth/pup/client"
	"github.com/jurekbarth/pup/client/event"
	"github.com/jurekbarth/pup/client/reporter"
)

// Cmd is the root command.
var Cmd = kingpin.New("pup", "")

// Command registers a command.
var Command = Cmd.Command

// Init function.
var Init func() (*client.Config, *client.Project, error)

func init() {
	log.SetHandler(cli.Default)

	Cmd.Example(`pup`, "Deploy the project to the *latest* folder")
	Cmd.Example(`pup deploy <foldername>`, "Deploy the project to <foldername>")
	Cmd.Example(`pup url`, "Show the latest folder url.")
	workdir := Cmd.Flag("chdir", "Change working directory.").Default(".").Short('C').String()
	verbose := Cmd.Flag("verbose", "Enable verbose log output.").Short('v').Bool()
	// format := Cmd.Flag("format", "Output formatter.").Default("text").String()

	Cmd.PreAction(func(ctx *kingpin.ParseContext) error {
		os.Chdir(*workdir)

		if *verbose {
			log.SetLevel(log.DebugLevel)
			log.Debugf("pup version %s (os: %s, arch: %s)", Cmd.GetVersion(), runtime.GOOS, runtime.GOARCH)
		}

		Init = func() (*client.Config, *client.Project, error) {
			c, err := client.ReadConfig("pup.json")
			if err != nil {
				return nil, nil, errors.Wrap(err, "reading config")
			}

			events := make(event.Events)
			p := client.New(c, events)

			go reporter.Text(events)

			return c, p, nil
		}

		return nil
	})
}
