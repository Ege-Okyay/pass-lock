package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Ege-Okyay/pass-lock/helpers"
	"github.com/Ege-Okyay/pass-lock/types"
)

var SelfDestructCommand = types.Command{
	Name:        "self-destruct",
	Description: "Delete all stored data and remove passlock configuration.",
	Usage:       "passlock self-destruct",
	ArgCount:    0,
	Execute: func(args []string) {
		if !helpers.VerifySetup() {
			return
		}

		_, _, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		helpers.PrintBanner("WARNING: This will delete ALL your saved data and configuration files.")

		fmt.Print("Type 'sudo delete passlock' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v\n", err)
		}

		input = helpers.TrimNewline(input)

		if input != "sudo delete passlock" {
			helpers.ErrorMessage("Invalid input. Aborting self-destruct.")
			return
		}

		passlockDir := helpers.GetAppDataPath()

		if err := os.RemoveAll(passlockDir); err != nil {
			log.Fatalf("Failed to delete passlock directory: %v\n", err)
		}

		helpers.SuccessMessage("All passlock data has been destroyed. Goodbye!")
	},
}
