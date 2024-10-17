package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// DeriveKey generates a 256-bit cryptographic key from the given password using SHA-256.
func DeriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// GenerateAESKey creates a random 256-bit AES encryption key.
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key) // Read random bytes into the key.
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt encrypts the provided data using AES encryption in CFB mode.
func Encrypt(data, key []byte) (string, error) {
	block, err := aes.NewCipher(key) // Create a new AES cipher block.
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize] // Create an initialization vector (IV).
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)           // Create a CFB encrypter.
	stream.XORKeyStream(cipherText[aes.BlockSize:], data) // Encrypt the data.

	return hex.EncodeToString(cipherText), nil // Return encrypted data as a hex string.
}

// Decrypt decrypts AES-encrypted data from a hex string in CFB mode.
func Decrypt(ciphertextHex string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex) // Decode hex string to bytes.
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key) // Create a new AES cipher block.
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize] // Extract the IV from the ciphertext.
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv) // Create a CFB decrypter.
	stream.XORKeyStream(ciphertext, ciphertext) // Decrypt the data.

	return string(ciphertext), nil // Return decrypted data as a string.
}
