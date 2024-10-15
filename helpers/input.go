package helpers

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Ege-Okyay/pass-lock/types"
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

func VerifyPassword() ([]types.PlockEntry, []byte, error) {
	for {
		password, err := ReadPassword("Password: ")
		if err != nil {
			return nil, nil, fmt.Errorf("error reading password: %w", err)
		}

		derivedKey := DeriveKey(password)

		entries, err := LoadFromFile(filepath.Join(GetAppDataPath(), "keys.plock"), derivedKey)
		if err != nil {
			ErrorMessage("Incorrect password. Please try again.")
			PrintSeparator()
			continue
		}

		for _, entry := range entries {
			if entry.Key == "password" {
				storedPassword, err := Decrypt(entry.Value, derivedKey)
				if err != nil || storedPassword != password {
					ErrorMessage("Incorrect password. Please try again.")
					PrintSeparator()
					continue
				}

				SuccessMessage("Password verified!")

				return entries, derivedKey, nil
			}
		}

		ErrorMessage("Master password not found. Please ensure setup is completed.")
		return nil, nil, fmt.Errorf("master password missing")
	}
}
