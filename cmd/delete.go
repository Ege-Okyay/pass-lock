package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var DeleteCommand = types.Command{
	Name:        "delete",
	Description: "Delete a key-value pair from the vault.",
	Usage:       "passlock delete <key>",
	ArgCount:    1,
	Execute: func(args []string) {
		key := args[0]

		if err := helpers.ValidateInput(key, "Key"); err != nil {
			helpers.ErrorMessage(err.Error())
			helpers.PrintSeparator()
			return
		}

		status := helpers.VerifySetup()
		if !status {
			return
		}

		entries, derivedKey, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		var updatedEntries []types.PlockEntry
		found := false
		for _, entry := range entries {
			if entry.Key == key {
				found = true
			} else {
				updatedEntries = append(updatedEntries, entry)
			}
		}

		if !found {
			helpers.ErrorMessage(fmt.Sprintf("Key '%s' not found.", key))
			helpers.PrintSeparator()
			return
		}

		err = helpers.SaveToFile(updatedEntries, filepath.Join(helpers.GetAppDataPath(), "data.plock"), derivedKey)
		if err != nil {
			log.Fatalf("Error saving updated entries: %v\n", err)
		}

		helpers.SuccessMessage(fmt.Sprintf("Key '%s' deleted successfully.", key))
	},
}
