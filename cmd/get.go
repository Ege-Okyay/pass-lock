package cmd

import (
	"fmt"
	"log"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var GetCommand = types.Command{
	Name:        "get",
	Description: "change this later",
	Usage:       "passlock get <key>",
	ArgCount:    1,
	Execute: func(args []string) {
		key := args[0]

		if err := helpers.ValidateInput(key, "Key"); err != nil {
			helpers.ErrorMessage(err.Error())
			helpers.PrintSeparator()
			return
		}

		if !helpers.VerifySetup() {
			return
		}

		entries, derivedKey, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		entryFound := false
		for _, entry := range entries {
			if entry.Key == key {
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

		if !entryFound {
			helpers.ErrorMessage(fmt.Sprintf("Key '%s' not found.", key))
			helpers.PrintSeparator()
		}
	},
}
