package cmd

import (
	"fmt"

	"github.com/Ege-Okyay/pass-lock/types"
)

var TestCommand = types.Command{
	Name:        "test",
	Description: "test command",
	Usage:       "passlock test",
	Execute: func(args []string) {
		fmt.Println("Testing...")
	},
}
