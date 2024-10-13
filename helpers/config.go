package helpers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/Ege-Okyay/pass-lock/types"
)

func GetAppDataPath() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		log.Fatal("APPDATA environment variable not found.")
	}
	return filepath.Join(appData, "passlock")
}

func SaveToFile(content interface{}, filename string, aesKey []byte) error {
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}

	encryptedContent, err := Encrypt(jsonContent, aesKey)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(encryptedContent), 0644)
}

func LoadFromFile(filename string, aesKey []byte) ([]types.PlockEntry, error) {
	encryptedContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decryptedContent, err := Decrypt(string(encryptedContent), aesKey)
	if err != nil {
		return nil, err
	}

	var entries []types.PlockEntry
	if err := json.Unmarshal([]byte(decryptedContent), &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func AddDataEntry(aesKey []byte, filename, key, value string) error {
	dataFile := filepath.Join(GetAppDataPath(), filename)
	var entries []types.PlockEntry

	if content, err := LoadFromFile(dataFile, aesKey); err == nil {
		entries = content
	}

	encryptedValue, err := Encrypt([]byte(value), aesKey)
	if err != nil {
		return err
	}

	entries = append(entries, types.PlockEntry{Key: key, Value: encryptedValue})

	return SaveToFile(entries, dataFile, aesKey)
}

func CheckKeysFileExists() (bool, error) {
	keysFile := filepath.Join(GetAppDataPath(), "keys.plock")

	if _, err := os.Stat(keysFile); os.IsNotExist(err) {
		return false, nil
	}

	fileInfo, err := os.Stat(keysFile)
	if err != nil {
		return false, err
	}

	return fileInfo.Size() > 0, nil
}
