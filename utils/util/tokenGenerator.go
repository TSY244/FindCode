package util

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
)

// GenerateToken Generate a token
func GenerateToken() string {
	// Define the character pool
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/"

	// Generate 32 random bytes
	tokenLength := 32
	randomBytes := make([]byte, tokenLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %v", err)
	}

	// Map random bytes to character pool
	token := make([]byte, tokenLength)
	temp := byte(len(charset))
	for i, b := range randomBytes {
		token[i] = charset[b%temp]
	}

	base64Token := Base64Encode(string(token))

	file, err := os.Create("token.txt")
	if err != nil {
		log.Fatalf("GenerateToken -> os.Create: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(base64Token)
	if err != nil {
		log.Fatalf("GenerateToken -> file.WriteString: %v", err)
	}

	return base64Token
}

// Base64Encode Base64 encode the input string
func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
