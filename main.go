package main

import (
	"os"

	"github.com/Ege-Okyay/pass-lock/cli"
)

func main() {
	args := os.Args[1:]

	cli.RunCommand(args)
}
