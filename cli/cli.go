package cli

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/pass-lock/cmd"
	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var commands = map[string]types.Command{
	"setup":         cmd.SetupCommand,
	"set":           cmd.SetCommand,
	"get":           cmd.GetCommand,
	"delete":        cmd.DeleteCommand,
	"get-all":       cmd.GetAllCommand,
	"edit":          cmd.EditCommand,
	"self-destruct": cmd.SelfDestructCommand,
}

// Setup parses input and executes the appropriate command.
func Setup() {
	args := os.Args[1:] // Ignore program name.

	// Print general help if no command or "help" is provided.
	if len(args) == 0 || args[0] == "help" {
		helpers.PrintHelp(commands)
		return
	}

	// Check if the command exists.
	cmdName := args[0]
	cmd, exists := commands[cmdName]
	if !exists {
		helpers.HandleUnknownCommand(commands, cmdName) // Handle invalid commands.
		return
	}

	// Show command-specific help if the user requests it.
	if len(args) > 1 && helpers.IsHelpFlag(args[1]) {
		helpers.PrintCommandHelp(cmd)
		return
	}

	// Validate argument count.
	expectedArgs := cmd.ArgCount
	providedArgs := len(args) - 1 // Exclude the command itself.

	if providedArgs < expectedArgs {
		fmt.Printf("Error: '%s' needs at least %d arguments, but got %d.\n", cmdName, expectedArgs, providedArgs)
		fmt.Printf("See 'passlock %s --help' for usage.\n", cmdName)
		return
	}

	if providedArgs > expectedArgs {
		fmt.Printf("Warning: '%s' received more arguments than expected.\n", cmdName)
	}

	// Execute the command with provided arguments.
	cmd.Execute(args[1:])
}
