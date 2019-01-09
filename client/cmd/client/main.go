package main

import (
	"github.com/jurekbarth/pup/client/internal/app"
	_ "github.com/jurekbarth/pup/client/internal/deploy"
	"github.com/jurekbarth/pup/client/internal/util"
)

var version = "1"

func main() {

	err := run()

	if err == nil {
		return
	}

	util.Fatal(err)

}

// run the cli.
func run() error {
	return app.Run(version)
}
