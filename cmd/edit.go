package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

// Command to edit the value of an existing key.
var EditCommand = types.Command{
	Name:        "edit",
	Description: "Edit the value of an existing key.",
	Usage:       "passlock edit <key>",
	ArgCount:    1,
	Execute: func(args []string) {
		key := args[0]

		// Validate the key format.
		if err := helpers.ValidateInput(key, "Key"); err != nil {
			helpers.ErrorMessage(err.Error())
			helpers.PrintSeparator()
			return
		}

		// Ensure setup is complete.
		if !helpers.VerifySetup() {
			return
		}

		// Verify password and load existing entries.
		entries, derivedKey, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		var oldValue string
		entryFound := false

		// Search for the entry by key.
		for i, entry := range entries {
			if entry.Key == key {
				// Decrypt the existing value.
				oldValue, err = helpers.Decrypt(entry.Value, derivedKey)
				if err != nil {
					helpers.ErrorMessage("Error decrypting the value.")
					helpers.PrintSeparator()
					return
				}

				entryFound = true

				// Prompt user to enter a new value.
				helpers.SuccessMessage(fmt.Sprintf("Editing value for key '%s'", key))
				fmt.Printf("Old value: %s\n", oldValue)
				fmt.Print("Enter new value (leave empty to keep the old value): ")

				newValue, err := helpers.ReadLine()
				if err != nil {
					log.Fatalf("Error reading new value: %v\n", err)
				}

				// If no new value is provided, keep the old one.
				if newValue == "" {
					newValue = oldValue
				}

				// Encrypt the new value and update the entry.
				encryptedValue, err := helpers.Encrypt([]byte(newValue), derivedKey)
				if err != nil {
					log.Fatalf("Error encrypting the new value: %v\n", err)
				}
				entries[i].Value = encryptedValue

				// Save the updated entries.
				err = helpers.SaveToFile(entries, filepath.Join(helpers.GetAppDataPath(), "data.plock"), derivedKey)
				if err != nil {
					log.Fatalf("Error saving updated entry: %v\n", err)
				}

				// Confirm successful update.
				helpers.SuccessMessage("Entry updated successfully.")
				helpers.PrintSeparator()
				return
			}
		}

		// Handle case where the key was not found.
		if !entryFound {
			helpers.ErrorMessage(fmt.Sprintf("Key '%s' not found.", key))
			helpers.PrintSeparator()
		}
	},
}
