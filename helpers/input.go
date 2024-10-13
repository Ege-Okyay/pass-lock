package helpers

import (
	"errors"
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

func ValidateInput(input, fieldName string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New(fieldName + " cannot be empty.")
	}
	return nil
}
