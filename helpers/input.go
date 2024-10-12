package helpers

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))

	fmt.Println()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(passwordBytes)), nil
}
