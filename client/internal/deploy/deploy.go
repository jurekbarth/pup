package deploy

import (
	"os"
	"time"

	"github.com/jurekbarth/pup/client"
	"github.com/jurekbarth/pup/client/internal/root"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"
)

func init() {
	cmd := root.Command("deploy", "Deploy the project.").Default()
	stage := cmd.Arg("foldername", "Target folder name.").Default("latest").String()
	cmd.Example(`pup deploy`, "Deploy the project the staging environment.")
	cmd.Example(`pup deploy master`, "Deploy the project to the master folder.")

	cmd.Action(func(_ *kingpin.ParseContext) error {
		return deploy(*stage)
	})
}

func deploy(stage string) error {
	_, p, err := root.Init()

	if isMissingConfig(err) {
		return errors.New("cannot find ./pup.json configuration file")
	}
	// unrelated error
	if err != nil {
		return errors.Wrap(err, "initializing")
	}

	// TODO Time measurement?
	ts := time.Now()

	if err := p.Deploy(client.Deploy{
		Stage:     stage,
		Timestamp: ts,
	}); err != nil {
		return err
	}

	return nil
}

// isMissingConfig returns true if the error represents a missing up.json.
func isMissingConfig(err error) bool {
	err = errors.Cause(err)
	e, ok := err.(*os.PathError)
	return ok && e.Path == "pup.json"
}
