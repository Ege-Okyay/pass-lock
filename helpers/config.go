package helpers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Ege-Okyay/passlock/types"
)

// GetUserConfigDir returns the path to the user's configuration directory
// for the application, handling different operating systems.
func GetUserConfigDir() string {
	var configDir string

	// Determine the operating system
	switch runtime.GOOS {
	case "windows":
		// Windows uses APPDATA environment variable
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			log.Fatal("APPDATA environment variable not found.")
		}
	case "darwin":
		// macOS uses a standard path in the user's home directory
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	default:
		log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
	}

	return filepath.Join(configDir, "passlock")
}

// SaveToFile encrypts the given content with the AES key and writes it to a file.
func SaveToFile(content interface{}, filename string, aesKey []byte) error {
	// Convert the content to JSON format.
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// Encrypt the JSON content with the provided AES key.
	encryptedContent, err := Encrypt(jsonContent, aesKey)
	if err != nil {
		return err
	}

	// Write the encrypted content to the specified file.
	return os.WriteFile(filename, []byte(encryptedContent), 0644)
}

// LoadFromFile reads, decrypts, and unmarshals JSON content from a file.
func LoadFromFile(filename string, aesKey []byte) ([]types.PlockEntry, error) {
	// Check if the file exists; if not, return no entries.
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// If the file is empty, return an empty slice of entries.
	if fileInfo.Size() == 0 {
		return []types.PlockEntry{}, nil
	}

	// Read the encrypted content from the file.
	encryptedContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Decrypt the content using the AES key.
	decryptedContent, err := Decrypt(string(encryptedContent), aesKey)
	if err != nil {
		return nil, err
	}

	// Unmarshal the decrypted content into a slice of PlockEntry structs.
	var entries []types.PlockEntry
	if err := json.Unmarshal([]byte(decryptedContent), &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// AddDataEntry appends a new encrypted key-value pair to the specified data file.
func AddDataEntry(aesKey []byte, filename, key, value string) error {
	// Determine the full path to the data file.
	dataFile := filepath.Join(GetUserConfigDir(), filename)

	// Load existing entries from the file, if available.
	var entries []types.PlockEntry
	if content, err := LoadFromFile(dataFile, aesKey); err == nil {
		entries = content
	}

	// Encrypt the new value with the AES key.
	encryptedValue, err := Encrypt([]byte(value), aesKey)
	if err != nil {
		return err
	}

	// Append the new entry to the list of entries.
	entries = append(entries, types.PlockEntry{Key: key, Value: encryptedValue})

	// Save the updated entries back to the file.
	return SaveToFile(entries, dataFile, aesKey)
}

// CheckKeysFileExists verifies whether the 'keys.plock' file exists and is not empty.
func CheckKeysFileExists() (bool, error) {
	// Get the full path to the keys file.
	keysFile := filepath.Join(GetUserConfigDir(), "keys.plock")

	// Check if the file exists; if not, return false.
	if _, err := os.Stat(keysFile); os.IsNotExist(err) {
		return false, nil
	}

	// Retrieve file information and handle errors.
	fileInfo, err := os.Stat(keysFile)
	if err != nil {
		return false, err
	}

	// Return true if the file has content.
	return fileInfo.Size() > 0, nil
}

// VerifySetup checks if the vault setup is completed by verifying the existence of 'keys.plock'.
func VerifySetup() bool {
	// Check if the keys file exists and is not empty.
	exists, err := CheckKeysFileExists()
	if err != nil {
		log.Fatalf("Error checking keys file: %v\n", err)
	}

	// If the setup is incomplete, print an error message.
	if !exists {
		ErrorMessage("Setup is not completed. Please use the 'setup' command to initialize.")
		PrintSeparator()
		return false
	}

	// Return true if the setup is completed.
	return true
}
