package cmd

import (
	"log"
	"path/filepath"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var SetCommand = types.Command{
	Name:        "set",
	Description: "change this later",
	Usage:       "passlock set <key> <value>",
	ArgCount:    2,
	Execute: func(args []string) {
		key, value := args[0], args[1]

		if err := helpers.ValidateInput(key, "Key"); err != nil {
			helpers.ErrorMessage(err.Error())
			helpers.PrintSeparator()
			return
		}

		if err := helpers.ValidateInput(value, "Value"); err != nil {
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

		for {
			password, err := helpers.ReadPassword("Password: ")
			if err != nil {
				log.Fatalf("Error reading password: %v\n", err)
			}

			derivedKey := helpers.DeriveKey(password)

			entries, err := helpers.LoadFromFile(filepath.Join(helpers.GetAppDataPath(), "keys.plock"), derivedKey)
			if err != nil {
				helpers.ErrorMessage("Incorrect password. Please try again.")
				helpers.PrintSeparator()
				continue
			}

			var storedPassword string
			for _, entry := range entries {
				if entry.Key == "password" {
					storedPassword, err = helpers.Decrypt(entry.Value, derivedKey)
					if err != nil {
						helpers.ErrorMessage("Incorrect password. Please try again.")
						helpers.PrintSeparator()
						continue
					}
					break
				}
			}

			if storedPassword == password {
				helpers.SuccessMessage("Password verified! Adding new entry...")

				err := helpers.AddDataEntry(derivedKey, "data.plock", key, value)
				if err != nil {
					log.Fatalf("Error adding new entry: %v\n", err)
				}

				helpers.SuccessMessage("Entry added successfully.")
				break
			} else {
				helpers.ErrorMessage("Incorrect password. Please try again.")
				helpers.PrintSeparator()
			}
		}
	},
}
