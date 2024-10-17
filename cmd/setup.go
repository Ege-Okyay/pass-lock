package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

// Command to initialize the vault with a master password.
var SetupCommand = types.Command{
	Name:        "setup",
	Description: "Initialize the vault with a master password.",
	Usage:       "passlock setup",
	ArgCount:    0,
	Execute: func(args []string) {
		// Check if the vault is already set up.
		exists, err := helpers.CheckKeysFileExists()
		if err != nil {
			log.Fatalf("Error checking keys file: %v\n", err)
		}
		if exists {
			helpers.ErrorMessage("Setup already completed. Use the 'set' command to add entries.")
			return
		}

		helpers.PrintBanner("Welcome to Passlock")
		fmt.Println("Let's set up your password manager.")
		fmt.Println("\nSetup will begin shortly...")

		for {
			// Prompt the user to enter a password.
			password, err := helpers.ReadPassword("Password: ")
			if err != nil {
				log.Fatalf("Error reading password: %v\n", err)
			}

			// Validate the password input.
			if err := helpers.ValidateInput(password, "Password"); err != nil {
				helpers.ErrorMessage(err.Error())
				helpers.PrintSeparator()
				continue
			}

			// Confirm the password by prompting again.
			confirmPassword, err := helpers.ReadPassword("Confirm password: ")
			if err != nil {
				log.Fatalf("Error reading confirmation password: %v\n", err)
			}

			// Validate the confirmation input.
			if err := helpers.ValidateInput(confirmPassword, "Confirmation password"); err != nil {
				helpers.ErrorMessage(err.Error())
				helpers.PrintSeparator()
				continue
			}

			// Ensure both passwords match.
			if password != confirmPassword {
				helpers.ErrorMessage("Passwords do not match! Please try again.")
				helpers.PrintSeparator()
				continue
			}

			// Derive a cryptographic key from the password.
			derivedKey := helpers.DeriveKey(password)

			// Generate an AES encryption key.
			aesKey, err := helpers.GenerateAESKey()
			if err != nil {
				log.Fatalf("Error generating AES key: %v\n", err)
			}

			// Create the passlock data directory.
			passlockDir := helpers.GetUserConfigDir()
			if err := os.MkdirAll(passlockDir, os.ModePerm); err != nil {
				log.Fatalf("Error creating passlock directory: %v\n", err)
			}

			// Save the AES key in the keys file.
			if err := helpers.AddDataEntry(derivedKey, "keys.plock", "aes_key", string(aesKey)); err != nil {
				log.Fatalf("Error saving AES key: %v\n", err)
			}

			// Store the user's password securely.
			if err := helpers.AddDataEntry(derivedKey, "keys.plock", "password", password); err != nil {
				log.Fatalf("Error saving password: %v\n", err)
			}

			// Create an empty data file for future entries.
			dataFile := filepath.Join(passlockDir, "data.plock")
			if _, err := os.Create(dataFile); err != nil {
				log.Fatalf("Error creating data file: %v\n", err)
			}

			// Confirm successful setup.
			helpers.SuccessMessage("Setup Complete! Your vault is ready.")
			fmt.Println("Use: 'passlock set <key> <value>' to add entries.")
			helpers.PrintSeparator()
			break
		}
	},
}
