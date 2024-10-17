package types

// Command represents a CLI command with its name, description, usage, argument count,
// and the function to be executed when the command is called.
type Command struct {
	Name        string
	Description string
	Usage       string
	ArgCount    int
	Execute     func(args []string) // Function pointer to execute the command logic.
}

// CommandDistance represents a command name and its similarity score (used for suggestions).
type CommandDistance struct {
	Name  string // The command name.
	Score int    // Levenshtein distance score.
}

// PlockEntry stores a key-value pair, both encrypted and used in the password vault.
type PlockEntry struct {
	Key   string `json:"key"`   // The key associated with the value.
	Value string `json:"value"` // The encrypted value.
}
