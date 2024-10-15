package cmd

import (
	"log"
	"path/filepath"

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
			helpers.ErrorMessage("Setup is not completed. Please use the 'setup' command to setup passlock.")
			helpers.PrintSeparator()
			return
		}

		_, derivedKey, err := helpers.VerifyPassword()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		entries, err := helpers.LoadFromFile(filepath.Join(helpers.GetAppDataPath(), "data.plock"), derivedKey)
		if err != nil {
			helpers.ErrorMessage("Not sure what should I write here")
			helpers.PrintSeparator()
			return
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
			}
		}
	},
}
