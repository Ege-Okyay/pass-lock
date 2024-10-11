package types

type Command struct {
	Name        string
	Description string
	Usage       string
	Execute     func(args []string)
}

type CommandDistance struct {
	Name  string
	Score int
}
