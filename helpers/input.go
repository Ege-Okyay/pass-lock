package helpers

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Ege-Okyay/passlock/types"
	"golang.org/x/term"
)

// ReadPassword reads a password from the user input, hiding the typed characters.
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	// Read password input from the terminal.
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Print a newline after input.
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(passwordBytes)), nil // Return trimmed input.
}

// ValidateInput ensures that the provided input is non-empty.
func ValidateInput(input, fieldName string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New(fieldName + " cannot be empty.")
	}
	return nil
}

// VerifyPasswordAndLoadData prompts the user for the master password,
// validates it, and loads encrypted data if the password matches.
func VerifyPasswordAndLoadData() ([]types.PlockEntry, []byte, error) {
	for {
		// Prompt the user to enter the password.
		password, err := ReadPassword("Password: ")
		if err != nil {
			return nil, nil, fmt.Errorf("error reading password: %w", err)
		}

		// Derive a cryptographic key from the password.
		derivedKey := DeriveKey(password)

		// Attempt to load the keys file using the derived key.
		keysEntries, err := LoadFromFile(filepath.Join(GetUserConfigDir(), "keys.plock"), derivedKey)
		if err != nil {
			// Inform the user if the password is incorrect.
			ErrorMessage("Incorrect password. Please try again.")
			PrintSeparator()
			continue // Retry password input.
		}

		// Search for the "password" entry in the keys file.
		for _, entry := range keysEntries {
			if entry.Key == "password" {
				// Decrypt the stored password to compare it with user input.
				storedPassword, err := Decrypt(entry.Value, derivedKey)
				if err != nil || storedPassword != password {
					ErrorMessage("Incorrect password. Please try again.")
					PrintSeparator()
					continue // Retry password input.
				}

				SuccessMessage("Password verified!")

				// Load the data file with the same derived key.
				dataEntries, err := LoadFromFile(filepath.Join(GetUserConfigDir(), "data.plock"), derivedKey)
				if err != nil {
					log.Fatalf("Error occurred while reading data.plock file: %v\n", err)
				}

				// Return the loaded data and derived key upon success.
				return dataEntries, derivedKey, nil
			}
		}

		// If no password entry is found, inform the user and exit the loop.
		ErrorMessage("Master password not found. Please ensure setup is completed.")
		return nil, nil, fmt.Errorf("master password missing")
	}
}

// ReadLine reads a single line of user input from the console and
// returns the trimmed input or an error.
func ReadLine() (string, error) {
	var input string

	// Capture user input from the console.
	_, err := fmt.Scanln(&input)
	if err != nil {
		// Handle the case where the user presses Enter without typing anything.
		if err.Error() == "unexpected newline" {
			return "", nil
		}
		// Return any other input errors.
		return "", err
	}

	// Return the entered input.
	return input, nil
}

// TrimNewline removes any leading or trailing whitespace from the given string.
func TrimNewline(s string) string {
	return strings.TrimSpace(s) // Trim spaces and newlines.
}
