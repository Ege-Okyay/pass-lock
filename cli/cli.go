package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/Ege-Okyay/pass-lock/cmd"
	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var commands = map[string]types.Command{
	"help":  cmd.HelpCommand,
	"setup": cmd.SetupCommand,
	"set":   cmd.SetCommand,
}

func Setup() {
	args := os.Args[1:]

	if len(args) == 0 {
		cmd.HelpCommand.Execute([]string{})
		return
	}

	cmdName := args[0]
	cmd, exists := commands[cmdName]

	if !exists {
		closestCommands := FindClosestCommands(cmdName, 5)

		fmt.Printf("passlock: '%s' is not a passlock command. See 'passlock help'.\n", cmdName)
		fmt.Printf("\nThe most similar commands are:\n")

		for _, command := range closestCommands {
			fmt.Printf("\t%s\n", command)
		}

		return
	}

	cmd.Execute(args[1:])
}

func PrintHelp(command types.Command) {
	fmt.Printf("Name: %s\n", command.Name)
	fmt.Printf("Description: %s\n", command.Description)
	fmt.Printf("Usage: %s\n", command.Usage)
}

func FindClosestCommands(cmdName string, maxResults int) []string {
	str1 := []rune(cmdName)
	var commandDistances []types.CommandDistance

	for name := range commands {
		str2 := []rune(name)
		score := helpers.Levenshtein(str1, str2)

		commandDistances = append(commandDistances, types.CommandDistance{Name: name, Score: score})
	}

	sort.Slice(commandDistances, func(i, j int) bool {
		return commandDistances[i].Score < commandDistances[j].Score
	})

	var closestCommands []string
	for i := 0; i < min(maxResults, len(commandDistances)); i++ {
		closestCommands = append(closestCommands, commandDistances[i].Name)
	}

	return closestCommands
}
