package cmd

import (
	"fmt"
	"log"

	"github.com/Ege-Okyay/passlock/helpers"
	"github.com/Ege-Okyay/passlock/types"
)

// Command to retrieve and display all key-value pairs from the vault.
var GetAllCommand = types.Command{
	Name:        "get-all",
	Description: "Retrieve and display all key-value pairs from the data vault.",
	Usage:       "passlock get-all",
	ArgCount:    0,
	Execute: func(args []string) {
		// Ensure setup is complete.
		if !helpers.VerifySetup() {
			return
		}

		// Load entries after password verification.
		entries, derivedKey, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		// Handle case where the vault is empty.
		if len(entries) == 0 {
			helpers.ErrorMessage("The vault is empty. Add entries with 'passlock set <key> <value>'.")
			helpers.PrintSeparator()
			return
		}

		helpers.PrintBanner("Displaying all stored key-value pairs:")
		for _, entry := range entries {
			// Attempt to decrypt each entry's value.
			decryptedValue, err := helpers.Decrypt(entry.Value, derivedKey)
			if err != nil {
				helpers.ErrorMessage(fmt.Sprintf("Failed to decrypt value for key '%s'.", entry.Key))
				helpers.PrintSeparator()
				continue
			}
			fmt.Printf("KEY: %-15s VALUE: %s\n", entry.Key, decryptedValue)
		}
	},
}
