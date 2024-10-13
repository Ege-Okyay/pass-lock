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
			fmt.Println("Setup already completed. Please use the 'set' command to add entries.")
			return
		}

		fmt.Println("=======================================")
		fmt.Println("           Welcome to passlock!       ")
		fmt.Println("=======================================")
		fmt.Println("We will get started by creating a password.")
		fmt.Println()

		for {
			password, err := helpers.ReadPassword("Password: ")
			if err != nil {
				log.Fatalf("Error reading password: %v\n", err)
			}

			if err := helpers.ValidateInput(password, "Password"); err != nil {
				fmt.Println(err)
				continue
			}

			confirmPassword, err := helpers.ReadPassword("Confirm password: ")
			if err != nil {
				log.Fatalf("Error reading confirmation password: %v\n", err)
			}

			if err := helpers.ValidateInput(confirmPassword, "Confirmation password"); err != nil {
				fmt.Println(err)
				continue
			}

			if password != confirmPassword {
				fmt.Println()
				fmt.Println("Passwords do not match! Please try again.")
				fmt.Println("---------------------------------------")
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

			fmt.Println()
			fmt.Println("=======================================")
			fmt.Println("         Setup Complete!              ")
			fmt.Println("=======================================")
			fmt.Println("You can start using passlock by running 'passlock set <key> <value>'")
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
			fmt.Println(err)
			return
		}

		if err := helpers.ValidateInput(value, "Value"); err != nil {
			fmt.Println(err)
			return
		}

		exists, err := helpers.CheckKeysFileExists()
		if err != nil {
			log.Fatalf("Error checking keys file: %v\n", err)
		}

		if !exists {
			fmt.Println("Setup is not completed. Please use the 'setup' command to setup passlock.")
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
				fmt.Println("Error loading keys file or incorrect password. Please try again.")
				continue
			}

			var storedPassword string
			for _, entry := range entries {
				if entry.Key == "password" {
					storedPassword, err = helpers.Decrypt(entry.Value, derivedKey)
					if err != nil {
						fmt.Println("Incorrect password. Please try again.")
						continue
					}
					break
				}
			}

			if storedPassword == password {
				fmt.Println("Password verified! Adding new entry...")

				err := helpers.AddDataEntry(derivedKey, "data.plock", key, value)
				if err != nil {
					log.Fatalf("Error adding new entry: %v\n", err)
				}

				fmt.Println("Entry added successfully.")
				break
			} else {
				fmt.Println("Incorrect password. Please try again.")
			}
		}
	},
}
