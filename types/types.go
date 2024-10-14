package types

type Command struct {
	Name        string
	Description string
	Usage       string
	ArgCount    int
	Execute     func(args []string)
}

type CommandDistance struct {
	Name  string
	Score int
}

type PlockEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
