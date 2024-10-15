package cmd

import (
	"log"

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
			helpers.ErrorMessage("Setup is not completed. Please use the 'setup' command to initialize.")
			helpers.PrintSeparator()
			return
		}

		_, derivedKey, err := helpers.VerifyPasswordAndLoadData("data.plock")
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		err = helpers.AddDataEntry(derivedKey, "data.plock", key, value)
		if err != nil {
			log.Fatalf("Error adding new entry: %v\n", err)
		}

		helpers.SuccessMessage("Entry added successfully.")
	},
}
