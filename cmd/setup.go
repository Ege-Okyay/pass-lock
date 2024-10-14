package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var SetupCommand = types.Command{
	Name:        "setup",
	Description: "Set up your passlock with a new password.",
	Usage:       "passlock setup",
	ArgCount:    0,
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
