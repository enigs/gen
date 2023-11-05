package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"os"
)

func GenerateToken() {
	// Set colors
	red := color.New(color.FgHiRed).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()

	//// Check if bytes are provided
	//if b <= 0 {
	//	fmt.Printf("\n%s Please specify the number of bytes using the %s flag.\n\n", red("Error:"), yellow("-b"))
	//	os.Exit(1)
	//}

	// Generate a new random 24-byte nonce.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// Generate masterKey
	var masterKey [32]byte
	if _, err := io.ReadFull(rand.Reader, masterKey[:]); err != nil {
		fmt.Printf("\n%s Unable to generate random bytes - %s.\n\n", red("Error:"), err)
		os.Exit(1)
	}

	// Copy masterKey to a byte slice
	masterKeyBytes := masterKey[:]

	// Base64 URL encode without padding
	b64MasterKey := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(masterKeyBytes)

	// Output the results with yellow color for "key:" and "hex:"
	fmt.Printf("\n%s %s\n", yellow("Master Key:"), b64MasterKey)

	// Encrypt the message and prepend the nonce.
	generatedMasterKey := secretbox.Seal(nonce[:], masterKeyBytes, &nonce, &masterKey)

	// Tokenize master key
	tokenizedMasterKey := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(generatedMasterKey)

	// Output the results with yellow color for "key:" and "hex:"
	fmt.Printf("%s %s\n\n", yellow("Token Key:"), tokenizedMasterKey)
}
