package main

import (
	"fmt"

	"github.com/martient/Bifrost-env-manager/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(fmt.Sprintf("version %s, commit %s, built at %s\n", version, commit, date))
}
