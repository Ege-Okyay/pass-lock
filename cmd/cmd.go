package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

/*
	Commands List:
		- setup DONE
		- set DONE
		- get
		- get-all
		- delete
		- self-destruct
*/

var HelpCommand = types.Command{
	Name:        "help",
	Description: "change this later",
	Usage:       "passlock help",
	Execute: func(args []string) {
		fmt.Println("Hello World.")
	},
}

var SetupCommand = types.Command{
	Name:        "setup",
	Description: "Set up your passlock with a new password.",
	Usage:       "passlock setup",
	Execute: func(args []string) {
		exists, err := helpers.CheckKeysFileExists()
		if err != nil {
			log.Fatalf("Error checking keys file: %v\n", err)
		}

		if exists {
			helpers.ErrorMessage("Setup already completed. Please use the 'set' command to add entries.")
			return
		}

		helpers.PrintBanner("Welcome to Passlock")
		fmt.Println("Let's setup up your password manager.")
		fmt.Println("\n🚀 Setup will begin shortly...")

		for {
			password, err := helpers.ReadPassword("Password: ")
			if err != nil {
				log.Fatalf("Error reading password: %v\n", err)
			}

			if err := helpers.ValidateInput(password, "Password"); err != nil {
				helpers.ErrorMessage(err.Error())
				helpers.PrintSeparator()
				continue
			}

			confirmPassword, err := helpers.ReadPassword("Confirm password: ")
			if err != nil {
				log.Fatalf("Error reading confirmation password: %v\n", err)
			}

			if err := helpers.ValidateInput(confirmPassword, "Confirmation password"); err != nil {
				helpers.ErrorMessage(err.Error())
				helpers.PrintSeparator()
				continue
			}

			if password != confirmPassword {
				helpers.ErrorMessage("Passwords do not match! Please try again.")
				helpers.PrintSeparator()
				continue
			}

			derivedKey := helpers.DeriveKey(password)

			aesKey, err := helpers.GenerateAESKey()
			if err != nil {
				log.Fatalf("Error generating AES key: %v\n", err)
			}

			passlockDir := helpers.GetAppDataPath()

			if err := os.MkdirAll(passlockDir, os.ModePerm); err != nil {
				log.Fatalf("Error creating passlock directory: %v\n", err)
			}

			if err := helpers.AddDataEntry(derivedKey, "keys.plock", "aes_key", string(aesKey)); err != nil {
				log.Fatalf("Error saving AES key: %v\n", err)
			}

			if err := helpers.AddDataEntry(derivedKey, "keys.plock", "password", password); err != nil {
				log.Fatalf("Error saving password: %v\n", err)
			}

			helpers.SuccessMessage("Setup Complete! Your vault is ready.")
			fmt.Println("Use: 'passlock set <key> <value>' to add entries.")
			helpers.PrintSeparator()
			break
		}
	},
}

var SetCommand = types.Command{
	Name:        "set",
	Description: "change this later",
	Usage:       "passlock set <key> <value>",
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
