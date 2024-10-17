package helpers

import "fmt"

// PrintBanner prints a formatted banner with a given title.
func PrintBanner(title string) {
	fmt.Println("========================================")
	fmt.Printf("%s\n", title) // Display the title in the center.
	fmt.Println("========================================")
}

// PrintSeparator prints a horizontal line for visual separation of content.
func PrintSeparator() {
	fmt.Println("----------------------------------------")
}

// SuccessMessage prints a success message with a checkmark.
func SuccessMessage(msg string) {
	// Unicode checkmark (✔) is used to indicate success.
	fmt.Printf("\xE2\x9C\x94 %s\n", msg)
}

// ErrorMessage prints an error message with a cross mark.
func ErrorMessage(msg string) {
	// Unicode cross (✗) is used to indicate an error.
	fmt.Printf("\xE2\x9C\x97 %s\n", msg)
}
