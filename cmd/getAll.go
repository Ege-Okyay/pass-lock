package cmd

import (
	"fmt"
	"log"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var GetAllCommand = types.Command{
	Name:        "get-all",
	Description: "Retrive and display all key-value pairs from the data vault.",
	Usage:       "passlock get-all",
	ArgCount:    0,
	Execute: func(args []string) {
		status := helpers.VerifySetup()
		if !status {
			return
		}

		entries, derivedKey, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		if len(entries) == 0 {
			helpers.ErrorMessage("The vault is empty. Add entries with 'passlock set <key> <value>'.")
			helpers.PrintSeparator()
			return
		}

		helpers.PrintBanner("Displaying all stored key-value pairs:")
		for _, entry := range entries {
			decryptedValue, err := helpers.Decrypt(entry.Value, derivedKey)
			if err != nil {
				helpers.ErrorMessage(fmt.Sprintf("Failed to decrpyt value for key '%s'.", entry.Key))
				helpers.PrintSeparator()
				continue
			}
			fmt.Printf("KEY: %-15s VALUE: %s\n", entry.Key, decryptedValue)
		}
	},
}
