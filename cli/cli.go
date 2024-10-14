package cli

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/pass-lock/cmd"
	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var commands = map[string]types.Command{
	"setup": cmd.SetupCommand,
	"set":   cmd.SetCommand,
}

func Setup() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" {
		helpers.PrintHelp(commands)
		return
	}

	cmdName := args[0]
	cmd, exists := commands[cmdName]

	if !exists {
		helpers.HandleUnknownCommand(commands, cmdName)
		return
	}

	if len(args) > 1 && helpers.IsHelpFlag(args[1]) {
		helpers.PrintCommandHelp(cmd)
		return
	}

	expectedArgs := cmd.ArgCount
	providedArgs := len(args) - 1

	if providedArgs < expectedArgs {
		fmt.Printf("Error: '%s' required at least %d arguments but got %d.\n", cmdName, expectedArgs, providedArgs)
		fmt.Printf("\nSee 'passlock %s --help' for usage.\n", cmdName)
		return
	}

	if providedArgs > expectedArgs {
		fmt.Printf("Warning: '%s' received more arguments than expected.\n", cmdName)
	}

	cmd.Execute(args[1:])
}
