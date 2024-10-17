package helpers

import (
	"fmt"
	"sort"

	"github.com/Ege-Okyay/pass-lock/types"
)

// PrintHelp displays the general help message with a list of available commands,
// their descriptions, and usage examples.
func PrintHelp(commands map[string]types.Command) {
	PrintBanner("Passlock - Secure Key/Value Store")
	fmt.Println("USAGE:")
	fmt.Println("\tpasslock <command> [arguments]")

	// List available commands.
	fmt.Println("AVAILABLE COMMANDS:")
	for _, cmd := range commands {
		fmt.Printf("\t%-15s %s\n", cmd.Name, cmd.Description)
	}

	fmt.Println("FLAGS:")
	fmt.Println("\t-h, --help")

	// Provide command usage examples.
	fmt.Println("EXAMPLES:")
	fmt.Println("\tpasslock setup \tInitialize the vault with a master password.")
	fmt.Println("\tpasslock set apiKey secret123 \tStore a new key-value pair.")
	fmt.Println("\tpasslock get apiKey \tRetrieve a specific key-value pair.")
	fmt.Println("\tpasslock get-all \tRetrieve all key-value pairs.")
	fmt.Println("\tpasslock delete apiKey \tDelete a specific key-value pair.")
	fmt.Println("\tpasslock edit apiKey \tEdit the value of a key.")
	fmt.Println("\tpasslock self-destruct \tDelete all data and configurations.")

	fmt.Println("\nUse 'passlock <command> --help' for more information about a command.")
}

// PrintCommandHelp shows help for a specific command, including the name,
// description, and usage format.
func PrintCommandHelp(cmd types.Command) {
	fmt.Println("NAME:")
	fmt.Printf("\t%s\n", cmd.Name)

	fmt.Println("DESCRIPTION:")
	fmt.Printf("\t%s\n", cmd.Description)

	fmt.Println("USAGE:")
	fmt.Printf("\t%s\n", cmd.Usage)
}

// findClosestCommands returns the most similar commands to the given input
// by comparing the Levenshtein distance between command names.
func findClosestCommands(commands map[string]types.Command, cmdName string, maxResults int) []string {
	str1 := []rune(cmdName)
	var commandDistances []types.CommandDistance

	// Calculate the similarity score between the input and command names.
	for name := range commands {
		str2 := []rune(name)
		score := Levenshtein(str1, str2)

		commandDistances = append(commandDistances, types.CommandDistance{Name: name, Score: score})
	}

	// Sort commands by similarity score (ascending).
	sort.Slice(commandDistances, func(i, j int) bool {
		return commandDistances[i].Score < commandDistances[j].Score
	})

	// Collect the closest commands based on the maximum results limit.
	var closestCommands []string
	for i := 0; i < min(maxResults, len(commandDistances)); i++ {
		closestCommands = append(closestCommands, commandDistances[i].Name)
	}

	return closestCommands
}

// HandleUnknownCommand suggests similar commands when the user enters an
// unrecognized command.
func HandleUnknownCommand(commands map[string]types.Command, cmdName string) {
	closestCommands := findClosestCommands(commands, cmdName, 2)

	// Notify the user of the unrecognized command.
	fmt.Printf("passlock: '%s' is not a passlock command. See 'passlock help'.\n", cmdName)

	// Suggest the most similar commands.
	if len(closestCommands) > 0 {
		fmt.Println("\nThe most similar commands are:")
		for _, command := range closestCommands {
			fmt.Printf("\t%s\n", command)
		}
	}
}

// IsHelpFlag checks if the provided argument matches any of the help flags.
func IsHelpFlag(arg string) bool {
	// Return true if the argument is a recognized help flag.
	return arg == "-h" || arg == "--help" || arg == "help"
}
