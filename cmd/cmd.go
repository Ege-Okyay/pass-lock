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
		- set
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

			confirmPassword, err := helpers.ReadPassword("Confirm password: ")
			if err != nil {
				log.Fatalf("Error reading confirmation password: %v\n", err)
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

			encryptedAESKey, err := helpers.Encrypt(aesKey, derivedKey)
			if err != nil {
				log.Fatalf("Error encrypting AES key: %v\n", err)
			}

			passlockDir := helpers.GetAppDataPath()

			if err := os.MkdirAll(passlockDir, os.ModePerm); err != nil {
				log.Fatalf("Error creating passlock directory: %v\n", err)
			}

			keysFile := filepath.Join(passlockDir, "keys.plock")
			if err := helpers.SaveToFile(encryptedAESKey, keysFile); err != nil {
				log.Fatalf("Error saving encrypted AES key: %v\n", err)
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
