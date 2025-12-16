package main

import (
	"os"

	"github.com/clbiggs/image-builder/internal/cli"
)

func main() {
	os.Exit(cli.Main(os.Args))
}
