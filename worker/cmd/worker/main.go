package main

import (
	"time"

	"github.com/jurekbarth/pup/worker/internal/app"
	"github.com/jurekbarth/pup/worker/internal/util"
)

var version = "master"

func main() {

	err := run()

	// Wait for the logs to get published :)
	duration := time.Second * 2
	time.Sleep(duration)

	if err == nil {
		return
	}

	util.Fatal(err)

}

// run the cli.
func run() error {
	return app.Run()
}
