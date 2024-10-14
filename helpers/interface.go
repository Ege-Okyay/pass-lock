package helpers

import "fmt"

func PrintBanner(title string) {
	fmt.Println("========================================")
	fmt.Printf("      🔐 %s 🔐\n", title)
	fmt.Println("========================================")
}

func PrintSeparator() {
	fmt.Println("----------------------------------------")
}

func SuccessMessage(msg string) {
	fmt.Printf("✅ %s\n", msg)
}

func ErrorMessage(msg string) {
	fmt.Printf("❌ %s\n", msg)
}
