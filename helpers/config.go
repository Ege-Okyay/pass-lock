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

func SaveToFile(content, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func LoadFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func AddDataEntry(aesKey []byte, key, value string) error {
	dataFile := filepath.Join(GetAppDataPath(), "data.plock")
	var entries []types.DataEntry

	if content, err := LoadFromFile(dataFile); err == nil {
		json.Unmarshal([]byte(content), &entries)
	}

	encryptedValue, err := Encrypt([]byte(value), aesKey)
	if err != nil {
		return err
	}

	entries = append(entries, types.DataEntry{Key: key, Value: encryptedValue})

	jsonContent, _ := json.Marshal(entries)

	return SaveToFile(string(jsonContent), dataFile)
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
