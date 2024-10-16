package cmd

import (
	"log"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var SetCommand = types.Command{
	Name:        "set",
	Description: "Store a new key-value pair.",
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

		if !helpers.VerifySetup() {
			return
		}

		_, derivedKey, err := helpers.VerifyPasswordAndLoadData()
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
