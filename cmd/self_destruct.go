package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Ege-Okyay/passlock/helpers"
	"github.com/Ege-Okyay/passlock/types"
)

// Command to delete all data and remove passlock configuration.
var SelfDestructCommand = types.Command{
	Name:        "self-destruct",
	Description: "Delete all stored data and remove passlock configuration.",
	Usage:       "passlock self-destruct",
	ArgCount:    0,
	Execute: func(args []string) {
		// Ensure setup is complete.
		if !helpers.VerifySetup() {
			return
		}

		// Verify password before proceeding.
		_, _, err := helpers.VerifyPasswordAndLoadData()
		if err != nil {
			log.Fatalf("Password verification failed: %v\n", err)
		}

		helpers.PrintBanner("WARNING: This will delete ALL your saved data and configuration files.")

		// Prompt user to confirm the destructive action.
		fmt.Print("Type 'sudo delete passlock' to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v\n", err)
		}
		input = helpers.TrimNewline(input)

		// Abort if the confirmation is incorrect.
		if input != "sudo delete passlock" {
			helpers.ErrorMessage("Invalid input. Aborting self-destruct.")
			return
		}

		// Delete all data by removing the passlock directory.
		passlockDir := helpers.GetUserConfigDir()
		if err := os.RemoveAll(passlockDir); err != nil {
			log.Fatalf("Failed to delete passlock directory: %v\n", err)
		}

		// Confirm successful destruction of data.
		helpers.SuccessMessage("All passlock data has been destroyed. Goodbye!")
	},
}
