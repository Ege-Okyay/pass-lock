package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

// Command to delete a key-value pair from the vault.
var DeleteCommand = types.Command{
	Name:        "delete",
	Description: "Delete a key-value pair from the vault.",
	Usage:       "passlock delete <key>",
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

		// Filter out the entry to delete, if it exists.
		var updatedEntries []types.PlockEntry
		found := false
		for _, entry := range entries {
			if entry.Key == key {
				found = true
			} else {
				updatedEntries = append(updatedEntries, entry)
			}
		}

		// Handle case where the key was not found.
		if !found {
			helpers.ErrorMessage(fmt.Sprintf("Key '%s' not found.", key))
			helpers.PrintSeparator()
			return
		}

		// Save the updated entries to the file.
		err = helpers.SaveToFile(updatedEntries, filepath.Join(helpers.GetUserConfigDir(), "data.plock"), derivedKey)
		if err != nil {
			log.Fatalf("Error saving updated entries: %v\n", err)
		}

		// Confirm successful deletion.
		helpers.SuccessMessage(fmt.Sprintf("Key '%s' deleted successfully.", key))
	},
}
