package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"os"
)

func GenerateToken(b int) {
	// Set colors
	red := color.New(color.FgHiRed).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()

	// Check if bytes are provided
	if b <= 0 {
		fmt.Printf("\n%s Please specify the number of bytes using the %s flag.\n\n", red("Error:"), yellow("-b"))
		os.Exit(1)
	}

	// Generate random bytes
	randomBytes := make([]byte, b)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Printf("\n%s Unable to generate random bytes - %s.\n\n", red("Error:"), err)
		os.Exit(1)
	}

	// Base64 URL encode without padding
	base64URL := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	// Hex encode
	hexEncoded := hex.EncodeToString(randomBytes)

	// Output the results with yellow color for "key:" and "hex:"
	fmt.Printf("\n%s %s\n", yellow("Key:"), base64URL)
	fmt.Printf("%s %s\n\n", yellow("Hex:"), hexEncoded)
}
