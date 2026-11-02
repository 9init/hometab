package main

import (
	"fmt"
	"os"

	"github.com/9init/hometab/internal/cli"
)

var version = "dev"

func main() {
	if err := cli.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, "hometab:", err)
		os.Exit(1)
	}
}
