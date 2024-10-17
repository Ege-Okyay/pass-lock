package cmd

import (
	"fmt"
	"log"

	"github.com/Ege-Okyay/passlock/helpers"
	"github.com/Ege-Okyay/passlock/types"
)

// Command to retrieve and display a specific key-value pair.
var GetCommand = types.Command{
	Name:        "get",
	Description: "Retrieve and display the specified key-value pair from the data vault.",
	Usage:       "passlock get <key>",
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

		// Load entries after password verification.
		entries, derivedKey, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		// Search for the key in the vault.
		entryFound := false
		for _, entry := range entries {
			if entry.Key == key {
				// Decrypt the value if the key is found.
				storedEntry, err := helpers.Decrypt(entry.Value, derivedKey)
				if err != nil {
					helpers.ErrorMessage("Error decrypting the value.")
					helpers.PrintSeparator()
					continue
				}

				entryFound = true
				helpers.SuccessMessage("Found the value!")
				helpers.PrintBanner(storedEntry)
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
