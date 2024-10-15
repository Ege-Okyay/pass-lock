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

		exists, err := helpers.CheckKeysFileExists()
		if err != nil {
			log.Fatalf("Error checking keys file: %v\n", err)
		}
		if !exists {
			helpers.ErrorMessage("Setup is not completed. Please use the 'setup' command to initialize.")
			helpers.PrintSeparator()
			return
		}

		entries, derivedKey, err := helpers.VerifyPasswordAndLoadData("data.plock")
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		for _, entry := range entries {
			if entry.Key == key {
				storedEntry, err := helpers.Decrypt(entry.Value, derivedKey)
				if err != nil {
					helpers.ErrorMessage("I really don't know what this error should be")
					helpers.PrintSeparator()
					continue
				}

				helpers.SuccessMessage("Found the value!")
				helpers.PrintBanner(storedEntry)
				return
			} else {
				helpers.ErrorMessage(fmt.Sprintf("Key '%s' not found.", key))
				helpers.PrintSeparator()
				return
			}
		}
	},
}
