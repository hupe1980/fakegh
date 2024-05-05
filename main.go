package main

import (
	"github.com/hupe1980/fakegh/cmd"
)

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
